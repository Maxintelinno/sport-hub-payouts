package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sport-hub/sport-hub-payouts/internal/service"
)

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(service service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) GetSummary(c echo.Context) error {
	// Assuming owner_id is set in context by an auth middleware
	// For now, let's extract it from a header or a placeholder
	ownerID := c.Get("owner_id").(string)
	if ownerID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	summary, err := h.service.GetSummary(c.Request().Context(), ownerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, summary)
}
