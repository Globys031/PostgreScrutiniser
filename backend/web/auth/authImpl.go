package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TO DO: add a validator that would check if passed username is system username
type AuthImpl struct {
	Jwt      *JwtWrapper
	Logger   *utils.Logger
	Validate *validator.Validate
}

type AppUser struct {
	// Name of our application's main user
	Name string `json:"name" validate:"required,username"`
}

// Login as postgrescrutiniser
func (impl *AuthImpl) PostLogin(c *gin.Context) {
	// 1. Get request body data
	loginData := LoginRequest{}
	if err := c.BindJSON(&loginData); err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: "incorrect payload format",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorMsg)
		return
	}

	// 2. Validate username = postgrescrutiniser
	request := AppUser{Name: *loginData.Name}
	if err := impl.Validate.Struct(request); err != nil {
		err := fmt.Errorf("Username should be the name of our main application user")
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusBadRequest, &errorMsg)
		return
	}

	// 3. Check if password is correct
	// TO DO: add code here

	// 3. Generate JWT token
	token, err := impl.Jwt.GenerateToken(loginData, impl.Logger)
	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorMsg)
		return
	}

	// 4. Return token response (aka, login)
	tokenResponse := &LoginSuccessResponse{
		Token: &token,
	}
	jsonData, err := json.Marshal(tokenResponse)
	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}
	c.Data(http.StatusAccepted, "application/json", jsonData)
}
