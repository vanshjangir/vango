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
	Username string    `json:"username"`
	Rating   int       `json:"rating"`
	AddedAt  time.Time `json:"addedat"`
}

type GameDataForPlayer struct {
	GameId    int       `json:"gameid"`
	BlackName string    `json:"blackname"`
	Whitename string    `json:"whitename"`
	StartTime time.Time `json:"starttime"`
}

type SpectateServicePayload struct {
	GameId int `json:"gameid"`
	GameData any `json:"gamedata"`
}

type ClientRecvResult struct {
	data []byte
	err  error
}

type MsgType struct {
	Type string `json:"type"`
}

type MsgStart struct {
	Type   string `json:"type"`
	GameId int    `json:"gameid"`
	Color  int    `json:"color"`
}

type MsgMove struct {
	Type       string `json:"type"`
	Move       string `json:"move"`
	State      string `json:"state"`
	BlackRemTime int    `json:"blackRemTime"`
	WhiteRemTime int    `json:"whiteRemTime"`
}

type MsgMoveStatus struct {
	Type       string `json:"type"`
	Move       string `json:"move"`
	BlackRemTime int    `json:"blackRemTime"`
	WhiteRemTime int    `json:"whiteRemTime"`
	State      string `json:"state"`
	Code       string `json:"code"`
}

type MsgAbort struct {
	Type string `json:"type"`
}

type MsgGameOver struct {
	Type   string `json:"type"`
	Winner int    `json:"winner"`
	By     string `json:"by"`
}

type MsgChat struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type MsgSyncState struct {
	Type string `json:"type"`
	Gameid int `json:"gameid"`
	BlackName string `json:"blackname"`
	WhiteName string `json:"whitename"`
	State string `json:"state"`
	History StringArray `json:"history"`
	BlackRemTime int `json:"blackRemTime"`
	WhiteRemTime int `json:"whiteRemTime"`
}
