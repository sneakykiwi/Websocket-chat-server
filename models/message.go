package models

type Message struct {
	Value string `json:"value"`
	Sender int `json:"sender_id"`
	ID int `json:"id"`
	Timestamp int
}
