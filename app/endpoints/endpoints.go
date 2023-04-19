package endpoints

import (
	"fmt"
	"net/http"
	db "set-game/app/database"

	//m "set-game/app/models"
	lob_m "set-game/app/models/lobby"

	req_m "set-game/app/models/requests"
	res_m "set-game/app/models/responses"
	usr_m "set-game/app/models/user"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// func readBody(с echo.Context) ([]byte, error) {
// 	bytes := []byte{}
// 	_, err := с.Request().Body.Read(bytes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bytes, nil
// }

type Endpoints struct{}

func (endp *Endpoints) Register(с echo.Context) error {
	request := req_m.RegRequest{}
	err := с.Bind(&request)
	if err != nil {
		return с.JSON(http.StatusOK, res_m.NewRegResponseBad(map[string]interface{}{"message": err.Error()}))
	}
	if db.Users[request.Nickname] != "" {
		return с.JSON(http.StatusOK, res_m.NewRegResponseBad(map[string]interface{}{"message": "This user is already registrated"}))
	}
	//registerUser(request)
	user := usr_m.NewUser(request)
	db.RegisterUser(user)
	return с.JSON(http.StatusOK, res_m.NewRegResponseOk(user))
}

func (endp *Endpoints) Create(c echo.Context) error {
	request := req_m.RequestWithTocken{}
	c.Bind(&request)

	flag, _ := db.CheckTocken(request.AccessToken)
	if !flag {
		return c.JSON(http.StatusOK, res_m.CreateResponse{Success: false, Exception: map[string]interface{}{"message": "Wrong Access Tocken"}})
	}

	id, lobby := lob_m.NewLobby()
	db.Lobbies[id] = lobby
	db.EnterGame(request.AccessToken, id)
	return c.JSON(http.StatusOK, res_m.CreateResponse{Success: true, GameId: id})
}

func (endp *Endpoints) GetLobbies(с echo.Context) error {
	return с.JSON(http.StatusOK, db.GetLobbies())
}

func (endp *Endpoints) GetLobbiesInfo(с echo.Context) error {
	return с.JSON(http.StatusOK, db.Lobbies)
}

func (endp *Endpoints) JoinToGame(с echo.Context) error {
	request := req_m.JoinRequest{}
	с.Bind(&request)
	return с.JSON(http.StatusOK, db.EnterGame(request.AccessToken, request.GameId))
}

func (endp *Endpoints) GetUsers(с echo.Context) error {
	return с.JSON(http.StatusOK, db.Users)
}

func (endp *Endpoints) GetPlayers(с echo.Context) error {
	return с.JSON(http.StatusOK, db.PlayerDatas)
}

func (endp *Endpoints) Field(с echo.Context) error {
	accessTocken := req_m.RequestWithTocken{}
	с.Bind(&accessTocken)
	pd, flag := db.PlayerDatas[accessTocken.AccessToken]
	exaption := map[string]interface{}{}
	if !flag {
		exaption["mess"] = fmt.Sprintf("bad access token - %s", accessTocken.AccessToken)
		return с.JSON(http.StatusOK, res_m.BadCardsResponse(exaption))
	}
	gameId := pd.CurrentLobby
	_, flag = db.Lobbies[gameId]
	if !flag {
		exaption["mess"] = fmt.Sprintf("Bad lobby id - %d", gameId)
		return с.JSON(http.StatusOK, res_m.BadCardsResponse(exaption))
	}

	return с.JSON(http.StatusOK, res_m.GoodCardsResponse(db.Lobbies[gameId].ActiveCards))
}

func (endp *Endpoints) AddCards(с echo.Context) error {
	accessTocken := req_m.RequestWithTocken{}
	с.Bind(&accessTocken)
	pd, flag := db.PlayerDatas[accessTocken.AccessToken]
	exaption := map[string]interface{}{}
	if !flag {
		exaption["message"] = fmt.Sprintf("bad access token - %s", accessTocken.AccessToken)
		return с.JSON(http.StatusOK, res_m.BadCardsResponse(exaption))
	}
	gameId := pd.CurrentLobby
	lobby, flag := db.Lobbies[gameId]
	if !flag {
		exaption["message"] = fmt.Sprintf("Bad lobby id - %d", gameId)
		return с.JSON(http.StatusOK, res_m.BadCardsResponse(exaption))
	}
	flag, mess := lobby.MakeActiveCards(3)
	db.Lobbies[gameId] = lobby
	if !flag {
		exaption["message"] = mess
		return с.JSON(http.StatusOK, res_m.MixedCardsResponse(db.Lobbies[gameId].ActiveCards, exaption))
	}
	db.Lobbies[gameId] = lobby
	return с.JSON(http.StatusOK, res_m.GoodCardsResponse(db.Lobbies[gameId].ActiveCards))
}

func (endp *Endpoints) FindSet(с echo.Context) error {
	accessTocken := req_m.RequestWithTocken{}
	с.Bind(&accessTocken)
	pd, flag := db.PlayerDatas[accessTocken.AccessToken]
	exaption := map[string]interface{}{}
	if !flag {
		exaption["message"] = fmt.Sprintf("bad access token - %s", accessTocken.AccessToken)
		return с.JSON(http.StatusOK, res_m.BadCardsResponse(exaption))
	}
	gameId := pd.CurrentLobby
	lobby, flag := db.Lobbies[gameId]
	if !flag {
		exaption["message"] = fmt.Sprintf("Bad lobby id - %d", gameId)
		return с.JSON(http.StatusOK, res_m.BadCardsResponse(exaption))
	}
	flag, mess, set := lobby.FindSet()
	db.Lobbies[gameId] = lobby
	if !flag {
		exaption["mess"] = mess
		return с.JSON(http.StatusOK, res_m.MixedCardsResponse(set, exaption))
	}

	return с.JSON(http.StatusOK, res_m.GoodCardsResponse(set))
}

func (endp *Endpoints) Pick(c echo.Context) error {
	request := req_m.PickRequest{}
	c.Bind(&request)
	pd, flag := db.PlayerDatas[request.AccessToken]
	exaption := map[string]interface{}{}
	if !flag {
		exaption["message"] = fmt.Sprintf("bad access token - %s", request.AccessToken)
		return c.JSON(http.StatusOK, res_m.NewPickResponse(flag, exaption, -1))
	}
	gameId := pd.CurrentLobby
	lobby, flag := db.Lobbies[gameId]
	if !flag {
		exaption["message"] = fmt.Sprintf("Bad lobby id - %d", gameId)
		return c.JSON(http.StatusOK, res_m.NewPickResponse(flag, exaption, -1))
	}
	flag, mess := lobby.Pick(request.Cards...)
	if flag {
		pd.Score += 1
		db.PlayerDatas[request.AccessToken] = pd
	}
	exaption["message"] = mess
	db.Lobbies[gameId] = lobby
	return c.JSON(http.StatusOK, res_m.NewPickResponse(flag, exaption, pd.Score))
}

func (endp *Endpoints) LeaveFromLobby(c echo.Context) error {
	request := req_m.RequestWithTocken{}
	c.Bind(&request)
	flag, mess := db.LeaveFromLobby(request.AccessToken)
	response := map[string]interface{}{}
	response["success"] = flag
	response["message"] = mess
	return c.JSON(http.StatusOK, response)
}

func (endp *Endpoints) DeleteLobby(c echo.Context) error {
	request := req_m.JoinRequest{}
	c.Bind(&request)

	flag, mess := db.DeleteLobby(request.GameId)
	response := map[string]interface{}{}
	response["success"] = flag
	response["mess"] = mess
	return c.JSON(http.StatusOK, response)
}

func (endp *Endpoints) SayToConns(с echo.Context) error {
	mess := []byte{}
	body := с.Request().Body
	body.Read(mess)
	for conn := range db.Conns {
		conn.WriteMessage(websocket.BinaryMessage, mess)
	}
	return nil
}
