package app

import (
	"github.com/jovandeginste/workout-tracker/v2/pkg/container"
	"github.com/jovandeginste/workout-tracker/v2/pkg/controller"
	"github.com/labstack/echo/v4"
)

func (a *App) registerEquipmentController(apiGroup *echo.Group) {
	ec := controller.NewEquipmentController(container.NewContainer(a.db))

	apiGroup.GET("/equipment", ec.GetEquipmentList).Name = "equipment-list"
	apiGroup.GET("/equipment/:id", ec.GetEquipment).Name = "equipment-get"
	apiGroup.POST("/equipment", ec.CreateEquipment).Name = "equipment-create"
	apiGroup.PUT("/equipment/:id", ec.UpdateEquipment).Name = "equipment-update"
	apiGroup.DELETE("/equipment/:id", ec.DeleteEquipment).Name = "equipment-delete"
}
