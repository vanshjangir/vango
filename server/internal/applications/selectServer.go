package applications

import (
	"encoding/json"
	"os"
	"sync"
)

type GameServer struct {
    url      string
    load     int
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

func Pick() string {
    gsMu.Lock()
    defer gsMu.Unlock()

    var selected *GameServer
    minLoad := int(^uint(0) >> 1)

    for i := range gameServers {
        if gameServers[i].load < minLoad {
            selected = &gameServers[i]
            minLoad = gameServers[i].load
        }
    }
    if selected != nil {
        selected.load++
    }
    return selected.url
}

func Release(url string) {
    gsMu.Lock()
    defer gsMu.Unlock()
    for i := range gameServers {
        if gameServers[i].url == url {
            gameServers[i].load -= 1
            break
        }
    }
}
