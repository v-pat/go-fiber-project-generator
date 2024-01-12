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
	db, err := databases.ConnectToDatabase()
	if err != nil {
		log.Fatal("Error setting up database:", err)
	}
	defer db.Close()


	// Create a new Fiber app
	app := fiber.New()

	// Setup routes
	routes.Routes(app)

	// Start the server
	port := 3000 // Change this to your desired port
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}


`
