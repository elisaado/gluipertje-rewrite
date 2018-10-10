package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/elisaado/gluipertje-rewrite/db"
	"github.com/elisaado/gluipertje-rewrite/models"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c echo.Context) error {
	var us []models.User

	db.DB.All(&us)

	return c.JSON(http.StatusOK, db.SafeUsers(us))
}

func GetUserByIdOrUsername(c echo.Context) error {
	var u models.User
	idOrUsername := c.Param("idOrUsername")

	var err error

	// check if only numbers (ID)
	if match, _ := regexp.MatchString("^[0-9]+$", idOrUsername); match {
		id, err := strconv.Atoi(idOrUsername)
		if err != nil || id <= 0 || idOrUsername == "0" {
			return c.JSON(http.StatusBadRequest, "ID has to be an integer above 0")
		}

		err = db.DB.One("ID", id, &u)
	} else {
		// we know it's a username by now
		err = db.DB.One("Username", idOrUsername, &u)
	}

	if err == storm.ErrNotFound || u.ID == 0 {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, db.SafeUser(u))
}

func GetUserByToken(c echo.Context) error {
	var u models.User
	token := c.Param("token")

	err := db.DB.One("Token", token, &u)

	if err == storm.ErrNotFound || u.ID == 0 {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, db.SafeUser(u))
}

func CreateUser(c echo.Context) error {
	var u models.User

	if err := c.Bind(&u); err != nil {
		realError := err.(*echo.HTTPError)
		return c.JSON(realError.Code, realError.Message)
	}

	var missing []string
	if empty(u.Nickname) {
		missing = append(missing, "nickname")
	}
	if empty(u.Username) {
		missing = append(missing, "username")
	}
	if empty(u.Password) {
		missing = append(missing, "password")
	}

	if len(missing) > 0 {
		return c.JSON(http.StatusBadRequest, "Missing parameters "+strings.Join(missing, ", "))
	}

	if len(u.Nickname) > 50 {
		return c.JSON(http.StatusBadRequest, "nickname may not be longer than 50 characters")
	}
	if len(u.Username) > 50 {
		return c.JSON(http.StatusBadRequest, "username may not be longer than 50 characters")
	}

	if match, _ := regexp.MatchString("^[a-zA-Z0-9-_]+$", u.Username); !match {
		return c.JSON(http.StatusBadRequest, "Username may only contain 0-9, a-z, A-Z, - and _")
	}
	if match, _ := regexp.MatchString("^[0-9]+$", u.Username); match {
		return c.JSON(http.StatusBadRequest, "Username may not contain only numbers")
	}

	// duplicate user check
	var du models.User
	if err := db.DB.One("Username", strings.ToLower(u.Username), &du); err != storm.ErrNotFound {
		return c.JSON(http.StatusConflict, "User with username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}
	u.Password = string(hash)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Token = randomString(60)

	err = db.DB.Save(&u)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	u.Password = ""
	return c.JSON(http.StatusOK, u)
}

func RevokeToken(c echo.Context) error {
	var u models.User

	if err := c.Bind(&u); err != nil {
		realError := err.(*echo.HTTPError)
		return c.JSON(realError.Code, realError.Message)
	}

	p := u.Password
	err := db.DB.One("Username", u.Username, &u)
	if err == storm.ErrNotFound || u.ID == 0 {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	u.Token = randomString(60)
	db.DB.Save(&u)

	u.Password = ""
	return c.JSON(http.StatusOK, u)
}
