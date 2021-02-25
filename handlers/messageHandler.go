package handlers

import (
	db "awesomeProject/database"
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
)

func GetLast5Messages(c *fiber.Ctx) error{
	last_5_messages := db.GetLastMessages()
	err := c.JSON(last_5_messages); if err != nil{
		fmt.Print(err)
	}
	return nil

}

func SendMessage(c *websocket.Conn) {
	var (
		msg []byte
		err error
	)
	client := db.AppendClient(c)
	for {
		if _, msg, err = c.ReadMessage(); err != nil{
			log.Println(err)
			break
		}

		var message *models.Message
		_ = json.Unmarshal(msg, &message)
		user, _ := db.GetUser(message.Sender)

		db.InsertMessage(message)
		fmt.Printf("[%+v] ", user.Name)
		fmt.Print(message.Value, "\n")

		db.BroadcastMessage(client, message, user)

	}
}

