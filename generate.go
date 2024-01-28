package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	generator "vpat_codegen/generators"
	"vpat_codegen/model"

	"archive/zip"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	app := fiber.New()

	// API endpoint to receive StructDefinition with BodyParser and Database with QueryParser
	app.Post("/generate", func(c *fiber.Ctx) error {

		var appJson model.AppJson
		if err := c.BodyParser(&appJson); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		database := c.Query("database")

		err := GenerateApplicationCode(appJson, database)

		if err != nil {
			fmt.Println("Unable to generate application code : " + err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to generate application code : " + err.Error())
		}

		zipFile, err := CreateApplicationZip(appJson.AppName)

		if err != nil {
			fmt.Println("Unable to zip application code  : " + err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to zip application code  : " + err.Error())
		}

		err = os.RemoveAll("./generated")
		if err != nil {
			fmt.Println("Unable to clean generated directory  : " + err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to clean generated directory  : " + err.Error())
		}

		// Set appropriate headers for download
		c.Set("Content-Disposition", "attachment; filename="+appJson.AppName+".zip")
		c.Set("Content-Type", "application/zip")

		file, err := os.ReadFile(zipFile)
		if err != nil {
			fmt.Println("Unable to read generated zip file  : " + err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to read generated zip file  : " + err.Error())
		}

		// Send the zip file as response
		c.Response().SetBodyRaw(file)

		err = os.Remove(zipFile)
		if err != nil {
			fmt.Println("Unable to delete generated zip file  : " + err.Error())
		}

		return c.SendStatus(fiber.StatusOK)

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
	routesFile, err := os.Create("generated/routes/routes.go")
	if err != nil {
		panic("Unable to create routes file")
	}

	defer routesFile.Close()

	mainFile, err := os.Create("generated/main.go")
	if err != nil {
		panic("Unable to create routes file")
	}

	defer mainFile.Close()
}

func updateModFile() error {
	modFile, err := os.OpenFile("generated/go.mod", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open go.mod")
		return err
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
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.6
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)`)

	if err2 != nil {
		fmt.Println("Could not write text to go.mod")
		return err2
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
		fmt.Println("Go mod tidy failed:" + err.Error())
		return err
	}

	return nil
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

func GenerateApplicationCode(appJson model.AppJson, database string) error {
	// Parse the incoming JSON data as a StructDefinition
	var structDefs []model.StructDefinition

	structDefs = appJson.Tables

	// Parse the "database" query parameter
	if database != "postgres" && database != "mysql" {
		fmt.Println("Unabel to process request : Invalid database type")
		return errors.New("Only supported databases are mysql and postgres.")
	}

	createFiles(appJson.AppName)

	err := generator.CreateDatabase(database, cases.Title(language.English).String(appJson.AppName), structDefs, appJson.AppName)
	if err != nil {
		panic("Unabel to create and connect to database")
	}

	//creates model,controlles and service files
	_, err = CreateServices(structDefs, database, appJson.AppName)
	if err != nil {
		fmt.Println("Unabel to generate services : " + err.Error())
		return err
	}

	// Update routes.go to define API endpoints
	err = generator.UpdateRoutesFile(structDefs, database, appJson.AppName)
	if err != nil {
		fmt.Println("Unabel to generate routes : " + err.Error())
		return err
	}

	err = generator.CreateMainFile(structDefs, database, appJson.AppName)
	if err != nil {
		fmt.Println("Unabel to generate main.go : " + err.Error())
		return err
	}

	err = updateModFile()

	if err != nil {
		panic(err)
	}
	fmt.Println("Code generation completed.")
	return nil
}

func CreateApplicationZip(appName string) (string, error) {
	// Directory to zip
	dirToZip := "./generated"

	// Create a temporary zip file
	zipFile := appName + ".zip"
	zipWriter, err := os.Create(zipFile)
	if err != nil {
		return "", err
	}
	defer zipWriter.Close()

	// Create a new zip archive
	zipArchive := zip.NewWriter(zipWriter)
	defer zipArchive.Close()

	// Walk through the directory and add files to zip
	err = filepath.Walk(dirToZip, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(dirToZip, filePath)
			if err != nil {
				return err
			}
			fileToZip, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer fileToZip.Close()

			// Create a new file in the zip archive
			zipFile, err := zipArchive.Create(relPath)
			if err != nil {
				return err
			}

			// Copy the file to the zip writer
			_, err = io.Copy(zipFile, fileToZip)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return zipFile, nil
}
