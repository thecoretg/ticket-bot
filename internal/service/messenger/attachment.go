package messenger

import (
	"context"
	"fmt"
	"strconv"

	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/webex"
)

func (s *Service) ParseAndRespondAttachment(ctx context.Context, ach *webex.AttachmentAction) error {
	if isNotifierRule(ach) {
		return s.processNotifierRule(ctx, ach)
	}
	return nil
}

func (s *Service) processNotifierRule(ctx context.Context, ach *webex.AttachmentAction) error {
	rule, err := attachmentToNotifierRule(ach)
	if err != nil {
		// do something
	}
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
