package websocket

import (
	"bytes"
	"encoding/json"
	"github.com/duxweb/go-fast/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/gookit/event"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

var Socket *Service

var (
	maxMsgSize = int64(512)
	pongWait   = 60 * time.Second
	newLine    = []byte{'\n'}
	space      = []byte{' '}
	pingPeriod = 10 * time.Second
)

func Init() {
	Socket = New()
	Socket.Start()
}

func Release() {
	Socket.Pool.Release()
}

type Service struct {
	Clients    map[string]Clients
	Users      map[string]Users
	Broadcast  chan *Broadcast
	Register   chan *Client
	Unregister chan *Client
	Pool       *ants.Pool
}

type Client struct {
	Auth      string
	Login     func(data string) (map[string]any, error)
	User      *User
	Socket    *websocket.Conn
	Mutex     sync.Mutex
	Send      chan *Message
	Message   []byte
	Service   *Service
	accountId string
}

type Users map[string]*User
type Clients map[*Client]bool

type User struct {
	ID     string
	Auth   string
	Client *Client
}

// Message is return msg
type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Broadcast struct {
	Client *Client
	Msg    []byte
}

func New() *Service {
	pool, _ := ants.NewPool(200000)
	return &Service{
		Pool: pool,
		//Socket:     socket,
		Clients:    map[string]Clients{},
		Users:      map[string]Users{},
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Handler message handling
func (r *Service) Handler(auth string, accountId string) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {

		return websocket.New(func(c *websocket.Conn) {
			// Set the client information
			var client Client
			client.Socket = c
			client.Service = r
			client.Send = make(chan *Message)
			client.Auth = auth
			client.accountId = accountId

			// registration client
			r.Register <- &client

			_ = r.Pool.Submit(func() {
				client.ServiceRead()
			})
			_ = r.Pool.Submit(func() {
				client.ServiceWrite()
			})

		})(c)
	}
}

// ServiceRead Get the client message
func (c *Client) ServiceRead() {
	defer func() {
		c.Service.Unregister <- c
		_ = c.Socket.Close()
	}()
	// Set the maximum value
	c.Socket.SetReadLimit(maxMsgSize)
	// Set link timeout
	_ = c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(appData string) error {
		// Every time pong is received, the deadline is pushed back 60 seconds
		_ = c.Socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			// Turn off error handling
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				logger.Log().Debug().Err(err).Msg("websocket error")
			}
			break
		}
		// Read received message
		message := bytes.TrimSpace(bytes.Replace(msg, newLine, space, -1))
		c.Service.Broadcast <- &Broadcast{
			Client: c,
			Msg:    message,
		}
	}
}

// ServiceWrite Writes the client message
func (c *Client) ServiceWrite() {
	// 写入超时
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Socket.Close()
	}()
	for {
		select {
		// Write a message to the current websocket connection
		case message, ok := <-c.Send:
			_ = c.Socket.SetWriteDeadline(time.Now().Add(pingPeriod))
			if !ok {
				// Closed channel
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// NextWriter Returns a writer for the next message to be sent. The writer's Close method flushes the complete message to the network.
			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			content, _ := json.Marshal(message)
			w.Write(content)
			// Adds the queued chat message to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				msg := <-c.Send
				content, _ = json.Marshal(msg)
				w.Write(newLine)
				w.Write(content)
			}
			if err := w.Close(); err != nil {
				return
			}
		//Heartbeat maintenance
		case <-ticker.C:
			_ = c.Socket.SetWriteDeadline(time.Now().Add(pingPeriod))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Start 启动服务
func (r *Service) Start() {

	_ = r.Pool.Submit(func() {
		for {
			select {
			case client := <-r.Register:
				// register channel
				if r.Clients[client.Auth] == nil {
					r.Clients[client.Auth] = map[*Client]bool{}
				}
				r.Clients[client.Auth][client] = true
				client.SendMsg("connect", "successful connection to socket service")

				// login client
				if r.Users[client.Auth] == nil {
					r.Users[client.Auth] = map[string]*User{}
				}
				user := &User{
					ID:     client.accountId,
					Client: client,
				}
				r.Users[client.Auth][client.accountId] = user
				client.User = user
				client.SendMsg("login", "login successful")

				// Notify user to go online
				r.Pool.Submit(func() {
					event.Fire("websocket.online", event.M{
						"client": client,
					})
				})

			case client := <-r.Unregister:
				// Notify user offline
				r.Pool.Submit(func() {
					event.Fire("websocket.offline", event.M{
						"client": client,
					})
				})
				// logout channel
				client.Close()
			case data := <-r.Broadcast:
				// broadcast channel
				MessageStruct := Message{}
				err := json.Unmarshal(data.Msg, &MessageStruct)
				if err != nil {
					data.Client.SendMsg("err", "incorrect message format")
					continue
				}
				if MessageStruct.Type == "" {
					data.Client.SendMsg("err", "message type error")
					continue
				}
				switch MessageStruct.Type {
				case "ping":
					data.Client.SendMsg("pong", "")
				default:
					_ = r.Pool.Submit(func() {
						event.Fire("websocket."+MessageStruct.Type, event.M{
							"client":  data.Client,
							"message": &MessageStruct,
						})
					})
				}
			}
		}

	})
}

// Close Closed channel
func (c *Client) Close() {
	close(c.Send)
	if _, ok := c.Service.Clients[c.Auth]; !ok {
		return
	}
	if _, ok := c.Service.Clients[c.Auth][c]; !ok {
		return
	}
	delete(c.Service.Clients[c.Auth], c)
	if _, ok := c.Service.Users[c.Auth]; !ok {
		return
	}
	if _, ok := c.Service.Users[c.Auth][c.User.ID]; !ok {
		return
	}
	delete(c.Service.Users[c.Auth], c.User.ID)
}

// SendMsg send a message
func (c *Client) SendMsg(Type string, message string, datas ...any) bool {
	var data any
	if len(datas) > 0 {
		data = datas[0]
	}
	select {
	case _, ok := <-c.Send:
		if !ok {
			return false
		}
	default:
		c.Send <- &Message{
			Type:    Type,
			Message: message,
			Data:    data,
		}
	}
	return true
}

func (c *Client) SendUserMsg(accountId string, Type string, message string, datas ...any) bool {
	user, ok := c.Service.Users[c.Auth][accountId]
	if !ok {
		return false
	}
	user.Client.SendMsg(Type, message, datas)
	return true
}

// GetUser Get user information based on id
func (c *Client) GetUser(accountId string) *User {
	if user, ok := c.Service.Users[c.Auth][accountId]; ok {
		return user
	}
	return nil
}

// Event  registered events
func Event(name string, call func(client *Client, message *Message) error) {
	event.On("websocket."+name, event.ListenerFunc(func(e event.Event) error {
		client := e.Get("client").(*Client)
		message := e.Get("message").(*Message)
		return call(client, message)
	}))
}
