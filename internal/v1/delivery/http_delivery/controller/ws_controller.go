package controller

import "github.com/gofiber/websocket/v2"

type WSController struct {
}

func NewWSController() *WSController {
	return &WSController{}
}

func (c *WSController) Handle(con *websocket.Conn) {
	for {
		_, msg, err := con.ReadMessage()
		if err != nil {
			return
		}
		_ = con.WriteMessage(websocket.TextMessage, msg)
	}
}
