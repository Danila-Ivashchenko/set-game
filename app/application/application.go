package application

import (
	"github.com/labstack/echo/v4"
)

type IEndpoints interface {
	Register(c echo.Context) error
	GetUsers(c echo.Context) error
	GetPlayers(c echo.Context) error
	SayToConns(c echo.Context) error
	Create(c echo.Context) error
	GetLobbies(c echo.Context) error
	GetLobbiesInfo(c echo.Context) error
	JoinToGame(c echo.Context) error
	Field(c echo.Context) error
	AddCards(с echo.Context) error
	FindSet(с echo.Context) error
	Pick(с echo.Context) error
	LeaveFromLobby(с echo.Context) error
	DeleteLobby(с echo.Context) error
}

type IGame interface {
	Hello(c echo.Context) error
	Set(c echo.Context) error
}

type App struct {
	e    *echo.Echo
	endp IEndpoints
	g    IGame
}

func NewApp(endp IEndpoints, g IGame) *App {
	app := new(App)
	app.e = echo.New()
	app.endp = endp
	app.g = g

	//app.fillGroups()
	app.routing()
	return app
}

func (app *App) routing() {
	app.e.POST("/user/register", app.endp.Register)
	app.e.POST("/user/get_all", app.endp.GetUsers)
	app.e.POST("/set/player/list", app.endp.GetPlayers)
	app.e.POST("/say", app.endp.SayToConns)
	app.e.POST("/set/room/create", app.endp.Create)
	app.e.POST("/set/room/list", app.endp.GetLobbies)
	app.e.POST("/set/room/list/info", app.endp.GetLobbiesInfo)
	app.e.POST("/set/room/list/enter", app.endp.JoinToGame)
	app.e.POST("/set/room/leave", app.endp.LeaveFromLobby)
	app.e.POST("/set/room/delete", app.endp.DeleteLobby)
	app.e.POST("/set/field", app.endp.Field)
	app.e.POST("/set/field/add", app.endp.AddCards)
	app.e.POST("/set/field/find", app.endp.FindSet)
	app.e.POST("/set/pick", app.endp.Pick)
	app.e.GET("/ws/game/hello", app.g.Hello)
	app.e.GET("/ws/game/set", app.g.Set)
}

// Get port like "8080", not ":8080"
func (app *App) Engine(port string) {
	app.e.Logger.Fatal(app.e.Start(":" + port))
}
