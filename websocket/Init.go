package websocket

import (
	"encoding/json"
	"github.com/duxweb/go-fast/logger"
	"github.com/olahol/melody"
	"github.com/spf13/cast"
	"sync"
	"time"
)

var Service *ServiceT

func New() *ServiceT {
	return &ServiceT{
		Websocket: melody.New(),
		Clients:   &sync.Map{},
		Channels:  &sync.Map{},
		Agents:    map[string]*Agent{},
	}
}

func Init() {
	Service = New()
	Service.Run()
}

type Agent struct {
	auth    func(token string) (map[string]any, error)
	event   func(name string, client *Client) error
	message func(message *Message, client *Client) error
}

type ServiceT struct {
	Websocket *melody.Melody
	Clients   *sync.Map
	Channels  *sync.Map
	Agents    map[string]*Agent
}

type Message struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Client  string `json:"client,omitempty"`
	Channel string `json:"channel"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func (t *ServiceT) RegisterAgents(name string, agent *Agent) {
	t.Agents[name] = agent
}

func (t *ServiceT) Run() {
	// 连接处理
	t.Websocket.HandleConnect(func(session *melody.Session) {
		token := session.Request.Header.Get("token")
		app := session.Request.URL.Query().Get("app")
		agent, ok := t.Agents[app]
		if !ok {
			msg := "applications are not registered"
			logger.Log("websocket").Error().Msg(msg)
			_ = session.CloseWithMsg([]byte(msg))
		}
		data, err := agent.auth(token)
		if err != nil {
			logger.Log("websocket").Error().Msg(err.Error())
			_ = session.CloseWithMsg([]byte(err.Error()))
			return
		}
		clientID := data["client_id"].(string)

		// 把之前的客户端踢下线
		RemoveClient(clientID)

		// 创建新客户端
		AddClient(app, session, clientID, token, data)

		// 设置客户端在线
		err = EventOnline(clientID)
		if err != nil {
			logger.Log("websocket").Error().Err(err).Msg("Client Online")
			_ = session.CloseWithMsg([]byte(err.Error()))
			return
		}
	})

	// 销毁处理
	t.Websocket.HandleDisconnect(func(session *melody.Session) {
		str, ok := session.Get("clientID")
		if !ok {
			return
		}
		clientID := cast.ToString(str)
		logger.Log("websocket").Debug().Str("clientID", clientID).Msg("Client Disconnect")

		// 移除客户端
		RemoveConnClient(session)

		// 发送离线
		err := EventOffline(clientID)
		if err != nil {
			logger.Log("websocket").Error().Err(err).Msg("Client Online")
		}
	})

	// ping 处理
	t.Websocket.HandlePong(func(session *melody.Session) {
		str, ok := session.Get("clientID")
		if !ok {
			return
		}
		clientID := cast.ToString(str)
		logger.Log("websocket").Debug().Str("clientID", clientID).Msg("ping")
		_ = EventPing(clientID)
	})

	// 收到消息
	t.Websocket.HandleMessage(func(s *melody.Session, msg []byte) {
		name, _ := s.Get("clientID")
		clientID := cast.ToString(name)
		logger.Log("websocket").Debug().Str("clientID", clientID).Str("message", string(msg)).Msg("Ws Received")

		data := Message{}
		err := json.Unmarshal(msg, &data)
		if err != nil {
			logger.Log("websocket").Error().Str("clientID", clientID).Str("message", string(msg)).Err(err).Msg("Ws Received")
			return
		}
		client, err := GetClient(clientID)
		if err != nil {
			return
		}

		agent := Service.Agents[client.app]
		err = agent.message(&data, client)
		if err != nil {
			logger.Log("websocket").Error().Err(err).Msg("Ws Received")
			return
		}

		msgData := map[string]any{
			"type": "receive",
			"data": map[string]any{
				"id":   data.Id,
				"time": time.Now().Format("2006-01-02 15:04:05"),
			},
		}

		logger.Log("websocket").Debug().Any("data", msgData).Msg("ws send")
		client.Send(msgData)
	})

}
