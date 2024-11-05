package main

import (
	"fed/internal/server"
)

func main() {
	serv := server.New()
	serv.Go()
}

// func main() {
// 	app := fiber.New(fiber.Config{
// 		AppName:           "TestFiber",
// 		CaseSensitive:     false,
// 		EnablePrintRoutes: true,
// 		ErrorHandler:      fiber.DefaultErrorHandler,
// 		StrictRouting:     true,
// 	})
// 	// 	package main

// 	// import (
// 	// 	"fed/internal/server"
// 	// )

// 	// func main() {
// 	// 	serv := server.New()
// 	// 	serv.Go()
// 	// }

// 	// app.Get("/coffee", func(c *fiber.Ctx) error {
// 	// 	return c.Redirect("/teapot")
// 	// })

// 	// app.Get("/teapot", func(c *fiber.Ctx) error {
// 	// 	return c.Status(fiber.StatusTeapot).SendString("🍵 short and stout 🍵")
// 	// })

// 	// app.Get("/api/*", func(c *fiber.Ctx) error {
// 	// 	return c.SendString("API path: " + c.Params("*"))
// 	// 	// => API path: user/john
// 	// })
// 	// app.Get("/:value?", func(c *fiber.Ctx) error {
// 	// 	return c.SendString("value: " + c.Params("value"))
// 	// 	// => Get request with value: hello world
// 	// })

// 	// app.Use(func(c *fiber.Ctx) error {
// 	// 	return fiber.NewError(782, "Custom error message")
// 	// })
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		fmt.Println("1st route!")
// 		return c.Next()
// 	})

// 	app.Get("*", func(c *fiber.Ctx) error {
// 		fmt.Println("2nd route!")
// 		return c.Next()
// 	})

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		fmt.Println("3rd route!")
// 		return c.SendString("Hello, World!")
// 	})
// 	app.Listen(":8080")
// }
