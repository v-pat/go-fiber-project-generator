package templates

const MainTemplate = `
package main

import (
	"fmt"
	"log"
	"{{.AppName}}/routes"
	"{{.AppName}}/databases"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	// Add other necessary imports
)

func main() {
	// Setup database connection
	databases.ConnectToDb()


	// Create a new Fiber app
	app := fiber.New()

	// Setup routes
	routes.Routes(app)

	// Start the server
	port := 8080 // Change this to your desired port
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}


`
