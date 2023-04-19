package game

import (
	"fmt"
	db "set-game/app/database"

	f "set-game/set/functions"
	m "set-game/set/models"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Game struct {
	upgrader websocket.Upgrader
	cards    []m.Card
}

func NewGame() *Game {
	g := new(Game)
	g.upgrader = websocket.Upgrader{}

	g.cards = f.GenerateBoard()
	return g
}

func (g *Game) Set(c echo.Context) error {
	ws, err := g.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	for {
		db.Conns[ws] = true
		var err error
		var mess []byte

		defer c.Logger().Error(err)
		defer delete(db.Conns, ws)
		defer ws.Close()
		err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", g.cards)))
		if err != nil {
			break
		}

		_, mess, err = ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(string(mess))
	}
	return nil
}

func (g *Game) Hello(c echo.Context) error {
	ws, err := g.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	func() {
		db.Conns[ws] = true
		var err error
		var mess []byte

		defer c.Logger().Error(err)
		defer delete(db.Conns, ws)
		defer ws.Close()

		for {
			err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", g.cards)))
			if err != nil {
				break
			}

			_, mess, err = ws.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println(string(mess))
		}
	}()
	return nil
}
