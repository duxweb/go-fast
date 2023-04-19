package logger

import (
	"bufio"
	"embed"
	"encoding/json"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rotisserie/eris"
	"net/http"
	"os"
)

//go:embed static/*
var template embed.FS

func InitDashboard() {

	global.App.Use("/dux/log", filesystem.New(filesystem.Config{
		Root:       http.FS(template),
		PathPrefix: "",
		Browse:     true,
	}))

	global.App.Get("/log", func(c *fiber.Ctx) error {
		return c.Render("logger/template/index", fiber.Map{})
	})

	global.App.Get("/log/data", func(c *fiber.Ctx) error {
		filePath := `./data/default/error.log`

		FileHandle, err := os.Open(filePath)
		if err != nil {
			return eris.New("log file not found")
		}
		defer FileHandle.Close()
		lineScanner := bufio.NewScanner(FileHandle)

		data := make([]map[string]any, 0)

		for lineScanner.Scan() {
			format := map[string]any{}
			_ = json.Unmarshal([]byte(lineScanner.Text()), &format)
			data = append(data, format)
		}

		return response.New(c).Send("ok", data)
	})
}
