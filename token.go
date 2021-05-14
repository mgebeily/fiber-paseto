package pasetoware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/o1egl/paseto"
)

type Config struct {
	Filter func(*fiber.Ctx) bool

	// Encryption key
	SymmetricKey []byte

	SuccessHandler fiber.Handler

	ErrorHandler fiber.ErrorHandler

	TokenLookup string
	AuthScheme  string
	ContextKey  string
}

func New(config *Config) fiber.Handler {
	if config.SuccessHandler == nil {
		config.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if err.Error() == "TODO: What is the error here" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed JWT")
			} else {
				return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
			}
		}
	}

	if config.TokenLookup == "" {
		config.TokenLookup = "header:" + fiber.HeaderAuthorization
	}

	if config.ContextKey == "" {
		config.ContextKey = "user"
	}

	if config.AuthScheme == "" {
		config.AuthScheme = "Bearer"
	}

	getToken := createTokenExtractor(config.TokenLookup, config.AuthScheme)
	decryptor := paseto.NewV2()

	return func(c *fiber.Ctx) error {
		if config.Filter != nil && config.Filter(c) {
			return c.Next()
		}

		var result paseto.JSONToken
		var footer string
		token := getToken(c)

		err := decryptor.Decrypt(token, config.SymmetricKey, &result, &footer)

		if err != nil {
			return config.ErrorHandler(c, err)
		}

		c.Locals(config.ContextKey, result)
		return config.SuccessHandler(c)
	}
}

func createTokenExtractor(tokenLookup string, authScheme string) func(c *fiber.Ctx) string {
	sources := strings.Split(tokenLookup, ",")
	checks := make([]func(c *fiber.Ctx) string, 0)
	authScheme = authScheme + " "

	for _, source := range sources {
		parts := strings.Split(source, ":")

		switch parts[0] {
		case "header":
			checks = append(checks, func(c *fiber.Ctx) string {
				values := strings.Split(c.Get(parts[1]), authScheme)

				if len(values) == 2 {
					return values[1]
				} else {
					return ""
				}
			})
		case "query":
			checks = append(checks, func(c *fiber.Ctx) string {
				return c.Query(parts[1])
			})
		case "param":
			checks = append(checks, func(c *fiber.Ctx) string {
				return c.Params(parts[1])
			})
		case "cookie":
			checks = append(checks, func(c *fiber.Ctx) string {
				return c.Cookies(parts[1])
			})
		}
	}

	return func(c *fiber.Ctx) string {
		for _, check := range checks {
			v := check(c)
			if v != "" {
				return v
			}
		}

		return ""
	}
}
