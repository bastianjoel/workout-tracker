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

func (a *App) registerWorkoutController(apiGroup *echo.Group, apiGroupPublic *echo.Group) {
	wc := controller.NewWorkoutController(container.NewContainer(a.db))

	workoutGroup := apiGroup.Group("/workouts")
	workoutGroup.GET("", wc.GetWorkouts).Name = "workouts-list"
	workoutGroup.POST("", wc.CreateWorkout).Name = "workouts-create"
	workoutGroup.GET("/recent", wc.GetRecentWorkouts).Name = "workouts-recent"
	workoutGroup.GET("/calendar", wc.GetWorkoutCalendar).Name = "workouts-calendar"
	workoutGroup.GET("/:id", wc.GetWorkout).Name = "workout-get"
	workoutGroup.GET("/:id/breakdown", wc.GetWorkoutBreakdown).Name = "workout-breakdown"
	workoutGroup.GET("/:id/stats-range", wc.GetWorkoutRangeStats).Name = "workout-range-stats"
	workoutGroup.GET("/:id/download", wc.DownloadWorkout).Name = "workout-download"
	workoutGroup.PUT("/:id", wc.UpdateWorkout).Name = "workout-update"
	workoutGroup.POST("/:id/toggle-lock", wc.ToggleWorkoutLock).Name = "workout-toggle-lock"
	workoutGroup.POST("/:id/refresh", wc.RefreshWorkout).Name = "workout-refresh"
	workoutGroup.POST("/:id/share", wc.CreateWorkoutShare).Name = "workout-share"
	workoutGroup.DELETE("/:id", wc.DeleteWorkout).Name = "workout-delete"
	workoutGroup.DELETE("/:id/share", wc.DeleteWorkoutShare).Name = "workout-share-delete"

	apiGroupPublic.GET("/workouts/public/:uuid", wc.GetPublicWorkout).Name = "workout-public"
	apiGroupPublic.GET("/workouts/public/:uuid/breakdown", wc.GetPublicWorkoutBreakdown).Name = "workout-public-breakdown"
	apiGroupPublic.GET("/workouts/public/:uuid/stats-range", wc.GetPublicWorkoutRangeStats).Name = "workout-public-range-stats"
}
