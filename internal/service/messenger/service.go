package messenger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/internal/service/cwsvc"
	"github.com/thecoretg/ticketbot/internal/service/notifier"
	"github.com/thecoretg/ticketbot/internal/service/user"
	"github.com/thecoretg/ticketbot/internal/service/webexsvc"
	"github.com/thecoretg/ticketbot/pkg/webex"
)

var ErrNotAnAdmin = errors.New("no admin user found")

type Service struct {
	UserSvc     *user.Service
	CWSvc       *cwsvc.Service
	WebexSvc    *webexsvc.Service
	NotifierSvc *notifier.Service
}

func New(us *user.Service, cw *cwsvc.Service, wx *webexsvc.Service, ns *notifier.Service) *Service {
	return &Service{
		UserSvc:     us,
		CWSvc:       cw,
		WebexSvc:    wx,
		NotifierSvc: ns,
	}
}

func (s *Service) ParseAndRespond(ctx context.Context, msg *webex.Message) error {
	if msg == nil {
		return errors.New("got nil message")
	}

	m, err := s.parseIncoming(ctx, msg)
	if err != nil {
		return fmt.Errorf("parsing incoming message: %w", err)
	}

	if m == nil {
		return errors.New("parse returned nil message")
	}

	if _, err := s.WebexSvc.PostMessage(m); err != nil {
		return fmt.Errorf("posting webex message: %w", err)
	}
	return nil
}

func (s *Service) parseIncoming(ctx context.Context, msg *webex.Message) (*webex.Message, error) {
	email := msg.PersonEmail
	if email == "" {
		return nil, errors.New("got empty email field")
	}

	if _, err := s.getValidUser(ctx, email); err != nil {
		if errors.Is(err, ErrNotAnAdmin) {
			return notAnAdminMessage(email), nil
		}
	}

	txt := msg.Text
	if msg.Markdown != "" {
		txt = msg.Markdown
	}
	slog.Info("messenger: got message", "sender", email, "text", txt)

	switch txt {
	case "list rules":
		slog.Debug("messenger: message matches list rules message", "sender", email, "text", txt)
		return s.makeListRulesMsg(ctx, email)
	case "create rule":
		slog.Debug("messenger: message matches create rule message", "sender", email, "text", txt)
		return s.makeCreateRuleMsg(ctx, email)
	default:
		slog.Warn("messenger: received invalid command", "sender", email, "text", txt)
		return invalidCommandMessage(email), nil
	}
}

func (s *Service) makeListRulesMsg(ctx context.Context, email string) (*webex.Message, error) {
	m := webex.NewMessageToPerson(email, "test")
	msg := &m

	rules, err := s.NotifierSvc.ListNotifierRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing notifier rules: %w", err)
	}
	slog.Debug("makeListRulesMsg: got rules", "total", len(rules))

	if len(rules) == 0 {
		msg.Markdown = "No notifier rules found."
		return msg, nil
	}

	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].BoardName < rules[j].BoardName
	})

	msg.Attachments = []json.RawMessage{createNotifierRuleList(rules)}
	return msg, nil
}

func (s *Service) makeCreateRuleMsg(ctx context.Context, email string) (*webex.Message, error) {
	boards, err := s.CWSvc.ListBoards(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing cw boards: %w", err)
	}
	slog.Debug("makeCreateRuleMsg: got boards", "total", len(boards))
	sort.SliceStable(boards, func(i, j int) bool {
		return boards[i].Name < boards[j].Name
	})

	recips, err := s.WebexSvc.ListRecipients(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing webex recipients: %w", err)
	}
	slog.Debug("makeCreateRuleMsg: got webex recipients", "total", len(recips))
	sort.SliceStable(recips, func(i, j int) bool {
		return recips[i].Name < recips[j].Name
	})

	m := webex.NewMessageToPerson(email, "test")
	msg := &m
	msg.Attachments = []json.RawMessage{createNotifierRulePayload(boards, recips)}
	return msg, nil
}

func (s *Service) getValidUser(ctx context.Context, email string) (*models.APIUser, error) {
	u, err := s.UserSvc.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, models.ErrAPIUserNotFound) {
			return nil, ErrNotAnAdmin
		}
		return nil, fmt.Errorf("getting user by email: %w", err)
	}

	return u, nil
}

func notAnAdminMessage(email string) *webex.Message {
	txt := "Sorry, this command requires admin permissions."
	m := webex.NewMessageToPerson(email, txt)
	return &m
}

func invalidCommandMessage(email string) *webex.Message {
	txt := "That was in invalid command."
	m := webex.NewMessageToPerson(email, txt)
	return &m
}
