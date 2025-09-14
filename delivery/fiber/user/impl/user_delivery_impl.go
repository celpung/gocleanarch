package delivery_impl

import (
	"github.com/celpung/gocleanarch/application/user/domain/entity"
	"github.com/celpung/gocleanarch/application/user/domain/usecase"
	"github.com/celpung/gocleanarch/delivery/dto"
	delivery "github.com/celpung/gocleanarch/delivery/fiber/user"
	"github.com/celpung/gocleanarch/infrastructure/mapper"
	"github.com/celpung/gocleanarch/infrastructure/validation"
	"github.com/gofiber/fiber/v2"
)

type UserDeliveryStruct struct {
	UserUsecase usecase.UserUsecase
}

func (d *UserDeliveryStruct) Register(c *fiber.Ctx) error {
	var req dto.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}
	if err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	var e entity.User
	if err := mapper.CopyTo(&req, &e); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map request",
			"error":   err.Error(),
		})
	}

	user, err := d.UserUsecase.Create(&e)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	var res dto.UserResponse
	if err := mapper.CopyTo(user, &res); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map response",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Register success",
		"user":    res,
	})
}

func (d *UserDeliveryStruct) Login(c *fiber.Ctx) error {
	var req dto.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid login data",
			"error":   err.Error(),
		})
	}
	if err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	token, err := d.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Login failed",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login success",
		"token":   token,
	})
}

func (d *UserDeliveryStruct) GetAllUserData(c *fiber.Ctx) error {
	users, err := d.UserUsecase.Read()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user data",
			"error":   err.Error(),
		})
	}

	resp, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success fetch user data",
		"users":   resp,
	})
}

func (d *UserDeliveryStruct) UpdateUser(c *fiber.Ctx) error {
	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid update data",
			"error":   err.Error(),
		})
	}
	if err := validation.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	var payload entity.UpdateUserPayload
	if err := mapper.CopyTo(&req, &payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map update payload",
			"error":   err.Error(),
		})
	}

	user, err := d.UserUsecase.Update(&payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	var resp dto.UserResponse
	if err := mapper.CopyTo(user, &resp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map response",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    resp,
	})
}

func (d *UserDeliveryStruct) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	if err := d.UserUsecase.SoftDelete(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func NewUserDelivery(usecase usecase.UserUsecase) delivery.UserDelivery {
	return &UserDeliveryStruct{UserUsecase: usecase}
}
