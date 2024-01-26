package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	generator "vpat_codegen/generators"
	"vpat_codegen/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	// API endpoint to receive StructDefinition with BodyParser and Database with QueryParser
	app.Post("/generate", func(c *fiber.Ctx) error {
		// Parse the incoming JSON data as a StructDefinition
		var structDefs []model.StructDefinition
		var appJson model.AppJson
		if err := c.BodyParser(&appJson); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		structDefs = appJson.Tables

		// Parse the "database" query parameter
		database := c.Query("database")
		if database != "postgres" && database != "mysql" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid database type"})
		}

		createFiles(appJson.AppName)
		err := generator.CreateDatabase(database, cases.Title(language.English).String(appJson.AppName))
		if err != nil {
			panic("Unabel to create and connect to database")
		}

		//creates model,controlles and service files
		errMsgMap, err := CreateServices(structDefs, database, appJson.AppName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errMsgMap)
		}

		// Update routes.go to define API endpoints
		err = generator.UpdateRoutesFile(structDefs, database, appJson.AppName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate routes file"})
		}

		err = generator.CreateMainFile(structDefs, database, appJson.AppName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate main file"})
		}

		updateModFile()

		return c.JSON(fiber.Map{"message": "CRUD code generated and organized successfully"})
	})

	// Start the Fiber server on port 3000
	app.Listen(":3000")
}

func createFiles(name string) {
	err := os.MkdirAll("./generated", os.ModeDir)

	if err != nil {
		panic("Unable to create generated dir")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	initCommand := exec.Command("go", "mod", "init", name)

	initCommand.Dir = "./generated/"
	initCommand.Stdin = os.Stdin
	initCommand.Stdout = &stdout
	initCommand.Stderr = &stderr

	err = initCommand.Run()
	if err != nil {
		fmt.Println("Err :" + err.Error())
		panic("Go mod init failed")
	}

	err = os.MkdirAll("generated/databases", os.ModeDir)
	if err != nil {
		panic("Unable to create databases dir")
	}

	err = os.MkdirAll("generated/service", os.ModeDir)

	if err != nil {
		panic("Unable to create service dir")
	}

	err = os.MkdirAll("generated/controller", os.ModeDir)
	if err != nil {
		panic("Unable to create controller dir")
	}

	err = os.MkdirAll("generated/model", os.ModeDir)
	if err != nil {
		panic("Unable to create model dir")
	}

	err = os.MkdirAll("generated/routes", os.ModeDir)
	if err != nil {
		panic("Unable to create model dir")
	}
	_, err = os.Create("generated/routes/routes.go")
	if err != nil {
		panic("Unable to create routes file")
	}

	_, err = os.Create("generated/main.go")
	if err != nil {
		panic("Unable to create routes file")
	}
}

func updateModFile() {
	modFile, err := os.OpenFile("generated/go.mod", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open go.mod")
	}

	defer modFile.Close()

	_, err2 := modFile.WriteString(`
		

		require (
			github.com/gofiber/fiber/v2 v2.49.2
			golang.org/x/text v0.8.0
		)
		
		require (
			github.com/andybalholm/brotli v1.0.5 // indirect
			github.com/go-sql-driver/mysql v1.7.1
			github.com/google/uuid v1.3.1 // indirect
			github.com/klauspost/compress v1.17.0 // indirect
			github.com/mattn/go-colorable v0.1.13 // indirect
			github.com/mattn/go-isatty v0.0.19 // indirect
			github.com/mattn/go-runewidth v0.0.15 // indirect
			github.com/rivo/uniseg v0.4.4 // indirect
			github.com/valyala/bytebufferpool v1.0.0 // indirect
			github.com/valyala/fasthttp v1.50.0 // indirect
			github.com/valyala/tcplisten v1.0.0 // indirect
			golang.org/x/sys v0.12.0 // indirect
		)
		
		`)

	if err2 != nil {
		fmt.Println("Could not write text to go.mod")

	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	initCommand := exec.Command("go", "mod", "tidy")

	initCommand.Dir = "./generated/"
	initCommand.Stdin = os.Stdin
	initCommand.Stdout = &stdout
	initCommand.Stderr = &stderr

	err = initCommand.Run()
	if err != nil {
		fmt.Println("Err :" + err.Error())
		panic("Go mod tidy failed")
	}
}

func CreateServices(structDefs []model.StructDefinition, database string, appName string) (fiber.Map, error) {
	for _, structDef := range structDefs {
		// Generate Go struct definition
		structCode, err := generator.GenerateStructFromJSON(structDef.JSONExample, structDef.StructName)
		if err != nil {
			return fiber.Map{"error": "Failed to generate struct"}, err
		}

		// Generate CRUD methods and save in service package
		serviceFileName := fmt.Sprintf("generated/service/%s_service.go", strings.ToLower(structDef.StructName))
		if err := generator.GenerateServiceFile(serviceFileName, structDef, structCode, database, appName); err != nil {
			fmt.Sprintln(err.Error())
			return fiber.Map{"error": "Failed to generate service file"}, err
		}

		// Generate controller methods and save in controller package
		controllerFileName := fmt.Sprintf("generated/controller/%s_controller.go", strings.ToLower(structDef.StructName))
		if err := generator.GenerateControllerFile(controllerFileName, structDef, appName); err != nil {
			return fiber.Map{"error": "Failed to generate controller file"}, err
		}

	}

	return nil, nil
}
