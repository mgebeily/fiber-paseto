package pasetoware

import (
	"testing"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Fatalf("%s != %s", a, b)
	}
}

func Test_Query_Source(t *testing.T) {
	app := fiber.New()

	app.Use(New(&Config{
		TokenLookup:  "query:token",
		SymmetricKey: []byte("01234567890123456789012345678901"),
	}))

	app.Get("/:value", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("value"))
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/value", nil))

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 401)

	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/value?token=test", nil))

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 401)

	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/value?token=v2.local.R3VNW9sJWJlOLIGTIHkIGqA1QUz-4SkDiI6pZMGC0h8oGCqnAHC6ww6HiSdf-Nqh6kOGmdxdYslLDdRsfvcIoKl_B0H2GxL0HweD4CqpDISf91Xy68RkQwurP66DJ8GAGzyyIGE4wpsCf_GFvTk-mkb4r9kpfKnErUL8AvJCa1YvM5v9vP8jOfZ7rhQRtQ.c29tZSBmb290ZXI", nil))

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 200)
}

func Test_Header_Source(t *testing.T) {
	app := fiber.New()

	app.Use(New(&Config{
		SymmetricKey: []byte("01234567890123456789012345678901"),
	}))

	r := httptest.NewRequest(fiber.MethodGet, "/test", nil)

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("value"))
	})

	resp, err := app.Test(r)

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 401)

	r = httptest.NewRequest(fiber.MethodGet, "/test", nil)
	r.Header.Set("Authorization", "Bearer wrong")

	resp, err = app.Test(r)

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 401)

	r = httptest.NewRequest(fiber.MethodGet, "/test", nil)
	r.Header.Set("Authorization", "Bearer v2.local.R3VNW9sJWJlOLIGTIHkIGqA1QUz-4SkDiI6pZMGC0h8oGCqnAHC6ww6HiSdf-Nqh6kOGmdxdYslLDdRsfvcIoKl_B0H2GxL0HweD4CqpDISf91Xy68RkQwurP66DJ8GAGzyyIGE4wpsCf_GFvTk-mkb4r9kpfKnErUL8AvJCa1YvM5v9vP8jOfZ7rhQRtQ.c29tZSBmb290ZXI")

	resp, err = app.Test(r)

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 200)
}

func Test_Param_Source(t *testing.T) {
	app := fiber.New()

	app.Use(New(&Config{
		TokenLookup:  "param:token",
		SymmetricKey: []byte("01234567890123456789012345678901"),
	}))

	app.Get("/token/:token", func(c *fiber.Ctx) error {
		return c.SendString(c.Params("token"))
	})

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/token", nil))

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 401)

	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/token/test", nil))

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 401)

	resp, err = app.Test(httptest.NewRequest(fiber.MethodGet, "/token/v2.local.R3VNW9sJWJlOLIGTIHkIGqA1QUz-4SkDiI6pZMGC0h8oGCqnAHC6ww6HiSdf-Nqh6kOGmdxdYslLDdRsfvcIoKl_B0H2GxL0HweD4CqpDISf91Xy68RkQwurP66DJ8GAGzyyIGE4wpsCf_GFvTk-mkb4r9kpfKnErUL8AvJCa1YvM5v9vP8jOfZ7rhQRtQ.c29tZSBmb290ZXI", nil))

	assertEqual(t, err, nil)
	assertEqual(t, resp.StatusCode, 200)
}

func Test_Cookie_Source(t *testing.T) {

}

func Test_Multiple_Sources(t *testing.T) {

}

func Test_Filter_Pass(t *testing.T) {

}

func Test_Filter_Fail(t *testing.T) {

}
