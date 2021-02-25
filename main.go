package main

import (
	db "awesomeProject/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"awesomeProject/handlers"
)



func main() {
	app := fiber.New()
	
	db.Connect()

	ws := app.Group("/ws")

	app.Post("/create-user", handlers.UserCreate)
	app.Post("/delete-user", handlers.UserDelete)
	app.Get("/view-user", handlers.UserGet)
	app.Get("/get-last-5-messages", handlers.GetLast5Messages)
	ws.Get("/chat", websocket.New(handlers.SendMessage))

	log.Fatal(app.Listen(":3000"))
}