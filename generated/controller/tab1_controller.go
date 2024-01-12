
package controller

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"vpat/model"
	"vpat/service"
)

// Createtab1 creates a new tab1.
func Createtab1(ctx *fiber.Ctx) error {
	var tab1 model.Tab1
	if err := ctx.BodyParser(&tab1); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := service.Createtab1(tab1)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tab1"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(tab1)
}

// Gettab1ByID retrieves a tab1 by ID.
func Gettab1ByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	intId,err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id not valid"})
	}
	tab1, err := service.Gettab1ByID(intId)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tab1 not found"})
	}

	return ctx.JSON(tab1)
}

// Updatetab1 updates an existing tab1 by ID.
func Updatetab1(ctx *fiber.Ctx) error {
	//id := ctx.Params("id")
	var updatedtab1 model.Tab1
	if err := ctx.BodyParser(&updatedtab1); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := service.Updatetab1(updatedtab1)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update tab1"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "tab1 updated"})
}

// Deletetab1ByID deletes a tab1 by ID.
func Deletetab1ByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	intId,err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id not valid"})
	}
	err = service.Deletetab1ByID(intId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete tab1"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "tab1 deleted"})
}

