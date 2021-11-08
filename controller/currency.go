package controller

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mhaqiw/game-currency/domain"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type CurrencyHandler struct {
	CurrencyService domain.CurrencyService
}

func NewCurrencyHandler(e *echo.Echo, m domain.CurrencyService) {
	handler := &CurrencyHandler{
		CurrencyService: m,
	}

	e.GET("/currency", handler.GetAll)
	e.POST("/currency", handler.Post)
}

func (h *CurrencyHandler) GetAll(c echo.Context) error {
	response, err := h.CurrencyService.GetAll(c.Request().Context())
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CurrencyHandler) Post(c echo.Context) error {
	var request domain.CurrencyRequestPayload
	err := c.Bind(&request)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusUnprocessableEntity, domain.Response{Message: err.Error()})
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	resp, err := h.CurrencyService.Create(c.Request().Context(), request)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}
