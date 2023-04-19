package database

import (
	"github.com/gorilla/websocket"
	lob_m "set-game/app/models/lobby"
	usr_m "set-game/app/models/user"
)


var (
	LobbiesCount = 0
	Users        = map[string]string{}
	AccessTokens = map[string]string{}
	Conns        = map[*websocket.Conn]bool{}
	Lobbies      = map[int]lob_m.Lobby{}
	PlayerDatas  = map[string]usr_m.PlayerData{}
)
