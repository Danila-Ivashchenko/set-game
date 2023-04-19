package lobby

import (
	"fmt"
	set_m "set-game/set/models"
	set_f "set-game/set/functions"

	err_m "set-game/app/models/errors"
)

const(
	maxPlayers = 10
)

var(
	LobbiesCount int = 0
)

type Lobby struct {
	Alive        bool         `json:"alive"`
	MaxPlayers   int          `json:"max_players"`
	Players      []string     `json:"players"`
	Cards        []set_m.Card `json:"stashed_cards"`
	ActiveCards  []set_m.Card `json:"active_cards"`
	CountPlayers int          `json:"count_players"`
}

func NewLobby() (int, Lobby) {
	l := Lobby{}
	l.Alive = true
	l.MaxPlayers = maxPlayers
	l.Players = make([]string, maxPlayers)
	l.CountPlayers = 0
	l.Cards = set_f.GenerateBoard()
	l.toAcriveCards(12)
	l.PrepareCards()
	id := LobbiesCount
	LobbiesCount++
	return id, l
}

func JoinToLobby(l Lobby, user string) (Lobby, error) {
	if l.CountPlayers >= l.MaxPlayers {
		return l, err_m.NewErrAnswer("This lobby is full")
	}
	for _, value := range l.Players {
		if value == user {
			return l, err_m.NewErrAnswer("You are already in the lobby")
		}
	}
	l.Players[l.CountPlayers] = user
	l.CountPlayers++

	return l, nil
}

func (l *Lobby) toAcriveCards(n int) bool {
	if n > len(l.Cards) {
		return false
	}
	cards := make([]set_m.Card, n)
	for i := 0; i < n; i++ {
		cards[i] = l.Cards[i]
	}
	l.Cards = l.Cards[n:]
	l.ActiveCards = append(l.ActiveCards, cards...)
	return true
}

func (l *Lobby) MakeActiveCards(n int) (bool, string) {
	if !l.Alive {
		return false, "This game is ended"
	}
	flag := l.toAcriveCards(n)
	if !flag {
		return false, "There are no stashed cards here"
	}
	flag = l.PrepareCards()
	if !flag {
		l.Alive = false
		return false, "There are no sets here"
	}
	return true, ""
}

func (l *Lobby) PrepareCards() bool {
	cards := set_f.FindSets(l.ActiveCards)
	if len(cards) == 0 {
		for i := range l.ActiveCards {
			for j := range l.Cards {
				l.ActiveCards[i] = l.Cards[j]
				cards = set_f.FindSets(l.ActiveCards)
				if len(cards) > 0 {
					break
				}
			}
			if len(cards) > 0 {
				break
			}
		}
	}
	return len(cards) > 0
}

func (l *Lobby) IsAlive() bool {
	if !l.Alive {
		return false
	}
	flag := len(append(l.ActiveCards, l.Cards...)) != 0

	l.Alive = flag
	return flag
}

func (l *Lobby) FindSet() (bool, string, []set_m.Card) {
	set := []set_m.Card{}
	l.IsAlive()
	if !l.Alive {
		return false, "This game is ended", set
	}
	set = set_f.FindSet(l.ActiveCards)
	if len(set) == 0 {
		return false, "There are no sets in active cards", set
	}

	return true, "", set
}

func (l *Lobby) FindActiveCard(id int) (set_m.Card, bool) {
	for _, value := range l.ActiveCards {
		if value.Id == id {
			return value, true
		}
	}
	return set_m.Card{}, false
}

func (l *Lobby) DeleteActiveCard(id int) bool {
	for i, value := range l.ActiveCards {
		if value.Id == id {
			l.ActiveCards[len(l.ActiveCards)-1], l.ActiveCards[i] = l.ActiveCards[i], l.ActiveCards[len(l.ActiveCards)-1]
			l.ActiveCards = l.ActiveCards[:len(l.ActiveCards)-1]
			return true
		}
	}
	return false
}

func (l *Lobby) Pick(Ids ...int) (bool, string) {
	if len(Ids) != 3 {
		return false, "There are not three cards"
	}
	cards := make([]set_m.Card, 3)
	flag := true
	for i := range cards {
		cards[i], flag = l.FindActiveCard(Ids[i])
		if !flag {
			return false, fmt.Sprintf("Cards with id - %d not exist", Ids[i])
		}
	}

	flag = set_f.Is_set(cards[0], cards[1], cards[2])
	mess := ""
	if flag {
		for i := range cards {
			l.DeleteActiveCard(cards[i].Id)
		}
	} else {
		mess = "Is not a set"
	}

	return flag, mess
}
