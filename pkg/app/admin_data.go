package app

import (
	"strconv"

	"github.com/jovandeginste/workout-tracker/v2/pkg/database"
	"github.com/labstack/echo/v4"
)

func (a *App) getUser(c echo.Context) (*database.User, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}

	u, err := database.GetUserByID(a.db, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
