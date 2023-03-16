package util

import (
	"bytes"
	"fmt"
	"github.com/duxweb/go-fast/handlers"
	"github.com/go-resty/resty/v2"
	"github.com/xuri/excelize/v2"
	"time"
)

func ExcelImport(url string) ([][]string, error) {
	resp, err := resty.New().SetTimeout(10 * time.Second).R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, handlers.Error(resp.String())
	}
	reader := bytes.NewReader(resp.Body())
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	return rows, nil

}
