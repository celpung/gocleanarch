package user_delivery_implementation

import (
	"strconv"

	"github.com/celpung/gocleanarch/delivery/dto"
	"github.com/celpung/gocleanarch/delivery/fiber/user_delivery"
	"github.com/celpung/gocleanarch/domain/user/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserDeliveryStruct struct {
	UserUsecase usecase.UserUsecaseInterface
}

func (d *UserDeliveryStruct) Register(c *fiber.Ctx) error {
	// userID := c.Locals("userID")
	// email := c.Locals("email")
	// role := c.Locals("role")

	var req dto.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input data",
			"error":   err.Error(),
		})
	}

	userEntity := dto.UserCreateRequestDTO(&req)
	user, err := d.UserUsecase.Create(userEntity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Register success!",
		"user":    dto.UserResponseDTO(user),
	})
}

func (d *UserDeliveryStruct) Login(c *fiber.Ctx) error {
	var login dto.UserLoginRequest
	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to bind login data!",
			"error":   err.Error(),
		})
	}

	token, err := d.UserUsecase.Login(login.Email, login.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Login failed!",
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success fetch user data!",
		"user":    dto.UserResponseListDTO(users),
	})
}

func (d *UserDeliveryStruct) UpdateUser(c *fiber.Ctx) error {
	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to bind update data!",
			"error":   err.Error(),
		})
	}

	userEntity := dto.UserUpdateRequestDTO(&req)
	user, err := d.UserUsecase.Update(userEntity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user!",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully!",
		"user":    dto.UserResponseDTO(user),
	})
}

func (d *UserDeliveryStruct) DeleteUser(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
	}

	if err := d.UserUsecase.SoftDelete(uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func NewUserDelivery(usecase usecase.UserUsecaseInterface) user_delivery.UserDeliveryInterface {
	return &UserDeliveryStruct{
		UserUsecase: usecase,
	}
}
