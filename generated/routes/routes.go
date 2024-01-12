
package routes

import (
    "github.com/gofiber/fiber/v2"
    "vpat/controller"
)

// Setup sets up routes for the structs resource.
func Routes(app *fiber.App) {
	
    // Define routes for tab1 resource
    tab1_group := app.Group("/tab1")

    // Create a tab1
    tab1_group.Post("/", controller.Createtab1)

    // Get a tab1 by ID
    tab1_group.Get("/:id", controller.Gettab1ByID)

    // Update a tab1 by ID
    tab1_group.Put("/:id", controller.Updatetab1)

    // Delete a tab1 by ID
    tab1_group.Delete("/:id", controller.Deletetab1ByID)
	
    // Define routes for tab2 resource
    tab2_group := app.Group("/tab2")

    // Create a tab2
    tab2_group.Post("/", controller.Createtab2)

    // Get a tab2 by ID
    tab2_group.Get("/:id", controller.Gettab2ByID)

    // Update a tab2 by ID
    tab2_group.Put("/:id", controller.Updatetab2)

    // Delete a tab2 by ID
    tab2_group.Delete("/:id", controller.Deletetab2ByID)
	
}
