package requests


type RegRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type RequestWithTocken struct {
	AccessToken string `json:"accessTocken"`
}

type JoinRequest struct {
	AccessToken string `json:"accsessTocken"`
	GameId      int    `json:"gameId"`
}

type PickRequest struct {
	AccessToken string `json:"accessTocken"`
	Cards       []int  `json:"cards"`
}