package user

import (
	"crypto/md5"
	"encoding/hex"
	req_m "set-game/app/models/requests"

	"golang.org/x/crypto/bcrypt"
)

func generateToken(username string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

type User struct {
	Nickname    string `json:"nickname"`
	AccessToken string `json:"accessToken"`
}

func NewUser(request req_m.RegRequest) User {
	return User{Nickname: request.Nickname, AccessToken: generateToken(request.Nickname)}
}

type PlayerData struct {
	UserName     string
	Score        int
	CurrentLobby int
}
