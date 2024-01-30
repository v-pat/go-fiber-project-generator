package main

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"
	"vpat_codegen/model"
	"vpat_codegen/server"

	"github.com/gofiber/fiber/v2"
)

func main() {
	server.Serve()
	// CmdHandler()
}

func CmdHandler() model.Errors {
	// Check if the command matches "vpat gen config <fileName>"
	args := os.Args[1:]
	if len(args) != 4 || args[0] != "vpat" || args[1] != "gen" || args[2] != "config" {
		fmt.Println("Usage: vpat gen config <fileName>")
		return model.NewErr("Usage: vpat gen config <fileName>", fiber.StatusBadRequest)
	}

	// Get the file name from the command-line arguments
	fileName := args[3]

	// Check the file extension
	fileExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
	if fileExt != "json" && fileExt != "txt" {
		fmt.Println("Error: Unsupported file type. Only JSON or TXT files are supported.")
		return model.NewErr("Error: Unsupported file type. Only JSON or TXT files are supported.", fiber.StatusBadRequest)
	}

	// Read the content of the file
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return model.NewErr("Error reading file:"+err.Error(), fiber.StatusBadRequest)
	}

	// Initialize a variable to hold the parsed config
	var appJson model.AppJson

	// Parse JSON or TXT based on the file extension
	switch fileExt {
	case "json":
		err = json.Unmarshal(data, &appJson)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return model.NewErr("Error parsing JSON: "+err.Error(), fiber.StatusBadRequest)
		}
	case "txt":
		// Assuming the TXT file contains JSON data
		err = json.Unmarshal(data, &appJson)
		if err != nil {
			fmt.Println("Error parsing JSON from TXT:", err)
			return model.NewErr("Error parsing JSON from TXT: "+err.Error(), fiber.StatusBadRequest)
		}
	}

	//call this function
	//generator.Generate()
	return model.NewErr("", fiber.StatusOK)
}
