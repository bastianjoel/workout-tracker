package app

import (
	"net/http"

	"github.com/jovandeginste/workout-tracker/v2/pkg/api"
	"github.com/labstack/echo/v4"
)

func (a *App) registerAPIV2UserRoutes(e *echo.Group) {
	e.GET("/whoami", a.apiV2WhoamiHandler).Name = "api-v2-whoami"
	e.GET("/totals", a.apiV2TotalsHandler).Name = "api-v2-totals"
	e.GET("/records", a.apiV2RecordsHandler).Name = "api-v2-records"
	e.GET("/:id", a.apiV2UserShowHandler).Name = "api-v2-user-show"
}

// apiV2WhoamiHandler returns current user information
// @Summary      Show current user information
// @Description  Returns the profile and settings of the currently authenticated user
// @Tags         users
// @Produce      json
// @Success      200  {object}  api.Response[api.UserProfileResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      403  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /whoami [get]
func (a *App) apiV2WhoamiHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	resp := api.Response[api.UserProfileResponse]{
		Results: api.NewUserProfileResponse(user),
	}

	return c.JSON(http.StatusOK, resp)
}

// apiV2TotalsHandler returns user's workout totals
// @Summary      Get workout totals
// @Description  Returns aggregated totals for the current user's workouts
// @Tags         users
// @Produce      json
// @Success      200  {object}  api.Response[api.TotalsResponse]
// @Failure      500  {object}  api.Response[any]
// @Router       /totals [get]
func (a *App) apiV2TotalsHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	totals, err := user.GetDefaultTotals()
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	resp := api.Response[api.TotalsResponse]{
		Results: api.NewTotalsResponse(totals),
	}

	return c.JSON(http.StatusOK, resp)
}

// apiV2RecordsHandler returns user's workout records
// @Summary      Get workout records
// @Description  Returns personal records for the current user across all workout types
// @Tags         users
// @Produce      json
// @Success      200  {object}  api.Response[[]api.WorkoutRecordResponse]
// @Failure      500  {object}  api.Response[any]
// @Router       /records [get]
func (a *App) apiV2RecordsHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	records, err := user.GetAllRecords()
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	resp := api.Response[[]api.WorkoutRecordResponse]{
		Results: api.NewWorkoutRecordsResponse(records),
	}

	return c.JSON(http.StatusOK, resp)
}

// apiV2UserShowHandler returns a specific user's workout records
// @Summary      Get user's public profile
// @Description  Returns public profile information and records for a specific user
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  api.Response[[]api.WorkoutRecordResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      403  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /{id} [get]
func (a *App) apiV2UserShowHandler(c echo.Context) error {
	u, err := a.getUser(c)
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	if u.IsAnonymous() {
		return a.renderAPIV2Error(c, http.StatusForbidden, api.ErrNotAuthorized)
	}

	records, err := u.GetAllRecords()
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	resp := api.Response[[]api.WorkoutRecordResponse]{
		Results: api.NewWorkoutRecordsResponse(records),
	}

	return c.JSON(http.StatusOK, resp)
}
