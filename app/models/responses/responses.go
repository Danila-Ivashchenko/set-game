package responses

import (
	usr_m "set-game/app/models/user"
	set_m "set-game/set/models"
)

// response for user/register

type RegResponse struct {
	Success     bool                   `json:"success"`
	Exaption    map[string]interface{} `json:"exaption"`
	Nickname    string                 `json:"nickname"`
	AccessToken string                 `json:"accessToken"`
}

func NewRegResponseBad(exaption map[string]interface{}) RegResponse {
	respose := RegResponse{}
	respose.Exaption = exaption
	respose.Success = false
	return respose
}

func NewRegResponseOk(user usr_m.User) RegResponse {
	respose := RegResponse{}
	respose.Nickname = user.Nickname
	respose.AccessToken = user.AccessToken
	respose.Success = true
	return respose
}

// response for set/create

type CreateResponse struct {
	Success   bool                   `json:"success"`
	Exception map[string]interface{} `json:"exception"`
	GameId    int                    `json:"gameId"`
}

// response for set/field, /set/field/add, set/fild/find

type CardsResponse struct {
	Success   bool                   `json:"success"`
	Exception map[string]interface{} `json:"exception"`
	Cards     []set_m.Card           `json:"cards"`
}

func BadCardsResponse(exaption map[string]interface{}) CardsResponse {
	respnose := CardsResponse{}
	respnose.Success = false
	respnose.Exception = exaption
	return respnose
}

func GoodCardsResponse(cards []set_m.Card) CardsResponse {
	respnose := CardsResponse{}
	respnose.Success = true
	respnose.Cards = cards
	return respnose
}

func MixedCardsResponse(cards []set_m.Card, exaption map[string]interface{}) CardsResponse {
	respnose := CardsResponse{}
	respnose.Success = true
	respnose.Exception = exaption
	respnose.Cards = cards
	return respnose
}

// response for /set/pick

type PickResponse struct {
	IsSet    bool                   `json:"isSet"`
	Exaption map[string]interface{} `json:"exaption"`
	Score    int                    `json:"score"`
}

func NewPickResponse(isset bool, exaprion map[string]interface{}, score int) PickResponse {
	return PickResponse{IsSet: isset, Exaption: exaprion, Score: score}
}

// response for /set/room/list/enter

type LobbyEnterResponce struct {
	Success  bool                   `json:"success"`
	Exaption map[string]interface{} `json:"exaption"`
	GameId   int                    `json:"gameId"`
}

func BadLobbyEnterResponce(exaption map[string]interface{}) LobbyEnterResponce {
	ler := LobbyEnterResponce{}
	ler.Success = false
	ler.Exaption = exaption
	ler.GameId = -1
	return ler
}

func GoodLobbyEnterResponce(gameId int) LobbyEnterResponce {
	ler := LobbyEnterResponce{}
	ler.Success = true
	ler.GameId = gameId
	return ler
}
