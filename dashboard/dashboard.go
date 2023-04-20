package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/response"
	"github.com/gofiber/fiber/v2"
	"github.com/icza/backscanner"
	"github.com/rotisserie/eris"
	"os"
)

func Init() {

	//global.App.Use("/dashboard/static", filesystem.New(filesystem.Config{
	//	Root:       http.FS(template),
	//	PathPrefix: "",
	//	Browse:     true,
	//}))

	//global.App.Get("/dashboard", func(c *fiber.Ctx) error {
	//	return c.Render("dashboard/template/index", fiber.Map{})
	//})

	global.App.Get("/dashboard/log", func(c *fiber.Ctx) error {
		filePath := `./data/logs/default.log`

		FileHandle, err := os.Open(filePath)
		if err != nil {
			return eris.New("log file not found")
		}
		fi, err := FileHandle.Stat()
		if err != nil {
			panic(err)
		}
		defer FileHandle.Close()

		last := c.QueryInt("line", int(fi.Size()))
		scanner := backscanner.New(FileHandle, last)

		fmt.Println(last)

		data := make([]map[string]any, 0)
		current := 0
		for {
			line, pos, err := scanner.Line()
			if err != nil {
				break
			}
			if line == "" {
				continue
			}
			format := map[string]any{}
			_ = json.Unmarshal([]byte(line), &format)
			data = append(data, format)
			current = pos
		}

		return response.New(c).Send("ok", fiber.Map{
			"current": current,
			"last":    last,
			"list":    data,
		})
	})
}
