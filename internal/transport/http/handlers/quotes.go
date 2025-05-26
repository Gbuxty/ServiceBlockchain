package handlers

import (
	"BlockchainCurrency/config"
	"BlockchainCurrency/internal/domain"
	"BlockchainCurrency/internal/service"
	"BlockchainCurrency/pkg/logger"
	"context"
	"errors"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

type QuotesHandlers struct {
	service *service.QuotesService
	logger  logger.Logger
	cfg     *config.CredentialsAuth
}

func NewQuotesHandlers(
	service *service.QuotesService,
	logger logger.Logger,
	cfg *config.CredentialsAuth,

) *QuotesHandlers {
	return &QuotesHandlers{
		service: service,
		logger:  logger,
		cfg:     cfg,
	}
}

func (h *QuotesHandlers) GetQuotes(c *gin.Context) {
	quotes, err := h.service.GetQuotes()
	if err != nil {
		h.logger.Errorf("Failed to get tickers: %v", err)
		c.JSON(http.StatusInternalServerError, GetQuotesResponse{
			Success: false,
			Error:   "Failed to retrieve tickers",
		})
		return
	}

	c.JSON(http.StatusOK, quotes)
}

func (h *QuotesHandlers) Quit(c *gin.Context) {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.service.Shutdown(ctxShutdown); err != nil {
		if errors.Is(err, domain.ErrShutdownAlreadyCalled) {
			h.logger.Errorf("Failed to shutdown service: %v", err)
			c.JSON(http.StatusBadRequest, QuitResponse{Success: false, Message: "shutdown already called"})
			return
		}

		h.logger.Errorf("Failed to shutdown service: %v", err)
		c.JSON(http.StatusInternalServerError, QuitResponse{Success: false, Message: "failed exit"})
		return
	}

	c.JSON(http.StatusOK, QuitResponse{Success: true, Message: "Success exit"})
}
