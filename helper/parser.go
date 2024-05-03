package helper

import (
	"encoding/json"
	"github.com/derekstavis/go-qs"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

func Qs(c echo.Context) (*gjson.Result, error) {
	paramsMaps, err := qs.Unmarshal(c.QueryString())
	if err != nil {
		return nil, err
	}
	paramsJson, err := json.Marshal(paramsMaps)
	if err != nil {
		return nil, err
	}
	params := gjson.ParseBytes(paramsJson)
	return &params, nil
}

func Body(c echo.Context) (*gjson.Result, error) {
	payload := map[string]any{}
	err := (&echo.DefaultBinder{}).BindBody(c, &payload)
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
