package handlers

import (
	db "awesomeProject/database"
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func UserCreate(c *fiber.Ctx) error{
	var body []byte
	body = c.Body()
	var data *models.User
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.ID == 0{
		fmt.Println("No ID has been given, skipping request")
		if err := c.JSON(fiber.Map{"err": "please provide an id"}); err != nil{
			return err
		}
		return nil
	}


	db.InsertUser(data)
	if err := c.JSON(data); err != nil {
		return err
	}
	return nil
}

func UserDelete(c *fiber.Ctx) error {
	var body []byte
	body = c.Body()
	type Response struct {
		ID int `json:"id"`
	}
	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.ID == 0{
		fmt.Println("No ID has been given, skipping request")
		if err := c.JSON(fiber.Map{"err": "please provide an id"}); err != nil{
			return err}
		return nil
	}
	db.DeleteUser(data.ID)
	return nil
}

func UserGet(c *fiber.Ctx) error{
	body := c.Body()
	type Response struct{
		ID int `json:"id"`
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data.ID == 0{
		fmt.Println("No ID has been given, skipping request")
		if err := c.JSON(fiber.Map{"err": "please provide an id"}); err != nil{
			return err}
		return nil
	}

	user, err := db.GetUser(data.ID); if err != nil {
		if err := c.JSON(fiber.Map{"err": "User not found"}); err != nil{
			return err}
	}else {
		_ = c.JSON(user)
	}
	return nil

}