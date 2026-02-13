package controller

import (
	"net/http"
	"time"

	"github.com/jovandeginste/workout-tracker/v2/pkg/api"
	"github.com/jovandeginste/workout-tracker/v2/pkg/container"
	"github.com/jovandeginste/workout-tracker/v2/pkg/database"
	"github.com/jovandeginste/workout-tracker/v2/pkg/templatehelpers"
	"github.com/labstack/echo/v4"
)

type MeasurementController interface {
	GetMeasurements(c echo.Context) error
	CreateMeasurement(c echo.Context) error
	DeleteMeasurement(c echo.Context) error
}

type measurementController struct {
	context *container.Container
}

func NewMeasurementController(c *container.Container) MeasurementController {
	return &measurementController{context: c}
}

// GetMeasurements returns a paginated list of measurements for the current user
// @Summary      List measurements
// @Tags         measurements
// @Security     ApiKeyAuth
// @Security     ApiKeyQuery
// @Security     CookieAuth
// @Param        page      query     int false "Page"
// @Param        per_page  query     int false "Per page"
// @Produce      json
// @Success      200  {object}  api.PaginatedResponse[api.MeasurementResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /measurements [get]
func (mc *measurementController) GetMeasurements(c echo.Context) error {
	user := mc.context.GetUser(c)

	var pagination api.PaginationParams
	if err := c.Bind(&pagination); err != nil {
		return renderApiError(c, http.StatusBadRequest, err)
	}
	pagination.SetDefaults()

	var totalCount int64
	if err := mc.context.GetDB().Model(&database.Measurement{}).Where("user_id = ?", user.ID).Count(&totalCount).Error; err != nil {
		return renderApiError(c, http.StatusInternalServerError, err)
	}

	var measurements []*database.Measurement
	db := mc.context.GetDB().Where("user_id = ?", user.ID).
		Order("date DESC").
		Limit(pagination.PerPage).
		Offset(pagination.GetOffset())

	if err := db.Find(&measurements).Error; err != nil {
		return renderApiError(c, http.StatusInternalServerError, err)
	}

	results := api.NewMeasurementsResponse(measurements)

	resp := api.PaginatedResponse[api.MeasurementResponse]{
		Results:    results,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		TotalPages: pagination.CalculateTotalPages(totalCount),
		TotalCount: totalCount,
	}

	return c.JSON(http.StatusOK, resp)
}

// CreateMeasurement creates or updates a measurement for a specific date
// @Summary      Create or update measurement
// @Tags         measurements
// @Security     ApiKeyAuth
// @Security     ApiKeyQuery
// @Security     CookieAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.Response[api.MeasurementResponse]
// @Failure      400  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /measurements [post]
func (mc *measurementController) CreateMeasurement(c echo.Context) error {
	user := mc.context.GetUser(c)

	d := &measurementPayload{units: user.PreferredUnits()}
	if err := c.Bind(d); err != nil {
		return renderApiError(c, http.StatusBadRequest, err)
	}

	m, err := user.GetMeasurementForDate(d.Time())
	if err != nil {
		return renderApiError(c, http.StatusInternalServerError, err)
	}

	d.Update(m)

	if err := m.Save(mc.context.GetDB()); err != nil {
		return renderApiError(c, http.StatusInternalServerError, err)
	}

	resp := api.Response[api.MeasurementResponse]{
		Results: api.NewMeasurementResponse(m),
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteMeasurement deletes a measurement for a specific date
// @Summary      Delete measurement by date
// @Tags         measurements
// @Security     ApiKeyAuth
// @Security     ApiKeyQuery
// @Security     CookieAuth
// @Param        date  path  string  true  "Date (YYYY-MM-DD)"
// @Produce      json
// @Success      204  {string}  string ""
// @Failure      400  {object}  api.Response[any]
// @Failure      404  {object}  api.Response[any]
// @Router       /measurements/{date} [delete]
func (mc *measurementController) DeleteMeasurement(c echo.Context) error {
	u := mc.context.GetUser(c)

	dateStr := c.Param("date")
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return renderApiError(c, http.StatusBadRequest, err)
	}

	m, err := u.GetMeasurementForDate(t)
	if err != nil {
		return renderApiError(c, http.StatusNotFound, err)
	}

	if err := m.Delete(mc.context.GetDB()); err != nil {
		return renderApiError(c, http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

type measurementPayload struct {
	Date       string  `form:"date" json:"date"`
	Steps      float64 `form:"steps" json:"steps"`
	WeightUnit string  `form:"weight_unit" json:"weight_unit"`
	HeightUnit string  `form:"height_unit" json:"height_unit"`

	Weight           float64 `form:"weight" json:"weight"`
	Height           float64 `form:"height" json:"height"`
	FTP              float64 `form:"ftp" json:"ftp"`
	RestingHeartRate float64 `form:"resting_heart_rate" json:"resting_heart_rate"`
	MaxHeartRate     float64 `form:"max_heart_rate" json:"max_heart_rate"`

	units *database.UserPreferredUnits
}

func (m *measurementPayload) Time() time.Time {
	if m.Date == "" {
		return time.Now()
	}

	d, err := time.Parse("2006-01-02", m.Date)
	if err != nil {
		return time.Now()
	}

	return d
}

func (m *measurementPayload) ToSteps() *float64 {
	if m.Steps == 0 {
		return nil
	}

	d := m.Steps
	return &d
}

func (m *measurementPayload) ToFTP() *float64 {
	if m.FTP == 0 {
		return nil
	}

	d := m.FTP
	return &d
}

func (m *measurementPayload) ToRestingHeartRate() *float64 {
	if m.RestingHeartRate == 0 {
		return nil
	}

	d := m.RestingHeartRate
	return &d
}

func (m *measurementPayload) ToMaxHeartRate() *float64 {
	if m.MaxHeartRate == 0 {
		return nil
	}

	d := m.MaxHeartRate
	return &d
}

func (m *measurementPayload) ToHeight() *float64 {
	if m.Height == 0 {
		return nil
	}

	if m.HeightUnit == "" {
		m.HeightUnit = m.units.Height()
	}

	d := templatehelpers.HeightToDatabase(m.Height, m.HeightUnit)
	return &d
}

func (m *measurementPayload) ToWeight() *float64 {
	if m.Weight == 0 {
		return nil
	}

	if m.WeightUnit == "" {
		m.WeightUnit = m.units.Weight()
	}

	d := templatehelpers.WeightToDatabase(m.Weight, m.WeightUnit)
	return &d
}

func (m *measurementPayload) Update(measurement *database.Measurement) {
	setIfNotNil(&measurement.Weight, m.ToWeight())
	setIfNotNil(&measurement.Height, m.ToHeight())
	setIfNotNil(&measurement.Steps, m.ToSteps())
	setIfNotNil(&measurement.FTP, m.ToFTP())
	setIfNotNil(&measurement.RestingHeartRate, m.ToRestingHeartRate())
	setIfNotNil(&measurement.MaxHeartRate, m.ToMaxHeartRate())
}
