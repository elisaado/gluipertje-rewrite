package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elisaado/gluipertje-rewrite/config"
	"github.com/elisaado/gluipertje-rewrite/util"
	"golang.org/x/net/websocket"

	"github.com/asdine/storm"
	"github.com/elisaado/gluipertje-rewrite/db"
	"github.com/elisaado/gluipertje-rewrite/models"
	"github.com/labstack/echo"
)

var (
	MChan        chan string // channel with json encoded messages
	Subschribers []*websocket.Conn
)

func GetMessages(c echo.Context) error {
	var ms []models.Message

	db.DB.All(&ms)

	return c.JSON(http.StatusOK, db.SafeMessages(ms))
}

func GetMessagesByLimit(c echo.Context) error {
	var ms []models.Message

	limitStr := c.Param("limit")
	limit, err := strconv.Atoi(limitStr)
	if limit <= 0 || err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "limit has to be an integer above 0")
	}

	db.DB.All(&ms, storm.Limit(limit))

	return c.JSON(http.StatusOK, db.SafeMessages(ms))
}

func GetMessageById(c echo.Context) error {
	var m models.Message

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if id <= 0 || err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "limit has to be an integer above 0")
	}

	err = db.DB.One("ID", id, &m)

	if err == storm.ErrNotFound || m.ID == 0 {
		return c.JSON(http.StatusNotFound, "Message not found")
	}

	return c.JSON(http.StatusOK, db.SafeMessage(m))
}

func SendMessage(c echo.Context) error {
	var mr models.MessageRequest
	var m models.Message
	var u models.User

	if err := c.Bind(&mr); err != nil {
		realError := err.(*echo.HTTPError)
		return c.JSON(realError.Code, realError.Message)
	}

	var missing []string
	if empty(mr.Token) {
		missing = append(missing, "token")
	}
	if empty(mr.Text) {
		missing = append(missing, "text")
	}

	if len(missing) > 0 {
		return c.JSON(http.StatusBadRequest, "Missing parameters "+strings.Join(missing, ", "))
	}

	if len(mr.Text) >= 500 {
		return c.JSON(http.StatusBadRequest, "Text may not exceed 500 characters")
	}

	err := db.DB.One("Token", mr.Token, &u)
	if err == storm.ErrNotFound {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	m.From = u
	m.FromId = u.ID
	m.Type = "text"
	m.Text = mr.Text

	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	m.Text = mr.Text

	db.DB.Save(&m)

	m = db.SafeMessage(m)
	json, err := json.Marshal(&m)
	if err != nil {
		fmt.Println(err)
	}
	MChan <- string(json)
	return c.JSON(http.StatusOK, m)
}

func SendImage(c echo.Context) error {
	token := c.FormValue("token")
	text := c.FormValue("text")

	if empty(token) {
		return c.JSON(http.StatusBadRequest, "Missing parameter token")
	}

	var u models.User
	if err := db.DB.One("Token", token, &u); err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, "Missing parameter image")
	}

	if file.Size > 5000000 { // 5 mb
		return c.JSON(http.StatusBadRequest, "File size may not exceed 5 megabytes")
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server errror")
	}
	defer src.Close()

	// Destination
	dst, err := os.Create("images/" + randomString(5) + file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server errror")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal server errror")
	}

	var m models.Message
	now := time.Now()

	m.Type = "image"
	m.CreatedAt = now
	m.UpdatedAt = now
	m.Text = text
	m.From = u
	m.FromId = u.ID
	m.SRC = config.C.ExternalURL + "/api/" + dst.Name()

	db.DB.Save(&m)

	return c.JSON(http.StatusOK, db.SafeMessage(m))
}

func SubschribeToMessages(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		Subschribers = append(Subschribers, ws)
		select {}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func BroadcastNewMessages() {
	for m := range MChan {
		for i, subschriber := range Subschribers {
			err := websocket.Message.Send(subschriber, m)
			if err != nil {
				fmt.Println(err)
				subschriber.Close()
				Subschribers = util.RemoveConn(Subschribers, i)
			}
		}
	}
}
