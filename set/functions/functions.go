package functions

import (
	"math/rand"
	m "set-game/set/models"
	"time"
)

func shake(cards *[]m.Card) {
	rand.Seed(time.Now().UnixNano())
	var n int = 81
	indexes := make([]int, 81)

	for i := 0; i < n; i++ {
		indexes[i] = i
	}

	for n > 0 {
		k := rand.Intn(n)
		n--
		idx := indexes[k]
		(*cards)[idx], (*cards)[n] = (*cards)[n], (*cards)[idx]
	}
}

func GenerateBoard() []m.Card {
	cards := make([]m.Card, 81)
	counter := 0
	var n uint8 = 3
	var i, j, c, v uint8
	for i = 0; i < n; i++ {
		for j = 0; j < n; j++ {
			for c = 0; c < n; c++ {
				for v = 0; v < n; v++ {
					cards[counter] = m.Card{Id: counter, Shape: i, Count: j, Color: c, Fill: v}
					counter++
				}
			}
		}
	}
	shake(&cards)
	return cards
}

func Is_set(card1, card2, card3 m.Card) bool {
	Is_set := true

	counts := [3]uint8{card1.Count, card2.Count, card3.Count}
	colors := [3]uint8{card1.Color, card2.Color, card3.Color}
	shapes := [3]uint8{card1.Shape, card2.Shape, card3.Shape}
	fills := [3]uint8{card1.Fill, card2.Fill, card3.Fill}

	propertys := [4][3]uint8{counts, colors, shapes, fills}

	for i := 0; i < len(propertys); i++ {
		if propertys[i][0] == propertys[i][1] && propertys[i][0] == propertys[i][2] {
			continue
		}
		if propertys[i][0] != propertys[i][1] && propertys[i][0] != propertys[i][2] && propertys[i][1] != propertys[i][2] {
			continue
		}

		Is_set = false
		break
	}
	return Is_set
}

func FindSets(cards []m.Card) [][3]m.Card {
	sets := [][3]m.Card{}

	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			for c := j + 1; c < len(cards); c++ {
				if Is_set(cards[i], cards[j], cards[c]) {
					sets = append(sets, [3]m.Card{cards[i], cards[j], cards[c]})
				}
			}
		}
	}
	return sets
}

func FindSet(cards []m.Card) []m.Card {
	set := []m.Card{}

	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			for c := j + 1; c < len(cards); c++ {
				if Is_set(cards[i], cards[j], cards[c]) {
					return []m.Card{cards[i], cards[j], cards[c]}
				}
			}
		}
	}
	return set
}
