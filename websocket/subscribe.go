package websocket

import (
	"errors"
	"github.com/duxweb/go-fast/logger"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cast"
	"time"
)

// Sub 订阅频道
func (p *Client) Sub(channels ...string) error {
	if len(channels) < 1 {
		return errors.New("empty channel")
	}

	logger.Log("websocket").Debug().Str("client", p.clientID).Any("channels", channels).Msg("Client Sub")

	for _, topic := range channels {
		// 判断客户是否订阅该频道
		if !isExists(p.channel, topic) {
			p.channel.PushBack(topic)
		}

		// 判断频道是否存在
		if data, ok := Service.Channels.Load(topic); ok {
			item := data.(*Channel)
			// 判断频道是否包含客户
			if !isExists(item.channel, p.clientID) {
				item.channel.PushFront(p.clientID)
			}
		} else {
			// 创建新的频道映射
			clients := listLock{}
			clients.PushBack(p.clientID)
			Service.Channels.Store(topic, &Channel{
				clientID: topic,
				channel:  &clients,
			})
		}
	}
	return nil
}

// Unsub 取消订阅频道
func (p *Client) Unsub(args ...string) {

	var topics = make([]string, 0)
	if len(args) <= 0 {
		for e := p.channel.Front(); e != nil; e = e.Next() {
			topics = append(topics, cast.ToString(e.Value))
		}
	} else {
		topics = args
	}

	logger.Log("websocket").Debug().Str("client", p.clientID).Any("channels", topics).Msg("Client Unsub")

	for _, topic := range topics {
		// 判断频道是否存在
		data, ok := Service.Channels.Load(topic)
		if !ok {
			continue
		}

		// 从频道中移除客户端
		item := data.(*Channel)
		for e := item.channel.Front(); e != nil; e = e.Next() {
			if cast.ToString(e.Value) == p.clientID {
				item.channel.Remove(e)
			}
		}

		// 如果频道为空则删除频道
		if item.channel.Len() == 0 {
			Service.Channels.Delete(topic)
		}

		// 从客户端中移除频道
		for e := p.channel.Front(); e != nil; e = e.Next() {
			if cast.ToString(e.Value) == topic {
				p.channel.Remove(e)
			}
		}
	}
}

// Send 发布频道消息
func (p *Client) Send(data map[string]any) {
	err := SendClient(p.clientID, data)
	if err != nil {
		// 稍后重试发送
		_ = ants.Submit(func() {
			time.Sleep(3)
			_ = SendClient(p.clientID, data)
		})
	}
}

// Push 发布频道消息
func Push(channels []string, data map[string]any) {

	for _, topic := range channels {
		err := SendClient(topic, data)
		if err != nil {
			// 稍后重试发送
			_ = ants.Submit(func() {
				time.Sleep(3)
				_ = SendClient(topic, data)
			})
		}
	}
}

func isExists(list *listLock, value string) bool {
	for e := list.Front(); e != nil; e = e.Next() {
		if cast.ToString(e.Value) == value {
			return true
		}
	}
	return false
}
