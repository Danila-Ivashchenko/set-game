package database

import (
	"fmt"
	lob_m "set-game/app/models/lobby"
	res_m "set-game/app/models/responses"
	usr_m "set-game/app/models/user"
)

func RegisterUser(user usr_m.User) {
	AccessTokens[user.AccessToken] = user.Nickname
	Users[user.Nickname] = user.AccessToken
}

func AddPlayer(username string, currentLobby int) usr_m.PlayerData {
	pd := usr_m.PlayerData{}
	pd.UserName = username
	pd.CurrentLobby = currentLobby
	return pd
}

func CheckTocken(tocken string) (bool, string) {
	val, ok := AccessTokens[tocken]
	return ok, val
}

func CheckUserInGames(username string) bool {
	for _, lobby := range Lobbies {
		for _, value := range lobby.Players {
			if value == username {
				return true
			}
		}
	}
	return false
}

func GetLobbies() map[string][]map[string]int {
	games := []map[string]int{}
	for key := range Lobbies {
		games = append(games, map[string]int{"gameId": key})
	}
	return map[string][]map[string]int{"games": games}
}

func EnterGame(accessTocken string, gameId int) res_m.LobbyEnterResponce {
	_, flag := Lobbies[gameId]
	if !flag {
		return res_m.BadLobbyEnterResponce(map[string]interface{}{"message": "This lobby is not exist"})
	}
	username := AccessTokens[accessTocken]
	if CheckUserInGames(username) {
		return res_m.BadLobbyEnterResponce(map[string]interface{}{"message": "You are in the other lobby"})
	}
	var err error
	Lobbies[gameId], err = lob_m.JoinToLobby(Lobbies[gameId], username)
	if err != nil {
		return res_m.BadLobbyEnterResponce(map[string]interface{}{"message": err.Error()})
	}
	PlayerDatas[accessTocken] = AddPlayer(username, gameId)

	return res_m.GoodLobbyEnterResponce(gameId)
}

func LeaveFromLobby(accessTocken string) (bool, string) {
	pd, flag := PlayerDatas[accessTocken]
	if !flag {
		return false, "Bad access tocken"
	}
	lobby, flag := Lobbies[pd.CurrentLobby]
	if !flag {
		return false, fmt.Sprintf("Lobby № - %d doesn't exist", pd.CurrentLobby)
	}
	for i := range lobby.Players {
		if lobby.Players[i] == pd.UserName {
			lobby.Players[lobby.CountPlayers-1], lobby.Players[i] = lobby.Players[i], lobby.Players[lobby.CountPlayers-1]
			lobby.Players[lobby.CountPlayers-1] = ""
			lobby.CountPlayers--
			if lobby.CountPlayers == 0 {
				delete(Lobbies, pd.CurrentLobby)
			} else {
				Lobbies[pd.CurrentLobby] = lobby
			}
			delete(PlayerDatas, accessTocken)
			return true, ""
		}
	}
	return false, "Something bad"

}

func DeleteLobby(id int) (bool, string) {
	lobby, flag := Lobbies[id]
	if !flag {
		return false, "This lobby is not exist"
	}
	usernames := lobby.Players
	for _, username := range usernames {
		LeaveFromLobby(Users[username])
	}
	delete(Lobbies, id)
	return true, ""
}
