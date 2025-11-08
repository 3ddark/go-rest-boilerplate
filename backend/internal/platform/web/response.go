package web

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success  bool          `json:"success"`
	Status   string        `json:"status"`
	Messages []string      `json:"messages,omitempty"`
	Data     any           `json:"data,omitempty"`
	Error    *ErrorDetails `json:"error,omitempty"`
}

type ErrorDetails struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, data interface{}, messages ...string) error {
	if len(messages) == 0 {
		messages = []string{"Request processed successfully."}
	}
	return c.Status(statusCode).JSON(Response{
		Success:  true,
		Status:   "success",
		Messages: messages,
		Data:     data,
	})
}

func Error(c *fiber.Ctx, statusCode int, errCode, errMsg string, details interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Status:  "error",
		Error: &ErrorDetails{
			Code:    errCode,
			Message: errMsg,
			Details: details,
		},
	})
}

func ValidationError(c *fiber.Ctx, err error) error {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, err.Field()+" is "+err.Tag())
	}
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success:  false,
		Status:   "validation_error",
		Messages: errors,
	})
}

func CustomError(c *fiber.Ctx, statusCode int, messages ...string) error {
	if len(messages) == 0 {
		messages = []string{"An unexpected error occurred."}
	}
	return c.Status(statusCode).JSON(Response{
		Success:  false,
		Status:   "error",
		Messages: messages,
	})
}

func Unauthorized(c *fiber.Ctx, messages ...string) error {
	if len(messages) == 0 {
		messages = []string{"Unauthorized."}
	}
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success:  false,
		Status:   "unauthorized",
		Messages: messages,
	})
}

func Forbidden(c *fiber.Ctx, messages ...string) error {
	if len(messages) == 0 {
		messages = []string{"Forbidden."}
	}
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success:  false,
		Status:   "forbidden",
		Messages: messages,
	})
}

func NotFound(c *fiber.Ctx, messages ...string) error {
	if len(messages) == 0 {
		messages = []string{"Resource not found."}
	}
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success:  false,
		Status:   "not_found",
		Messages: messages,
	})
}

func WithMessage(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success:  statusCode >= 200 && statusCode < 300,
		Status:   "message",
		Messages: []string{message},
	})
}

func Respond(c *fiber.Ctx, statusCode int, data interface{}, messages ...string) error {
	success := statusCode >= 200 && statusCode < 300
	status := "success"
	if !success {
		status = "error"
	}

	// If data is an error, handle it as an error response
	if err, ok := data.(error); ok {
		return CustomError(c, statusCode, err.Error())
	}

	// If data is a validation error, handle it specifically
	if _, ok := data.(validator.ValidationErrors); ok {
		return ValidationError(c, data.(error))
	}

	return c.Status(statusCode).JSON(Response{
		Success:  success,
		Status:   status,
		Messages: messages,
		Data:     data,
	})
}

func ValidationErrors(c *fiber.Ctx, errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		// Customize this part to create more user-friendly messages
		messages = append(messages, "Field '"+err.Field()+"' failed on the '"+err.Tag()+"' tag.")
	}
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success:  false,
		Status:   "validation_error",
		Messages: messages,
	})
}

func FromFiberError(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var messages []string

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		messages = append(messages, e.Message)
	} else {
		messages = append(messages, "Internal Server Error")
	}

	return c.Status(code).JSON(Response{
		Success:  false,
		Status:   "error",
		Messages: messages,
	})
}

func ParseBody(c *fiber.Ctx, body interface{}) error {
	if err := c.BodyParser(body); err != nil {
		return FromFiberError(c, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}
	return nil
}

func ValidateStruct(c *fiber.Ctx, s interface{}) error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return ValidationErrors(c, errs)
		}
		return FromFiberError(c, fiber.NewError(fiber.StatusInternalServerError, "Error validating request"))
	}
	return nil
}

func ParseAndValidate(c *fiber.Ctx, body interface{}) error {
	if err := ParseBody(c, body); err != nil {
		return err
	}
	if err := ValidateStruct(c, body); err != nil {
		return err
	}
	return nil
}

func GetQuery(c *fiber.Ctx, key string, defaultValue ...string) string {
	val := c.Query(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func GetParams(c *fiber.Ctx, key string, defaultValue ...string) string {
	val := c.Params(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func GetHeader(c *fiber.Ctx, key string, defaultValue ...string) string {
	val := c.Get(key)
	if val == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return val
}

func GetLocals(c *fiber.Ctx, key string) interface{} {
	return c.Locals(key)
}

func SetLocals(c *fiber.Ctx, key string, value interface{}) {
	c.Locals(key, value)
}

func GetQueryInt(c *fiber.Ctx, key string, defaultValue ...int) (int, error) {
	valStr := GetQuery(c, key)
	if valStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, fiber.NewError(fiber.StatusBadRequest, "Missing required query parameter: "+key)
	}
	val := c.QueryInt(key)
	if val == 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Invalid value for query parameter: "+key)
	}
	return val, nil
}

func GetParamsInt(c *fiber.Ctx, key string, defaultValue ...int) (int, error) {
	valStr := GetParams(c, key)
	if valStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return 0, fiber.NewError(fiber.StatusBadRequest, "Missing required path parameter: "+key)
	}
	val, err := c.ParamsInt(key)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, "Invalid value for path parameter: "+key)
	}
	return val, nil
}

func GetQueryBool(c *fiber.Ctx, key string, defaultValue ...bool) (bool, error) {
	valStr := GetQuery(c, key)
	if valStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return false, nil // Typically, missing boolean is false
	}
	return strings.ToLower(valStr) == "true", nil
}
