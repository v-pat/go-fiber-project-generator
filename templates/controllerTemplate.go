package templates

const ControllerTemplate = `
package controller

import (
	"github.com/gofiber/fiber/v2"
	"{{.AppName}}/model"
	"{{.AppName}}/service"
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
