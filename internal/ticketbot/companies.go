package ticketbot

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tctg-automation/pkg/connectwise"
	"tctg-automation/pkg/util"
)

func (s *server) processCompanyPayload(c *gin.Context) {
	w := &connectwise.WebhookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("invalid request body"))
		return
	}

	if w.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company ID cannot be 0"})
		return
	}
	switch w.Action {
	case "deleted":
		if err := s.dbHandler.DeleteCompany(w.ID); err != nil {
			slog.Error("deleting company", "id", w.ID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   fmt.Sprintf("couldn't delete company %v", err),
				"company": w.ID,
			})
			return
		}

		slog.Info("company deleted", "id", w.ID)
		c.Status(http.StatusNoContent)
		return
	default:
		if err := processCompanyUpdate(c.Request.Context(), w.ID); err != nil {

		}
	}
}

func (s *server) processCompanyUpdate(ctx context.Context, companyID int) error {
	cwc, err := s.cwClient.GetCompany(ctx, companyID, nil)
	if err != nil {
		return checkCWError("getting company via CW API", err, companyID)
	}

	c := NewCompany(companyID, cwc.Name)
	if err := s.dbHandler.UpsertCompany(c); err != nil {
		return fmt.Errorf("processing company in db: %w", err)
	}

	return nil
}
