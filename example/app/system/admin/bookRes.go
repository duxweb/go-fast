package admin

import (
	"context"
	"example/app/system/models"
	"fmt"

	"github.com/duxweb/go-fast/v2/action"
	"github.com/duxweb/go-fast/v2/resp"
	"gorm.io/gorm"
)

// QueryParams 查询参数
type QueryParams struct {
	Name string `query:"name" doc:"名称"`
}

// Info 输出信息结构
type Info struct {
	Id      uint   `json:"id" doc:"记录ID"`
	Name    string `json:"name" doc:"名称"`
	Content string `json:"content" doc:"内容"`
}

// Data 输入数据结构
type Data struct {
	Name    string `json:"name" doc:"名称"`
	Content string `json:"content" doc:"内容"`
}

// BookMeta 自定义Meta类型（继承现有Meta）
type BookMeta struct {
	action.Meta        // 继承现有的Meta字段
	Summary     string `json:"summary" doc:"摘要信息"`
}

// BookRes @Resource()
func BookRes() {
	// 使用泛型创建资源实例，
	res := action.New[models.Book, Info, QueryParams, Data, BookMeta, resp.EmptyMeta]()

	// 设置路由配置 - 可以在最开始就设置
	res.SetRoute("admin", "system.book", "/books")

	// 定义模型
	res.SetModel(models.Book{})

	// 定义全局查询
	res.Query(func(tx *gorm.DB, ctx context.Context) *gorm.DB {
		tx = tx.Preload("Role").Preload("Dept")
		return tx
	})

	// 定义列表筛选
	res.Filter(func(tx *gorm.DB, params *QueryParams) *gorm.DB {
		if params.Name != "" {
			tx = tx.Where("name LIKE ?", "%"+params.Name+"%")
		}
		return tx
	})

	// 定义列表和详情输出
	res.Transform(func(item *models.Book, index int, ctx *action.TransformContext) Info {
		return Info{
			Id:      item.ID,
			Name:    item.Name,
			Content: item.Content,
		}
	})

	// 设置自定义列表元数据
	res.ListMeta(func(data []models.Book) BookMeta {
		return BookMeta{
			Summary: fmt.Sprintf("共有 %d 本书籍", len(data)),
		}
	})

	// 定义创建和编辑保存
	res.Format(func(model *models.Book, data *Data, ctx context.Context) (*models.Book, error) {
		return &models.Book{
			Name:    data.Name,
			Content: data.Content,
		}, nil
	})

	// 最后执行注册 - 一行搞定！
	res.Register()
}
