package templates

const RoutesTemplate = `
package routes

import (
    "github.com/gofiber/fiber/v2"
    "AppName/controller"
)

// Setup sets up routes for the structs resource.
func Routes(app *fiber.App, 
	{{range .}}
	{{.StructName}}Controller *controller.{{.StructName}}Controller,
	{{end}}
	) {
	{{range .}}
    // Define routes for {{.StructName}} resource
    {{.StructName}}_group := app.Group("/{{.StructName}}")

    // Create a {{.StructName}}
    {{.StructName}}_group.Post("/", {{.StructName}}Controller.Create{{.StructName}})

    // Get a {{.StructName}} by ID
    {{.StructName}}_group.Get("/:id", {{.StructName}}Controller.Get{{.StructName}}ByID)

    // Update a {{.StructName}} by ID
    {{.StructName}}_group.Put("/:id", {{.StructName}}Controller.Update{{.StructName}})

    // Delete a {{.StructName}} by ID
    {{.StructName}}_group.Delete("/:id", {{.StructName}}Controller.Delete{{.StructName}}ByID)
	{{end}}
}
`
