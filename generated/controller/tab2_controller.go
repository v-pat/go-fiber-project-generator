
package controller

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"vpat/model"
	"vpat/service"
)

// Createtab2 creates a new tab2.
func Createtab2(ctx *fiber.Ctx) error {
	var tab2 model.Tab2
	if err := ctx.BodyParser(&tab2); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := service.Createtab2(tab2)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tab2"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(tab2)
}

// Gettab2ByID retrieves a tab2 by ID.
func Gettab2ByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	intId,err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id not valid"})
	}
	tab2, err := service.Gettab2ByID(intId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tab2 not found"})
	}

	return ctx.JSON(tab2)
}

// Updatetab2 updates an existing tab2 by ID.
func Updatetab2(ctx *fiber.Ctx) error {
	//id := ctx.Params("id")
	var updatedtab2 model.Tab2
	if err := ctx.BodyParser(&updatedtab2); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := service.Updatetab2(updatedtab2)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update tab2"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "tab2 updated"})
}

// Deletetab2ByID deletes a tab2 by ID.
func Deletetab2ByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	intId,err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id not valid"})
	}
	err = service.Deletetab2ByID(intId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete tab2"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "tab2 deleted"})
}

