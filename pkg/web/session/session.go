package session

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
)

var (
	ErrNotAuthed = errors.New("no active session found")
)

func New(c echo.Context, user types.User) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["UserID"] = user.ID
	return sess.Save(c.Request(), c.Response())
}

func Close(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}

	delete(sess.Values, "UserID")
	sess.Save(c.Request(), c.Response())
	return nil
}

func GetAuthedUser(c echo.Context) (*types.User, error) {
	if user, ok := c.Get("user").(types.User); ok {
		return &user, nil
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}
	db := c.(database.DBContext).DB()

	foundID, ok := sess.Values["UserID"]
	if !ok {
		return nil, ErrNotAuthed
	}
	userID, ok := foundID.(string)
	if !ok {
		return nil, ErrNotAuthed
	}
	user, err := db.Users().FindByID(userID)
	if err != nil {
		return nil, err
	}
	c.Set("user", user)

	return &user, nil
}