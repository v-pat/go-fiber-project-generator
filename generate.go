package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	tmpl "vpat_codegen/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// DatabaseConnectionParams represents the parameters required for generating database connection code.
type DatabaseConnectionParams struct {
	DatabaseDriver     string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	DBURLFormat        string
	DatabaseDriverName string
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

		createFiles(appJson.AppName)
		err := createDatabase(database, cases.Title(language.English).String(appJson.AppName))
		if err != nil {
			panic("Unabel to create and connect to database")
		}

		for _, structDef := range structDefs {
			// Generate Go struct definition
			structCode, err := generateStructFromJSON(structDef.JSONExample, structDef.StructName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate struct"})
			}

			// Generate CRUD methods and save in service package
			serviceFileName := fmt.Sprintf("generated/service/%s_service.go", strings.ToLower(structDef.StructName))
			if err := generateServiceFile(serviceFileName, structDef, structCode, database, appJson.AppName); err != nil {
				fmt.Sprintln(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate service file"})
			}

			// Generate controller methods and save in controller package
			controllerFileName := fmt.Sprintf("generated/controller/%s_controller.go", strings.ToLower(structDef.StructName))
			if err := generateControllerFile(controllerFileName, structDef, appJson.AppName); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate controller file"})
			}

		}
		// Update routes.go to define API endpoints
		err = updateRoutesFile(structDefs, database, appJson.AppName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate routes file"})
		}

		err = createMainFile(structDefs, database, appJson.AppName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate main file"})
		}

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

func createDatabase(database string, appName string) error {
	params := DatabaseConnectionParams{
		DatabaseDriver:     "github.com/lib/pq",
		DBHost:             "localhost",
		DBPort:             "5432",
		DBName:             appName,
		DBUser:             "myuser",
		DBPassword:         "mypassword",
		DBURLFormat:        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DatabaseDriverName: "postgres",
	}

	if database == "postgres" {
		// Example usage: Generate code for PostgreSQL database connection and write to a file
		err := GenerateDatabaseConnectionCode(params, "postgres_connection.go")
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	} else if database == "mysql" {
		// Example usage: Generate code for MySQL database connection and write to a file
		params.DatabaseDriver = "github.com/go-sql-driver/mysql"
		params.DBPort = "3306"
		params.DatabaseDriverName = "mysql"
		params.DBUser = "root"
		params.DBPassword = "root"
		params.DBURLFormat = "%s:%s@tcp(%s:%s)/%s"
		params.DBName = appName

		err := GenerateDatabaseConnectionCode(params, "mysql_connection.go")
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

	}

	fmt.Println("Database connection code generated and written to files successfully.")

	return nil
}

// GenerateDatabaseConnectionCode generates code for connecting to a database (PostgreSQL or MySQL) and writes it to a file.
func GenerateDatabaseConnectionCode(params DatabaseConnectionParams, fileName string) error {
	// Create a new template
	tmpl, err := template.New("databaseConnection").Parse(tmpl.DatabaseConnectionTemplate)
	if err != nil {
		return err
	}

	// Create or open the file
	file, err := os.Create("generated/databases/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute the template and write the generated code to the file
	if err := tmpl.Execute(file, params); err != nil {
		return err
	}

	return nil
}

func createMainFile(structDefs []StructDefinition, database string, appName string) error {
	mainFilePath := "generated/main.go"

	// Open the main.go file for appending
	mainFile, err := os.OpenFile(mainFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer mainFile.Close()

	// Parse the routes template
	tmpl, err := template.New("main").Parse(tmpl.MainTemplate)
	if err != nil {
		return err
	}

	// Define data for the template
	type structname struct {
		StructName string
	}

	type dataType struct {
		StructNames []structname
		AppName     string
	}
	data := dataType{}
	names := []structname{}

	for _, structDef := range structDefs {
		names = append(names, structname{StructName: structDef.StructName})
	}

	data.StructNames = names
	data.AppName = appName

	// Execute the template and write to the routes file
	if err := tmpl.Execute(mainFile, data); err != nil {
		return err
	}

	return nil
}

func updateRoutesFile(structDefs []StructDefinition, database string, appName string) error {
	routesFilePath := "generated/routes/routes.go"

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

	type dataType struct {
		StructNames []structname
		AppName     string
	}
	data := dataType{}
	names := []structname{}

	for _, structDef := range structDefs {
		names = append(names, structname{StructName: structDef.StructName})
	}

	data.StructNames = names
	data.AppName = appName

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

	structVar.StructName = cases.Title(language.English).String(structName)

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
	filePath := fmt.Sprintf("generated/model/%s.go", strings.ToLower(structName))

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
		AppName             string
		StructName          string
		StructNameTitlecase string
		Database            string
		DBName              string
		StructCode          string
	}{
		AppName:             appName,
		StructName:          structDef.StructName,
		StructNameTitlecase: cases.Title(language.English).String(structDef.StructName),
		Database:            database,
		StructCode:          structCode,
		DBName:              cases.Title(language.English).String(appName),
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
		StructName          string
		AppName             string
		StructNameTitleCase string
	}{
		StructName:          structDef.StructName,
		AppName:             appName,
		StructNameTitleCase: cases.Title(language.English).String(structDef.StructName),
	}

	// Execute the template and write to the controller file
	if err := tmpl.Execute(controllerFile, data); err != nil {
		return err
	}

	return nil
}
