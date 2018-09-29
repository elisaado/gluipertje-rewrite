package db

import (
	"fmt"
	"sort"

	"github.com/asdine/storm"
	"github.com/elisaado/gluipertje-rewrite/models"
)

var DB *storm.DB

func InitDB() {
	var err error
	DB, err = storm.Open("./db/db.db")

	if err != nil {
		fmt.Println(err)
	}
}

func SafeUser(u models.User) models.User {
	u.Password = ""
	u.Token = ""
	return u
}

func SafeUsers(us []models.User) []models.User {
	for i := range us {
		us[i] = SafeUser(us[i])
	}

	return us
}

func SafeMessage(m models.Message) models.Message {
	if m.From == (models.User{}) {
		var u models.User
		DB.One("id", m.FromId, &u)

		m.From = u
	}

	m.From = SafeUser(m.From)

	return m
}

func SafeMessages(ms []models.Message) []models.Message {
	for i := range ms {
		ms[i] = SafeMessage(ms[i])
	}

	sort.Slice(ms, func(i, j int) bool {
		return ms[i].CreatedAt.After(ms[j].CreatedAt)
	})

	return ms
}
