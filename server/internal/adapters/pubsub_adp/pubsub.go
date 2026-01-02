package pubsub_adp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vanshjangir/vango/server/internal/domain"
)

type PubsubRepo struct {
	rdb *redis.Client
	ctx context.Context
	ps	*redis.PubSub
}

func (r *PubsubRepo) setupRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	useTLS := os.Getenv("REDIS_TLS") == "true"

	opts := &redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
	}

	if useTLS {
		opts.Username = "default"
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	r.rdb = redis.NewClient(opts)
	r.ctx = context.Background()

	_, err := r.rdb.Ping(r.ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v\n", err)
	}
	log.Println("Successfully connected to Redis!")
}

func NewPubsubRepo() *PubsubRepo {
	r := new(PubsubRepo)
	r.setupRedis()
	return r
}

func (r *PubsubRepo) GetGameFromUsername(
	username string,
) (*domain.GameDataForPlayer, error) {
	hashKey := "live_game"
	rawJsonString, err := r.rdb.HGet(r.ctx, hashKey, username).Result()
	if err != nil {
		return nil, fmt.Errorf("GetGameFromUsername: redis: %v", err)
	}

	var gameData domain.GameDataForPlayer
	err = json.Unmarshal([]byte(rawJsonString), &gameData)
	if err != nil {
		return nil, fmt.Errorf("GetGameFromUsername: json: %v", err)
	}

	return &gameData, nil
}

func (r *PubsubRepo) SetGameWithUsername(
	gameId int, blackName, whiteName string,
) error {
	hashKey := "live_game"
	jsonData, err := json.Marshal(domain.GameDataForPlayer{
		GameId:    gameId,
		BlackName: blackName,
		Whitename: whiteName,
		StartTime: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("SetGameWithUsername: %v", err)
	}

	err = r.rdb.HSet(r.ctx, hashKey, blackName, jsonData).Err()
	if err != nil {
		return fmt.Errorf("SetGameWithUsername HSet for black: %v", err)
	}

	err = r.rdb.HSet(r.ctx, hashKey, whiteName, jsonData).Err()
	if err != nil {
		return fmt.Errorf("SetGameWithUsername HSet for white: %v", err)
	}

	return nil
}

func (r *PubsubRepo) Send(game *domain.Game, msg any) error {
	var payload domain.SpectateServicePayload
	payload.GameData = msg
	payload.GameId = game.Id
	
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Send: Marshal: %v", err)
	}
	err = r.rdb.Publish(r.ctx, "spectator_game", string(data)).Err()
	if err != nil {
		return fmt.Errorf("Send: Publish: %v", err)
	}
	return nil
}

func (r *PubsubRepo) Subscribe(name string) error {
	r.ps = r.rdb.Subscribe(r.ctx, name)
	_, err := r.ps.Receive(r.ctx)
	if err != nil {
		return fmt.Errorf("Subscribe: %v", err)
	}
	return nil
}

func (r *PubsubRepo) Unsubscribe(name string) error {
	err := r.ps.Unsubscribe(r.ctx, name)
	if err != nil {
		return fmt.Errorf("Unsubscribe: %v", err)
	}
	return nil
}

func (r *PubsubRepo) Receive(msgChan chan string) {
    defer func() {
        if recover() != nil {
            log.Println("msgChan closed")
        }
    }()
	
	ch := r.ps.Channel()
	for msg := range ch {
		msgChan <- msg.Payload
	}
}
