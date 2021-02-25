package db

import (
	"awesomeProject/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"math/rand"
	"sync"
	"time"
)

var (
	users_db []*models.User
	messages_db []* models.Message
	ws_clients []* models.Client
	mu sync.Mutex
)


func Connect() {
	users_db = make([]*models.User, 0)
	messages_db = make([]*models.Message, 0)
	ws_clients = make([]*models.Client, 0)
	fmt.Println("Connected with database.")
}

// USER FUNCTIONS


func InsertUser(user *models.User){
	mu.Lock()
	if indexOfUser(user.ID, users_db) != -1{
		fmt.Println("User already exists!")
		return
	}
	users_db = append(users_db, user)
	mu.Unlock()
	fmt.Println("User has been successfully added!")
}

func DeleteUser(user_id int){
	mu.Lock()

	i := indexOfUser(user_id, users_db)
	if i == -1 {
		fmt.Println("User not found!")
		return
	}
	copy(users_db[:], users_db[i+1:])
	users_db[len(users_db)-1] = nil // or the zero value of T
	users_db = users_db[:len(users_db)-1]
	mu.Unlock()

}

func UpdateUser(user_id int, newName string){
	mu.Lock()
	index := indexOfUser(user_id, users_db)
	if index == -1 {
		fmt.Println("User not found!")
		return
	}
	users_db[index].Name = newName
	mu.Unlock()
}

func GetUser(user_id int) (models.User, error){
	index := indexOfUser(user_id, users_db)
	if index == -1{
		fmt.Println("User not found!")
		return models.User{ID: 0}, fiber.ErrNotFound
	}

	usr := users_db[index]

	return *usr, nil

}

func appendMessageToUser(message *models.Message, user models.User){
	mu.Lock()
	index := indexOfUser(user.ID, users_db)
	fmt.Println(string(len(users_db)))
	messages := append(users_db[index].Messages, *message)
	users_db[index].Messages = messages
	mu.Unlock()
}


func indexOfUser(user_id int, data []*models.User) int {
	for k, v := range data {
		if user_id == v.ID {
			return k
		}
	}
	return -1    //not found.
}

// MESSAGE FUNCTIONS

func GetLastMessages() fiber.Map{
	var data fiber.Map
	var messages []fiber.Map
	for i := 1; i <= 5; i++{
		if len(messages) - i < 0{
			break
		}
		msg := messages_db[len(messages_db) - i]
		sender, _ := GetUser(msg.Sender)
		messages = append(messages, fiber.Map{"value": msg.Value, "author": sender.Name})
	}
	if len(messages) > 0 {
		data["messages"] = messages;
	}
	data = fiber.Map{"error": "messages_not_found"}
	return data
}

func InsertMessage(message *models.Message){
	now := time.Now()
	message.Timestamp = int(now.Unix())
	mu.Lock()
	messages_db = append(messages_db, message)
	mu.Unlock()

	user, _ := GetUser(message.Sender)

	appendMessageToUser(message, user)
}

func DeleteMessage(message_id int){
	mu.Lock()

	i := indexOfMessage(message_id, messages_db)
	if i == -1 {
		fmt.Println("Message not found!")
		return
	}
	copy(messages_db[:], messages_db[i+1:])
	messages_db[len(messages_db)-1] = nil // or the zero value of T
	messages_db = messages_db[:len(messages_db)-1]
	mu.Unlock()

}

func UpdateMessage(message_id int, newValue string){
	mu.Lock()
	index := indexOfMessage(message_id, messages_db)
	if index == -1 {
		fmt.Println("Message not found!")
		return
	}
	messages_db[index].Value = newValue
	mu.Unlock()
}

func indexOfMessage(message_id int, data []*models.Message) int {
	for k, v := range data {
		if message_id == v.ID {
			return k
		}
	}
	return -1    //not found.
}

//CLIENT FUNCTIONS
func createClientID() int{
	possible_id := rand.Intn(100)

	if len(ws_clients) == 0{
		return possible_id
	}

	for i := 0; i<len(ws_clients); i++{
		if ws_clients[i].ID == possible_id {
			possible_id = rand.Intn(100)
		}
	}
	return possible_id
}



//func GetClient()

func AppendClient(conn *websocket.Conn) *models.Client{
	client := new(models.Client)
	client.Connection = conn
	client.ID = createClientID()
	ws_clients = append(ws_clients, client)
	return client
}

func BroadcastMessage(client *models.Client, message *models.Message, user models.User){
	fmt.Printf("%d", len(ws_clients))
	for i := 0; i < len(ws_clients); i++{
		if ws_clients[i].ID != client.ID{
			_ = ws_clients[i].Connection.WriteJSON(fiber.Map{"sender": user.Name, "message": message.Value})
		}
	}
}