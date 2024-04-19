package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

// Struct Movie
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Struct Director
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Slice Globals
var movies []Movie

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	movies = append(movies, Movie{ID: "1", Isbn: "335521", Title: "KongKing", Director: &Director{Firstname: "Uno", Lastname: "San"}})
	movies = append(movies, Movie{ID: "2", Isbn: "4416", Title: "Loyal of Cake", Director: &Director{Firstname: "Drow Ba", Lastname: "Koli"}})

	app.Post("/login", login)

	// Middleware ก่อนการเรียกใช้ Request | Method
	app.Use(checkMiddleware)

	// JWT Middleware
	/* app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})) */

	app.Get("/movies", getMovies)
	app.Get("/movies/:id", getMovie)
	app.Post("/movies", createMovie)
	app.Put("/movies/:id", updateMovie)
	app.Delete("/movies/:id", deleteMovie)

	app.Get("/test-html", testHTML)

	app.Listen(":8080")
}

func testHTML(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		//Key + Value
		"Title":  "Hello, World!",
		"Movies": "Lahnma",
	})
}

func checkMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	fmt.Printf(
		"URL = %s, Method = %s, Time = %s\n",
		c.OriginalURL(), c.Method(), start,
	)

	return c.Next()
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	// user, password ไม่ตรง Unauthorized
	if user.Email != memberUser.Email || user.Password != memberUser.Password {
		return fiber.ErrUnauthorized
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"email": user.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   t,
	})
}
