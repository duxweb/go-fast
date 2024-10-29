package helper

import (
	"github.com/duxweb/go-fast/cache"
	"github.com/duxweb/go-fast/database"
	"github.com/duxweb/go-fast/models"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

func VisitIncrement(ctx echo.Context, hasType string, hasID uint, driver string, path string) error {

	date := carbon.Now().ToDateStruct()
	path = lo.Ternary[string](path == "", path, ctx.Path())

	if strings.Contains(path, "/theme") || strings.Contains(path, "/manage") || strings.Contains(path, "/public") || strings.Contains(path, "/static") || strings.Contains(path, "/install") {
		return nil
	}

	ua := ctx.Request().UserAgent()

	uaParse, err := UaParser(ua)
	if err != nil {
		return err
	}
	browser := uaParse.UserAgent.ToString()
	ip := ctx.RealIP()

	visit := models.LogVisit{}
	database.Gorm().Model(models.LogVisit{}).FirstOrCreate(&visit, map[string]any{
		"has_type": hasType,
		"has_id":   hasID,
	})
	database.Gorm().Model(models.LogVisit{}).Where("id = ?", visit.ID).UpdateColumn("pv", gorm.Expr("pv + ?", 1))

	database.Gorm().Model(models.LogVisitData{}).Debug().Create(&models.LogVisitData{
		HasType: hasType,
		HasId:   hasID,
		Date:    date,
		Ip:      ip,
		Browser: browser,
		Driver:  driver,
	})

	visitData := models.LogVisitData{}
	database.Gorm().Model(models.LogVisitData{}).FirstOrCreate(&visitData, models.LogVisitData{
		HasType: hasType,
		HasId:   hasID,
		Date:    date,
		Ip:      ip,
		Browser: browser,
		Driver:  driver,
	})
	database.Gorm().Model(models.LogVisitData{}).Where("id = ?", visitData.ID).UpdateColumn("num", gorm.Expr("num + ?", 1))

	keys := []string{
		hasType,
		cast.ToString(hasID),
		driver,
		ip,
		browser,
	}

	key := strings.Join(keys, ".")
	_, err = cache.Injector().Get([]byte(key))

	if err != nil {
		seconds := carbon.NewCarbon().EndOfDay().DiffAbsInSeconds(carbon.NewCarbon())
		cache.Injector().Set([]byte(key), []byte("lock"), int(seconds))

		database.Gorm().Model(models.LogVisit{}).Where("id = ?", visit.ID).UpdateColumn("uv", gorm.Expr("uv + ?", 1))

		IpParse, _ := IpParser(ip)
		if visitData.Country == "" && IpParse != "" {
			address := strings.Split(IpParse, "|")
			if len(address) >= 4 {
				database.Gorm().Model(models.LogVisitData{}).Where("id = ?", visitData.ID).Updates(&models.LogVisitData{
					Country:  address[0],
					Province: address[2],
					City:     address[3],
				})
			}
		}
	}

	if uaParse.Device.Brand == "Spider" {
		visitSpider := models.LogVisitSpider{}
		database.Gorm().Model(&models.LogVisitSpider{}).FirstOrCreate(&visitSpider, models.LogVisitSpider{
			HasType: hasType,
			HasId:   hasID,
			Date:    date,
			Name:    uaParse.Device.Family,
			Path:    path,
		})
		database.Gorm().Model(models.LogVisitSpider{}).Where("id = ?", visitSpider.ID).UpdateColumn("num", gorm.Expr("num + ?", 1))
	}

	return nil

}
