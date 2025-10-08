package domain

import "time"

type User struct {
	Username      string
	Password      string
	Email         string
	Rating        int
	Highestrating int
}

type Player struct {
	User   User
	Color  int
	GameId *string
}

type WaitingPlayer struct {
	Id       string
	Username string
	Rating   int
	AddedAt  time.Time
}

type GameDataForPlayer struct {
	GameId	string `json:"gameid"`
	BlackName string `json:"blackname"`
	Whitename string `json:"whitename"`
	StartTime	int64 `json:"starttime"`
}
