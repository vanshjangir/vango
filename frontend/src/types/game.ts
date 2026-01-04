export interface MsgStart {
  type: "start";
  color: number;
  gameid: number;
}

export interface MsgStop {
  type: string;
}

export interface MsgPass {
  type: string;
}

export interface MsgAbort {
  type: string;
}

export interface MsgGameover {
  type: "gameover";
  winner: number;
  by: string;
}

export interface MsgMoveStatus {
  type: "movestatus";
  move: string;
  code: string;
  state: string;
  blackRemTime: number;
  whiteRemTime: number;
}

export interface MsgMove {
  type: "move";
  move: string;
  state: string;
  blackRemTime: number;
  whiteRemTime: number;
}

export interface MsgSync {
  type: "syncstate";
  gameid: string;
  blackname: string;
  whitename: string;
  state: string;
  history: string[];
  blackRemTime: number;
  whiteRemTime: number;
}

export interface MsgChat {
  type: "chat";
  message: string;
}

export interface GameState {
  gameId: string;
  pname: string;
  opname: string;
  color: number;
  state: number[][];
  history: string[];
}

export interface ChatMessage {
  type: "sent" | "received";
  text: string;
}

export interface UserProfileData {
  name: string;
  rating: number;
  gamesPlayed: number;
  wins: number;
  losses: number;
  highestRating: number;
  recentGames: {
    gameid: string;
    opponent: string;
    result: string;
    created_at: string;
  }[];
}
