package generators

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"vpat_codegen/model"

	"archive/zip"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Generate(appJson model.AppJson, dirPath string) (string, model.Errors) {

	err := GenerateApplicationCode(appJson, appJson.Database, appJson.Language, dirPath)

	if err != nil {
		fmt.Println("Unable to generate application code : " + err.Error())
		return "", model.NewErr("Unable to generate application code : "+err.Error(), fiber.StatusInternalServerError)
	}

	zipFile, err := CreateApplicationZip(appJson.AppName)

	if err != nil {
		fmt.Println("Unable to zip application code  : " + err.Error())
		return "", model.NewErr("Unable to zip application code  : "+err.Error(), fiber.StatusInternalServerError)
	}

	err = os.RemoveAll("./generated")
	if err != nil {
		fmt.Println("Unable to clean generated directory  : " + err.Error())
		return "", model.NewErr("Unable to clean generated directory  : "+err.Error(), fiber.StatusInternalServerError)
	}

	return zipFile, model.NewErr("Code Generated Successfull.", fiber.StatusOK)
}

func createFiles(name string, dirPath string) {
	err := os.MkdirAll(dirPath, os.ModeDir)

	if err != nil {
		panic("Unable to create generated dir")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	initCommand := exec.Command("go", "mod", "init", name)

	initCommand.Dir = dirPath + "/"
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	importCommand := exec.Command("goimports", "-l", "-w", ".")

	importCommand.Dir = "./generated/"
	importCommand.Stdin = os.Stdin
	importCommand.Stdout = &stdout
	importCommand.Stderr = &stderr

	err := importCommand.Run()
	if err != nil {
		fmt.Println("Goimports failed:" + err.Error())
		return err
	}

	getCommand := exec.Command("go", "get", "-u")

	getCommand.Dir = "./generated/"
	getCommand.Stdin = os.Stdin
	getCommand.Stdout = &stdout
	getCommand.Stderr = &stderr

	err = getCommand.Run()
	if err != nil {
		fmt.Println("go get failed:" + err.Error())
		return err
	}

	return nil
}

func CreateServices(structDefs []model.StructDefinition, database string, appName string) (fiber.Map, error) {
	for _, structDef := range structDefs {
		// Generate Go struct definition
		structCode, err := GenerateStructFromJSON(structDef.JSONExample, structDef.StructName)
		if err != nil {
			return fiber.Map{"error": "Failed to generate struct"}, err
		}

		// Generate CRUD methods and save in service package
		serviceFileName := fmt.Sprintf("generated/service/%s_service.go", strings.ToLower(structDef.StructName))
		if err := GenerateServiceFile(serviceFileName, structDef, structCode, database, appName); err != nil {
			fmt.Sprintln(err.Error())
			return fiber.Map{"error": "Failed to generate service file"}, err
		}

		// Generate controller methods and save in controller package
		controllerFileName := fmt.Sprintf("generated/controller/%s_controller.go", strings.ToLower(structDef.StructName))
		if err := GenerateControllerFile(controllerFileName, structDef, appName); err != nil {
			return fiber.Map{"error": "Failed to generate controller file"}, err
		}

	}

	return nil, nil
}

func GenerateApplicationCode(appJson model.AppJson, database string, lang string, dirPath string) error {
	// Parse the incoming JSON data as a StructDefinition
	var structDefs []model.StructDefinition

	structDefs = appJson.Tables

	// Parse the "database" query parameter
	if database != "postgres" && database != "mysql" {
		fmt.Println("Unabel to process request : Invalid database type")
		return errors.New("Only supported databases are mysql and postgres.")
	}

	createFiles(appJson.AppName, dirPath)

	err := CreateDatabase(database, cases.Title(language.English).String(appJson.AppName), structDefs, appJson.AppName)
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
	err = UpdateRoutesFile(structDefs, database, appJson.AppName)
	if err != nil {
		fmt.Println("Unabel to generate routes : " + err.Error())
		return err
	}

	err = CreateMainFile(structDefs, database, appJson.AppName)
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
