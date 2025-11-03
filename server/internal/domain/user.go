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
	GameId *int
}

type WaitingPlayer struct {
	Id       string 	`json:"id"`
	Username string		`json:"username"`
	Rating   int		`json:"rating"`
	AddedAt  time.Time	`json:"addedat"`
}

type GameDataForPlayer struct {
	GameId	int `json:"gameid"`
	BlackName string `json:"blackname"`
	Whitename string `json:"whitename"`
	StartTime	time.Time `json:"starttime"`
}
