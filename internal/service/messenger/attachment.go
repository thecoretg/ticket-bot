package messenger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/webex"
)

var (
	ErrConvertSubmission = errors.New("error converting submission to notifier rule")
	ErrRuleExists        = errors.New("this rule already exists")
	ErrInternalServer    = errors.New("internal server error")
)

func (s *Service) ParseAndRespondAttachment(ctx context.Context, ach *webex.AttachmentAction) error {
	slog.Debug("messenger: got attachment action", "id", ach.ID, "person_id", ach.PersonID)
	if isNotifierRule(ach) {
		return s.processNotifierRule(ctx, ach)
	}
	return nil
}

func (s *Service) processNotifierRule(ctx context.Context, ach *webex.AttachmentAction) error {
	response := &webex.Message{
		ToPersonID: ach.PersonID,
		Text:       "Rule created successfully!",
	}

	if err := s.addNotifierRule(ctx, ach); err != nil {
		switch {
		case errors.Is(err, ErrConvertSubmission):
			response.Text = "An error occured converting your submission to a rule. Please try again."
		case errors.Is(err, ErrRuleExists):
			response.Text = "A notifier rule with this board and recipient already exists."
		case errors.Is(err, ErrInternalServer):
			response.Text = "An internal server error occured. Please try again."
		default:
			response.Text = "An unknown error occured. Please try again."
		}
	}

	if _, err := s.WebexSvc.PostMessage(response); err != nil {
		return fmt.Errorf("posting webex message: %w", err)
	}

	return nil
}

func (s *Service) addNotifierRule(ctx context.Context, ach *webex.AttachmentAction) error {
	rule, err := attachmentToNotifierRule(ach)
	if err != nil {
		slog.Error("messenger: error converting submission to notifier rule", "submitted", rule, "error", err)
		return ErrConvertSubmission
	}

	exists, err := s.NotifierSvc.NotifierRules.ExistsByBoardAndRecipient(ctx, rule.CwBoardID, rule.WebexRecipientID)
	if err != nil {
		slog.Warn("messenger: error checking if notifier rule exists (continuing)",
			"board_id", rule.CwBoardID,
			"recipient_id", rule.WebexRecipientID,
			"error", err,
		)
	}

	if exists {
		slog.Info("messenger: submitted rule already exists", "board_id", rule.CwBoardID, "recipient_id", rule.WebexRecipientID)
		return ErrRuleExists
	}

	if _, err := s.NotifierSvc.AddNotifierRule(ctx, rule); err != nil {
		slog.Error("messenger: error adding notifier rule", "board_id", rule.CwBoardID, "recipient_id", rule.WebexRecipientID, "error", err)
		return ErrInternalServer
	}

	return nil
}

func attachmentToNotifierRule(ach *webex.AttachmentAction) (*models.NotifierRule, error) {
	boardIDVal := ach.Inputs["cw_board"]
	wxRecVal := ach.Inputs["webex_recipient"]
	boardID, err := strconv.Atoi(boardIDVal)
	if err != nil {
		return nil, fmt.Errorf("provided board id %s is not an integer", boardIDVal)
	}

	wxRecID, err := strconv.Atoi(wxRecVal)
	if err != nil {
		return nil, fmt.Errorf("provided webex recipient id %s is not an integer", wxRecVal)
	}

	return &models.NotifierRule{
		CwBoardID:        boardID,
		WebexRecipientID: wxRecID,
	}, nil
}

func isNotifierRule(ach *webex.AttachmentAction) bool {
	neededInputs := []string{"webex_recipient", "cw_board"}
	isRule := true
	for _, n := range neededInputs {
		if _, ok := ach.Inputs[n]; !ok {
			isRule = false
		}
	}

	return isRule
}
