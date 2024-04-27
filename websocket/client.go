package websocket

import (
	"container/list"
	"encoding/json"
	"errors"
	"github.com/duxweb/go-fast/logger"
	"github.com/olahol/melody"
	"github.com/spf13/cast"
	"log/slog"
	"sync"
)

// Client 客户端映射
type Client struct {
	app      string          // 应用
	clientID string          // 客户端ID
	token    string          //  客户授权 token
	conn     *melody.Session // ws 连接
	channel  *listLock       // 频道列表
	data     map[string]any  // 附加数据 - 登录授权传递
}

// Channel 频道映射
type Channel struct {
	clientID string
	channel  *listLock
}

type listLock struct {
	list.List
	sync.Mutex
}

// GetClient 获取客户端
func GetClient(clientID string) (*Client, error) {
	lastClient, ok := Service.Clients.Load(clientID)
	if !ok {
		return nil, errors.New("client is not online")
	}
	return lastClient.(*Client), nil
}

// SendClient 给客户端发消息
func SendClient(clientID string, message map[string]any) error {
	client, err := GetClient(clientID)
	if err != nil {
		logger.Log("websocket").Error("Send Client", err, slog.Any("message", message))
		return err
	}
	logger.Log("websocket").Debug("Send Client", slog.String("client", clientID), slog.Any("message", message))
	content, _ := json.Marshal(message)
	err = client.conn.Write(content)
	if err != nil {
		return err
	}
	return nil
}

// AddClient 添加客户端
func AddClient(app string, client *melody.Session, clientID string, token string, data map[string]any) {
	client.Set("clientID", clientID)
	logger.Log("websocket").Debug("Add Client", slog.String("client", clientID))
	Service.Clients.Store(clientID, &Client{
		conn:     client,
		clientID: clientID,
		channel:  &listLock{},
		app:      app,
		token:    token,
		data:     data,
	})
}

// RemoveConnClient 移除客户端
func RemoveConnClient(conn *melody.Session) {
	clientID, ok := conn.Get("clientID")
	if !ok {
		return
	}
	RemoveClient(cast.ToString(clientID))

}

// RemoveClient 移除客户端
func RemoveClient(clientID string) {
	client, err := GetClient(clientID)
	if err != nil {
		return
	}
	_ = client.conn.CloseWithMsg([]byte("Client disconnected"))

	// 取消全部订阅
	client.Unsub()

	Service.Clients.Delete(clientID)
	logger.Log("websocket").Debug("Del Client", slog.String("client", clientID))
}
