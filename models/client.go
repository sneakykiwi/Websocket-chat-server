package models

import "github.com/gofiber/websocket/v2"

type Client struct {
	Connection *websocket.Conn
	ID int
}