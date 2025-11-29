package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/types"
)

type UserHandler struct {
	userService interfaces.UserService
	startTime   time.Time
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{userService: userService, startTime: time.Now()}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request model.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
			},
		})
	}

	user := &model.User{
		Username: request.Username,
		Email:    request.Email,
	}

	if err := h.userService.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create user",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusCreated,
			Message: "User created successfully",
		},
		Data: model.CreateUserResponse{
			ID: user.ID,
		},
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var request model.GetUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
			},
		})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get user",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.ApiResponse{
		Status: types.ResponseStatus{
			Code:    http.StatusOK,
			Message: "User retrieved successfully",
		},
		Data: model.GetUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var request model.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
			},
		})
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var request model.DeleteUserRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ApiResponse{
			Status: types.ResponseStatus{
				Code:    http.StatusBadRequest,
				Message: "Invalid request",
			},
		})
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.Group("/users")
	{
		router.POST("", h.CreateUser)
		router.GET("/:id", h.GetUser)
		router.PUT("/:id", h.UpdateUser)
		router.DELETE("/:id", h.DeleteUser)
	}
}
