package app

import (
	"net/http"

	"github.com/jovandeginste/workout-tracker/v2/pkg/api"
	"github.com/jovandeginste/workout-tracker/v2/pkg/database"
	"github.com/labstack/echo/v4"
)

func (a *App) registerAPIV2EquipmentRoutes(apiGroup *echo.Group) {
	apiGroup.GET("/equipment", a.apiV2EquipmentHandler).Name = "api-v2-equipment"
	apiGroup.GET("/equipment/:id", a.apiV2EquipmentGetHandler).Name = "api-v2-equipment-get"
	apiGroup.POST("/equipment", a.apiV2EquipmentCreateHandler).Name = "api-v2-equipment-create"
	apiGroup.PUT("/equipment/:id", a.apiV2EquipmentUpdateHandler).Name = "api-v2-equipment-update"
	apiGroup.DELETE("/equipment/:id", a.apiV2EquipmentDeleteHandler).Name = "api-v2-equipment-delete"
}

// apiV2EquipmentHandler returns a paginated list of equipment for the current user
// @Summary      List equipment
// @Description  Returns a paginated list of equipment (shoes, bikes, etc.) for the current user
// @Tags         equipment
// @Produce      json
// @Param        page     query      int  false  "Page number" default(1)
// @Param        per_page query      int  false  "Items per page" default(20)
// @Success      200  {object}  api.PaginatedResponse[api.EquipmentResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /equipment [get]
func (a *App) apiV2EquipmentHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	// Parse pagination parameters
	var pagination api.PaginationParams
	if err := c.Bind(&pagination); err != nil {
		return a.renderAPIV2Error(c, http.StatusBadRequest, err)
	}
	pagination.SetDefaults()

	// Get total count
	var totalCount int64
	if err := a.db.Model(&database.Equipment{}).Where(&database.Equipment{UserID: user.ID}).Count(&totalCount).Error; err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	// Get paginated equipment
	var equipment []*database.Equipment
	db := a.db.Where(&database.Equipment{UserID: user.ID}).
		Order("name DESC").
		Limit(pagination.PerPage).
		Offset(pagination.GetOffset())

	if err := db.Find(&equipment).Error; err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	// Convert to API response
	results := api.NewEquipmentListResponse(equipment)

	resp := api.PaginatedResponse[api.EquipmentResponse]{
		Results:    results,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: pagination.CalculateTotalPages(totalCount),
		TotalCount: totalCount,
	}

	return c.JSON(http.StatusOK, resp)
}

// apiV2EquipmentGetHandler returns a single equipment by ID
// @Summary      Get equipment details
// @Description  Returns detailed information about a specific equipment item
// @Tags         equipment
// @Produce      json
// @Param        id   path      int  true  "Equipment ID"
// @Success      200  {object}  api.Response[api.EquipmentResponse]
// @Failure      404  {object}  api.Response[any]
// @Router       /equipment/{id} [get]
func (a *App) apiV2EquipmentGetHandler(c echo.Context) error {
	e, err := a.getEquipment(c)
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusNotFound, err)
	}

	resp := api.Response[api.EquipmentResponse]{
		Results: api.NewEquipmentResponse(e),
	}

	return c.JSON(http.StatusOK, resp)
}

// apiV2EquipmentCreateHandler creates a new equipment
// @Summary      Create equipment
// @Description  Creates a new equipment item
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Param        equipment body      database.Equipment  true  "Equipment data"
// @Success      201  {object}  api.Response[api.EquipmentResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /equipment [post]
func (a *App) apiV2EquipmentCreateHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	var e database.Equipment
	if err := c.Bind(&e); err != nil {
		return a.renderAPIV2Error(c, http.StatusBadRequest, err)
	}

	e.UserID = user.ID

	if err := e.Save(a.db); err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	resp := api.Response[api.EquipmentResponse]{
		Results: api.NewEquipmentResponse(&e),
	}

	return c.JSON(http.StatusCreated, resp)
}

// apiV2EquipmentUpdateHandler updates an existing equipment
// @Summary      Update equipment
// @Description  Updates an existing equipment item
// @Tags         equipment
// @Accept       json
// @Produce      json
// @Param        id        path      int                 true  "Equipment ID"
// @Param        equipment body      database.Equipment  true  "Equipment data"
// @Success      200  {object}  api.Response[api.EquipmentResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      403  {object}  api.Response[any]
// @Failure      404  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /equipment/{id} [put]
func (a *App) apiV2EquipmentUpdateHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	e, err := a.getEquipment(c)
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusNotFound, err)
	}

	e.DefaultFor = nil

	if e.UserID != user.ID {
		return a.renderAPIV2Error(c, http.StatusForbidden, api.ErrNotAuthorized)
	}

	if err := c.Bind(e); err != nil {
		return a.renderAPIV2Error(c, http.StatusBadRequest, err)
	}

	// Ensure user can't change ownership
	e.UserID = user.ID

	if err := e.Save(a.db); err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	resp := api.Response[api.EquipmentResponse]{
		Results: api.NewEquipmentResponse(e),
	}

	return c.JSON(http.StatusOK, resp)
}

// apiV2EquipmentDeleteHandler deletes an equipment
// @Summary      Delete equipment
// @Description  Permanently deletes an equipment item
// @Tags         equipment
// @Produce      json
// @Param        id   path      int  true  "Equipment ID"
// @Success      204
// @Failure      403  {object}  api.Response[any]
// @Failure      404  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /equipment/{id} [delete]
func (a *App) apiV2EquipmentDeleteHandler(c echo.Context) error {
	user := a.getCurrentUser(c)

	e, err := a.getEquipment(c)
	if err != nil {
		return a.renderAPIV2Error(c, http.StatusNotFound, err)
	}

	if e.UserID != user.ID {
		return a.renderAPIV2Error(c, http.StatusForbidden, api.ErrNotAuthorized)
	}

	if err := e.Delete(a.db); err != nil {
		return a.renderAPIV2Error(c, http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
