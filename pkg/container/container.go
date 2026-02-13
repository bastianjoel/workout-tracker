package container

import (
	"github.com/jovandeginste/workout-tracker/v2/pkg/database"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Container struct {
	db *gorm.DB
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{db: db}
}

func (c *Container) GetDB() *gorm.DB {
	return c.db
}

func (c *Container) GetUser(e echo.Context) *database.User {
	d := e.Get("user_info")

	u, ok := d.(*database.User)
	if !ok {
		u = database.AnonymousUser()
	}

	u.SetContext(e.Request().Context())

	return u
}
