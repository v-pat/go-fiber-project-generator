package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

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

type Tables struct {
	Tables []StructDefinition `json:"tables"`
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
		var tables Tables
		if err := c.BodyParser(&tables); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		structDefs = tables.Tables

		// Parse the "database" query parameter
		database := c.Query("database")
		if database != "postgres" && database != "mysql" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid database type"})
		}

		err := os.MkdirAll("service", os.ModeDir)

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

		for _, structDef := range structDefs {
			// Generate Go struct definition
			structCode, err := generateStructFromJSON(structDef.JSONExample, structDef.StructName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate struct"})
			}

			// Generate CRUD methods and save in service package
			serviceFileName := fmt.Sprintf("service/%s_service.go", strings.ToLower(structDef.StructName))
			if err := generateServiceFile(serviceFileName, structDef, structCode, database); err != nil {
				fmt.Sprintln(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate service file"})
			}

			// Generate controller methods and save in controller package
			controllerFileName := fmt.Sprintf("controller/%s_controller.go", strings.ToLower(structDef.StructName))
			if err := generateControllerFile(controllerFileName, structDef); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate controller file"})
			}

			// Update routes.go to define API endpoints
			err = updateRoutesFile(structDef, database)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate routes file"})
			}
		}

		return c.JSON(fiber.Map{"message": "CRUD code generated and organized successfully"})
	})

	// Start the Fiber server on port 3000
	app.Listen(":3000")
}

func updateRoutesFile(structDef StructDefinition, database string) error {
	routesFilePath := "routes/routes.go"

	// Open the routes.go file for appending
	routesFile, err := os.OpenFile(routesFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer routesFile.Close()

	// Parse the routes template
	tmpl, err := template.New("routes").Parse(routesTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	data := struct {
		StructName string
	}{
		StructName: structDef.StructName,
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
	tmpl, err := template.New("model").Parse(modelTemplate)
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

func generateServiceFile(fileName string, structDef StructDefinition, structCode, database string) error {

	serviceFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	// Parse the service template
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	data := struct {
		StructName string
		Database   string
		StructCode string
	}{
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

func generateControllerFile(fileName string, structDef StructDefinition) error {
	controllerFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer controllerFile.Close()

	// Parse the controller template
	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	data := struct {
		StructName string
	}{
		StructName: structDef.StructName,
	}

	// Execute the template and write to the controller file
	if err := tmpl.Execute(controllerFile, data); err != nil {
		return err
	}

	return nil
}

const routesTemplate = `
package routes

import (
    "github.com/gofiber/fiber/v2"
    "myapp/controller"
)

// Setup{{.StructName}}Routes sets up routes for the {{.StructName}} resource.
func Setup{{.StructName}}Routes(app *fiber.App, {{.StructName}}Controller *controller.{{.StructName}}Controller) {
    // Define routes for {{.StructName}} resource
    group := app.Group("/{{.StructName}}")

    // Create a {{.StructName}}
    group.Post("/", {{.StructName}}Controller.Create{{.StructName}})

    // Get a {{.StructName}} by ID
    group.Get("/:id", {{.StructName}}Controller.Get{{.StructName}}ByID)

    // Update a {{.StructName}} by ID
    group.Put("/:id", {{.StructName}}Controller.Update{{.StructName}})

    // Delete a {{.StructName}} by ID
    group.Delete("/:id", {{.StructName}}Controller.Delete{{.StructName}}ByID)
}
`

// Define a model template for generating the struct code in the model file.
const modelTemplate = `package model

// {{.StructName}} represents the {{.StructName}} struct.
type {{.StructName}} struct {
{{range .Fields}}
	{{.Name}} {{.Type}} ` + "`json:\"{{.Name}}\"`" +
	`{{end}}
}`

const controllerTemplate = `
package controller

import (
	"github.com/gofiber/fiber/v2"
	"myapp/model"
	"myapp/service"
)

// {{.StructName}}Controller handles requests for {{.StructName}}.
type {{.StructName}}Controller struct {
	Service *service.{{.StructName}}Service
}

// Create{{.StructName}} creates a new {{.StructName}}.
func (c *{{.StructName}}Controller) Create{{.StructName}}(ctx *fiber.Ctx) error {
	var {{.StructName}} model.{{.StructName}}
	if err := ctx.BodyParser(&{{.StructName}}); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := c.Service.Create{{.StructName}}(&{{.StructName}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create {{.StructName}}"})
	}

	return ctx.Status(fiber.StatusCreated).JSON({{.StructName}})
}

// Get{{.StructName}}ByID retrieves a {{.StructName}} by ID.
func (c *{{.StructName}}Controller) Get{{.StructName}}ByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	{{.StructName}}, err := c.Service.Get{{.StructName}}ByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "{{.StructName}} not found"})
	}

	return ctx.JSON({{.StructName}})
}

// Update{{.StructName}} updates an existing {{.StructName}} by ID.
func (c *{{.StructName}}Controller) Update{{.StructName}}(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var updated{{.StructName}} model.{{.StructName}}
	if err := ctx.BodyParser(&updated{{.StructName}}); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := c.Service.Update{{.StructName}}(id, &updated{{.StructName}})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update {{.StructName}}"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "{{.StructName}} updated"})
}

// Delete{{.StructName}}ByID deletes a {{.StructName}} by ID.
func (c *{{.StructName}}Controller) Delete{{.StructName}}ByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.Service.Delete{{.StructName}}ByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete {{.StructName}}"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "{{.StructName}} deleted"})
}

`

const serviceTemplate = `package service

import (
    "database/sql"
    "fmt"
    "myapp/model"
)

// {{.StructName}}Service represents the service for {{.StructName}}.
type {{.StructName}}Service struct {
    DB *sql.DB
}

{{.StructCode}}

// Create{{.StructName}} inserts a new {{.StructName}} record into the database.
func (s *{{.StructName}}Service) Create{{.StructName}}({{.StructName}} *model.{{.StructName}}) error {
    query := dbQuery("{{.Database}}", "insert", "{{.StructName}}s")
    // Execute the query to insert {{.StructName}} into the database
    fmt.Println("Executing query:", query)
    _, err := s.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Get{{.StructName}} retrieves a {{.StructName}} record from the database by ID.
func (s *{{.StructName}}Service) Get{{.StructName}}(id int) (*model.{{.StructName}}, error) {
    query := dbQuery("{{.Database}}", "select", "{{.StructName}}s")
    // Execute the query to retrieve {{.StructName}} from the database
    fmt.Println("Executing query:", query)
    // Implement query execution and scanning here
    {{.StructName}} := &model.{{.StructName}}{} // Replace with actual retrieval logic
    return {{.StructName}}, nil
}

// Update{{.StructName}} updates an existing {{.StructName}} record in the database.
func (s *{{.StructName}}Service) Update{{.StructName}}({{.StructName}} *model.{{.StructName}}) error {
    query := dbQuery("{{.Database}}", "update", "{{.StructName}}s")
    // Execute the query to update {{.StructName}} in the database
    fmt.Println("Executing query:", query)
    _, err := s.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

// Delete{{.StructName}} deletes a {{.StructName}} record from the database by ID.
func (s *{{.StructName}}Service) Delete{{.StructName}}(id int) error {
    query := dbQuery("{{.Database}}", "delete", "{{.StructName}}s")
    // Execute the query to delete {{.StructName}} from the database
    fmt.Println("Executing query:", query)
    _, err := s.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

`
