package session

import (
	"context"
	"errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	cntxt "github.com/m50/shinidex/pkg/context"
	"github.com/m50/shinidex/pkg/database"
	"github.com/m50/shinidex/pkg/types"
)

var (
	ErrNotAuthed      = errors.New("no active session found")
	ErrNotEchoContext = errors.New("not a valid echo context")
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

	delete(sess.Values, "UserID")
	sess.Save(c.Request(), c.Response())
	return nil
}

func IsLoggedInContext(ctx context.Context) bool {
	return IsLoggedIn(cntxt.ToEcho(ctx))
}

func IsLoggedIn(c echo.Context) bool {
	_, err := GetAuthedUser(c)
	return err == nil
}

func GetAuthedUserContext(ctx context.Context) (*types.User, error) {
	c := cntxt.ToEcho(ctx)
	if c == nil {
		return nil, ErrNotEchoContext
	}
	return GetAuthedUser(c)
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
	user, err := db.Users().FindByID(cntxt.FromEcho(c), userID)
	if err != nil {
		return nil, ErrNotAuthed
	}
	c.Set("user", user)

	return &user, nil
}
