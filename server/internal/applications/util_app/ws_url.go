package util_app

import (
	"encoding/json"
	"os"
	"sync"
)

type GameServer struct {
	Url  string `json:"url"`
	Load int    `json:"load"`
}

var gameServers []GameServer
var gsMu sync.Mutex

func Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&gameServers); err != nil {
		return err
	}

	return nil
}

func Fallback() {
	gameServers = append(gameServers, GameServer{
		Url:  "wss://vangowebsocket.vanshjangir.in",
		Load: 1,
	})
}

func Pick() string {
	gsMu.Lock()
	defer gsMu.Unlock()

	var selected *GameServer
	minLoad := int(^uint(0) >> 1)

	for i := range gameServers {
		if gameServers[i].Load < minLoad {
			selected = &gameServers[i]
			minLoad = gameServers[i].Load
		}
	}
	if selected != nil {
		selected.Load++
	}
	return selected.Url
}

func Release(url string) {
	gsMu.Lock()
	defer gsMu.Unlock()
	for i := range gameServers {
		if gameServers[i].Url == url {
			gameServers[i].Load -= 1
			break
		}
	}
}
