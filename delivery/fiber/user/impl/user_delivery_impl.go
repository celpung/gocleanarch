package delivery_impl

import (
	"net/http"
	"strconv"

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
	const (
		defaultPage  = 1
		defaultLimit = 10
		maxLimit     = 100
	)

	// Parse & normalize query
	page, err := strconv.Atoi(c.Query("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 1 {
		page = defaultPage
	}
	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	// Call usecase
	users, total, err := d.UserUsecase.Read(uint(page), uint(limit))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user data",
			"error":   err.Error(),
		})
	}

	// Map to DTO
	res, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
	}

	// Ceil division for total pages
	totalPage := (total + int64(limit) - 1) / int64(limit)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Users fetched successfully",
		"data": fiber.Map{
			"users":        res,
			"count":        total,
			"current_page": page,
			"total_page":   totalPage,
		},
	})
}

func (d *UserDeliveryStruct) SearchUser(c *fiber.Ctx) error {
	const (
		defaultPage  = 1
		defaultLimit = 10
		maxLimit     = 100
	)

	pageStr := c.Query("page", "")
	limitStr := c.Query("limit", "")

	page := defaultPage
	limit := defaultLimit

	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err != nil || v < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid page parameter",
			})
		} else {
			page = v
		}
	}
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err != nil || v < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid limit parameter",
			})
		} else {
			limit = v
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	keyword := c.Query("q", "")

	users, total, err := d.UserUsecase.Search(uint(page), uint(limit), keyword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data not found",
		})
	}

	res, err := mapper.MapStructList[entity.User, dto.UserResponse](users)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to map response list",
			"error":   err.Error(),
		})
	}

	totalPage := (total + int64(limit) - 1) / int64(limit)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Users fetched successfully",
		"data": fiber.Map{
			"users":        res,
			"count":        total,
			"current_page": page,
			"total_page":   totalPage,
		},
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
