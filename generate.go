package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	tmpl "vpat_codegen/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ModelType struct {
	StructName string
	Fields     []Field
}

// Field represents the structure of a field in a Go struct.
type Field struct {
	Name string
	Type string
}

type AppJson struct {
	AppName string             `json:"appName"`
	Tables  []StructDefinition `json:"tables"`
}

// StructDefinition represents the data required for generating CRUD methods.
type StructDefinition struct {
	StructName  string                 `json:"name"`
	JSONExample map[string]interface{} `json:"columns"`
	Endpoint    string                 `json:"endpoint"`
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	// API endpoint to receive StructDefinition with BodyParser and Database with QueryParser
	app.Post("/generate", func(c *fiber.Ctx) error {
		// Parse the incoming JSON data as a StructDefinition
		var structDefs []StructDefinition
		var appJson AppJson
		if err := c.BodyParser(&appJson); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		structDefs = appJson.Tables

		// Parse the "database" query parameter
		database := c.Query("database")
		if database != "postgres" && database != "mysql" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid database type"})
		}

		for _, structDef := range structDefs {
			// Generate Go struct definition
			structCode, err := generateStructFromJSON(structDef.JSONExample, structDef.StructName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate struct"})
			}

			// Generate CRUD methods and save in service package
			serviceFileName := fmt.Sprintf("service/%s_service.go", strings.ToLower(structDef.StructName))
			if err := generateServiceFile(serviceFileName, structDef, structCode, database, appJson.AppName); err != nil {
				fmt.Sprintln(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate service file"})
			}

			// Generate controller methods and save in controller package
			controllerFileName := fmt.Sprintf("controller/%s_controller.go", strings.ToLower(structDef.StructName))
			if err := generateControllerFile(controllerFileName, structDef, appJson.AppName); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate controller file"})
			}

		}
		// Update routes.go to define API endpoints
		err := updateRoutesFile(structDefs, database)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate routes file"})
		}

		return c.JSON(fiber.Map{"message": "CRUD code generated and organized successfully"})
	})

	// Start the Fiber server on port 3000
	app.Listen(":3000")
}

func createFiles(name string) {
	err := os.MkdirAll("generated", os.ModeDir)

	if err != nil {
		panic("Unable to create generated dir")
	}

	initCommand := exec.Command("go", "mod", "init", name)

	initCommand.Dir = "/generated"
	initCommand.Stdin = os.Stdin
	initCommand.Stdout = os.Stdout

	err = initCommand.Run()
	if err != nil {
		panic("Go mod init failed")
	}

	err = os.MkdirAll("service", os.ModeDir)

	if err != nil {
		panic("Unable to create service dir")
	}

	err = os.MkdirAll("controller", os.ModeDir)
	if err != nil {
		panic("Unable to create controller dir")
	}

	err = os.MkdirAll("model", os.ModeDir)
	if err != nil {
		panic("Unable to create model dir")
	}

	err = os.MkdirAll("routes", os.ModeDir)
	if err != nil {
		panic("Unable to create model dir")
	}
	_, err = os.Create("routes/routes.go")
	if err != nil {
		panic("Unable to create routes file")
	}
}

func updateRoutesFile(structDefs []StructDefinition, database string) error {
	routesFilePath := "routes/routes.go"

	// Open the routes.go file for appending
	routesFile, err := os.OpenFile(routesFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer routesFile.Close()

	// Parse the routes template
	tmpl, err := template.New("routes").Parse(tmpl.RoutesTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	type structname struct {
		StructName string
	}
	data := []structname{}

	for _, structDef := range structDefs {
		data = append(data, structname{StructName: structDef.StructName})
	}

	// Execute the template and write to the routes file
	if err := tmpl.Execute(routesFile, data); err != nil {
		return err
	}

	return nil
}

func generateStructFromJSON(jsonData map[string]interface{}, structName string) (string, error) {
	// Initialize the struct code
	structCode := fmt.Sprintf("type %s struct {\n", structName)

	var structVar ModelType

	structVar.StructName = structName

	// Iterate through JSON fields and generate struct fields
	for fieldName, fieldValue := range jsonData {
		fieldType := inferGoType(fieldValue)
		structField := fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, fieldName)
		structVar.Fields = append(structVar.Fields, Field{
			Name: fieldName,
			Type: fieldType,
		})
		structCode += structField
	}

	// Define the file path for the model file
	filePath := fmt.Sprintf("model/%s.go", strings.ToLower(structName))

	// Create the model file
	modelFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer modelFile.Close()

	// Parse the model template
	tmpl, err := template.New("model").Parse(tmpl.ModelTemplate)
	if err != nil {
		return "", err
	}

	// Execute the template and write the struct code to the model file
	if err := tmpl.Execute(modelFile, structVar); err != nil {
		return "", err
	}

	// Close the struct definition
	structCode += "}\n"

	return structCode, nil
}

func inferGoType(value interface{}) string {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return "int64"
	case float32, float64:
		return "float64"
	case string:
		return "string"
	case bool:
		return "bool"
	case []interface{}:
		// If it's an array, infer the element type
		if len(value.([]interface{})) > 0 {
			return "[]" + inferGoType(value.([]interface{})[0])
		}
	}

	// Default to interface{} for unsupported types
	return "interface{}"
}

func generateServiceFile(fileName string, structDef StructDefinition, structCode, database string, appName string) error {

	serviceFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	// Parse the service template
	tmpl, err := template.New("service").Parse(tmpl.ServiceTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	data := struct {
		AppName    string
		StructName string
		Database   string
		StructCode string
	}{
		AppName:    appName,
		StructName: structDef.StructName,
		Database:   database,
		StructCode: structCode,
	}

	// Execute the template and write to the service file
	if err := tmpl.Execute(serviceFile, data); err != nil {
		return err
	}

	return nil
}

func generateControllerFile(fileName string, structDef StructDefinition, appName string) error {
	controllerFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer controllerFile.Close()

	// Parse the controller template
	tmpl, err := template.New("controller").Parse(tmpl.ControllerTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	data := struct {
		StructName string
		AppName    string
	}{
		StructName: structDef.StructName,
		AppName:    appName,
	}

	// Execute the template and write to the controller file
	if err := tmpl.Execute(controllerFile, data); err != nil {
		return err
	}

	return nil
}
