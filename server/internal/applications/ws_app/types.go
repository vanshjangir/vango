package ws_app

const (
	TIMER_OUT = 1
	CLIENT_OUT = 2
	LOCAL_OUT = 3
)

type ClientRecvResult struct {
	data []byte
	err error
}

type MsgType struct {
    Type    string  `json:"type"`
}

type MsgStart struct {
	Type string	`json:"type"`
}

type MsgMove struct {
    Type    string  `json:"type"`
    Move    string  `json:"move"`
	PassedTime int	`json:"time"`
	State 	string	`json:"state"`
}

type MsgMoveStatus struct {
    Type    string  `json:"type"`
    Move    string  `json:"move"`
	PassedTime int	`json:"time"`
	State 	string	`json:"state"`
	InvalidTurn	bool 	`json:"invalidTurn"`
	InvalidMove	bool	`json:"invalidMove"`
}

type MsgAbort struct {
    Type    string  `json:"type"`
}

type MsgGameOver struct {
    Type    string  `json:"gameover"`
	By		string	`json:"by"`
	Winner	int		`json:"winner"`
}

type MsgChat struct {
    Type    string  `json:"type"`
    Text    string  `json:"text"`
}
