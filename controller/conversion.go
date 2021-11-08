package controller

import (
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mhaqiw/game-currency/domain"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type ConversionHandler struct {
	ConversionService domain.ConversionService
}

func NewConversionHandler(e *echo.Echo, m domain.ConversionService) {
	handler := &ConversionHandler{
		ConversionService: m,
	}
	e.GET("/conversion/rate/:from/:to/:amount", handler.GetRate)
	e.POST("/conversion", handler.Post)
}

func (h *ConversionHandler) Post(c echo.Context) error {
	var request domain.ConversionInsertRequestPayload
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

	resp, err := h.ConversionService.Create(c.Request().Context(), request)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *ConversionHandler) GetRate(c echo.Context) error {
	fromS := c.Param("from")
	from, err := strconv.ParseInt(fromS, 10, 64)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	toS := c.Param("to")
	to, err := strconv.ParseInt(toS, 10, 64)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	amountS := c.Param("amount")
	amount, err := strconv.ParseFloat(amountS, 64)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	request := domain.ConvertedCurrencyRequest{
		FromID: from,
		ToID:   to,
		Amount: amount,
	}

	response, err := h.ConversionService.GetConvertedCurrency(c.Request().Context(), request)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, domain.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}
