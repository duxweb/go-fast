package websocket

// EventOnline 客户端在线
func EventOnline(clientId string) error {
	client, err := GetClient(clientId)
	if err != nil {
		return err
	}

	agent := Service.Agents[client.app]
	err = agent.event("online", client)
	if err != nil {
		return err
	}

	return nil
}

// EventOffline 客户端离线
func EventOffline(clientId string) error {
	client, err := GetClient(clientId)
	if err != nil {
		return err
	}

	agent := Service.Agents[client.app]
	err = agent.event("offline", client)
	if err != nil {
		return err
	}
	return nil
}

// EventPing ping客户端
func EventPing(clientId string) error {
	client, err := GetClient(clientId)
	if err != nil {
		return err
	}

	agent := Service.Agents[client.app]
	err = agent.event("ping", client)
	if err != nil {
		return err
	}
	return nil
}
