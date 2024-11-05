package server

import (
	"errors"
	"fed/pkg/mysql"
	"fed/pkg/types"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *mysql.Users
}

type status struct {
	Code int16
	Text string
}

var (
	NotAllowed = status{
		Code: fiber.StatusMethodNotAllowed,
		Text: "Method Not Allowed",
	}
	NotFound = status{
		Code: fiber.StatusNotFound,
		Text: "Not Found",
	}
	InternalServerError = status{
		Code: fiber.StatusInternalServerError,
		Text: "Internal Server Error",
	}
	Ok = status{
		Code: fiber.StatusOK,
		Text: "OK",
	}
)

func (app *application) home() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		app.infoLog.Println("IP:", c.IP())
		return c.SendFile(regIndex)
	}
}
func (app application) Ok() func(*fiber.Ctx) error {
	return app.handleStatus(Ok, nil)
}
func (app application) NotAllowed() func(*fiber.Ctx) error {
	return app.handleStatus(NotAllowed, nil)
}

func (app application) NotFound() func(*fiber.Ctx) error {
	return app.handleStatus(NotFound, nil)
}
func (app application) ServerError() func(*fiber.Ctx) error {
	return app.handleStatus(InternalServerError, nil)
}

func (app *application) Code(stat status, c *fiber.Ctx) error {
	c.Status(int(stat.Code))

	// Рендерим шаблон с переданными данными
	err := c.Render("code.layout", stat) // Предполагается, что файл называется "code.tmpl".
	if err != nil {
		app.errorLog.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}

func (app application) handleStatus(stat status, methods []string) func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		app.infoLog.Println("IP", c.IP())
		app.infoLog.Println(stat.Code, stat.Text, c.OriginalURL())

		if methods != nil {
			c.Set("Allow", strings.Join(methods, ", "))
		}

		return app.Code(stat, c)
	}
}

func (app application) registration() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		app.infoLog.Println(c.OriginalURL())
		switch c.Method() {
		case fiber.MethodPost:
			ok, id := app.regUser(c)
			if ok {
				c.Status(int(Ok.Code))
				return c.JSON(struct {
					Id int `json:"id"`
				}{id})
			}
			c.Status(int(NotFound.Code))
			return c.JSON(fiber.Map{"error": "User registration failed"})
		case fiber.MethodGet:
			return app.showUser(c)
		default:
			return app.Code(NotAllowed, c)
		}
	}

}

func (app application) regUser(c *fiber.Ctx) (bool, int) {

	var user types.User

	if err := c.BodyParser(&user); err != nil {
		app.errorLog.Println(err)
		return false, 0
	}

	if user.Login == "" || user.Password == "" {
		app.infoLog.Println("Плохая отправка")
		return false, 0
	}
	app.infoLog.Println("Хорошая отправка")
	user.Created = time.Now()

	app.infoLog.Println(user)
	id, _ := app.db.Insert(&user)
	return true, id
}

func (app *application) showUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 1 {
		return app.Code(NotFound, c) // Страница не найдена.
	}

	u, err := app.db.Get(id)
	if err != nil {
		if errors.Is(err, types.ErrNoRecord) {
			return app.Code(NotFound, c)
		}
		app.errorLog.Println(err)
		return app.Code(InternalServerError, c)
	}

	return c.JSON(u) // Возвращаем пользователя в формате JSON
}
