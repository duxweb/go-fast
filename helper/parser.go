package helper

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

func Qs(c *fiber.Ctx) (*gjson.Result, error) {
	paramsJson, err := json.Marshal(c.Queries())
	if err != nil {
		return nil, err
	}
	params := gjson.ParseBytes(paramsJson)
	return &params, nil
}

func Body(c *fiber.Ctx) (*gjson.Result, error) {
	payload := map[string]any{}
	err := c.BodyParser(&payload)
	if err != nil {
		return nil, err
	}
	paramsJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	params := gjson.ParseBytes(paramsJson)
	return &params, nil
}
