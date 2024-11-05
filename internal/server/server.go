package server

import (
	"fed/pkg/mysql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type server struct {
	mux *fiber.App
	app application
}

const (
	tmpl     = "ui/html"
	regIndex = "ui/html/index.html"
	assets   = "ui/assets"
)

func New() *server {

	app := application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	userDB := "user:password@/db?parseTime=true"
	db, err := mysql.OpenDB(userDB)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.infoLog.Println("db is connected")
	app.db = mysql.New(db)

	engine := html.New(tmpl, ".tmpl")
	engine.Reload(true)
	engine.Debug(true)

	mux := fiber.New(fiber.Config{
		AppName:           "TestFiber",
		CaseSensitive:     false,
		EnablePrintRoutes: true,
		ErrorHandler:      fiber.DefaultErrorHandler,
		StrictRouting:     true,
		Views:             engine,
	})
	mux.Static("/assets/", assets, fiber.Static{
		Compress: true,
		Index:    "405.html",
	})
	// mux := http.NewServeMux()

	// fileServer := http.FileServer(modFileSystem{
	// 	fs:  http.Dir("ui/assets"),
	// 	app: app,
	// })
	// mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	// //mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assets))))

	// mux.HandleFunc("/assets", app.NotAllowed())

	// mux.HandleFunc("/home", app.home())
	mux.Get("/home", app.home())

	// mux.HandleFunc("/home/account", app.Ok())
	mux.Get("/home/account", app.Ok())

	// mux.HandleFunc("/db", app.registration())
	mux.Get("/db", app.registration())
	mux.Post("/db", app.registration())

	// mux.HandleFunc("/", app.NotFound())
	mux.Get("/*", app.NotFound())

	return &server{
		mux: mux,
		app: app,
	}
}

// func (sv *server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
// 	sv.app.infoLog.Println(pattern)
// 	sv.mux.HandleFunc(pattern, handler)
// }

func (sv *server) Go() {
	sv.app.infoLog.Println("Server is listening...")

	sv.app.errorLog.Fatal(sv.mux.Listen(":8080"))
}
