package ticketbot

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"tctg-automation/pkg/connectwise"
)

func (s *server) initiateTicketWebhook(ctx context.Context) error {
	currentHooks, err := s.cwClient.ListCallbacks(ctx, nil)
	if err != nil {
		return fmt.Errorf("listing callbacks: %w", err)
	}

	if err := s.processHook(ctx, s.ticketsWebhookURL(), "ticket", "owner", 1, currentHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	if err := s.processHook(ctx, s.contactsWebhookURL(), "contact", "owner", 1, currentHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	if err := s.processHook(ctx, s.companiesWebhookURL(), "company", "owner", 1, currentHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	if err := s.processHook(ctx, s.membersWebhookURL(), "member", "owner", 1, currentHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	return nil
}

func (s *server) processHook(ctx context.Context, url, entity, level string, objectID int, currentHooks []connectwise.Callback) error {
	hook := &connectwise.Callback{
		URL:      url,
		Type:     entity,
		Level:    level,
		ObjectId: objectID,
	}

	found := false
	for _, c := range currentHooks {
		if c.URL == hook.URL {
			if c.Type == hook.Type && c.Level == hook.Level && c.InactiveFlag == hook.InactiveFlag && !found {
				log.Printf("found existing webhook %d with URL %s", c.ID, c.URL)
				found = true
				continue
			} else {
				if err := s.cwClient.DeleteCallback(ctx, c.ID); err != nil {
					slog.Error("deleting unneeded hook", "url", url, "entity", entity, "level", level, "objectID", objectID)
					return fmt.Errorf("deleting webhook %d: %w", c.ID, err)
				}
				slog.Info("deleted unused webhook", "id", c.ID, "url", c.URL)
			}
		}
	}

	if !found {
		if _, err := s.cwClient.PostCallback(ctx, hook); err != nil {
			return fmt.Errorf("posting webhook: %w", err)
		}
		slog.Info("added new hook", "url", url, "entity", entity, "level", level, "objectID", objectID)
	}
	return nil
}

func (s *server) ticketsWebhookURL() string {
	return fmt.Sprintf("%s/tickets", s.rootUrl)
}

func (s *server) contactsWebhookURL() string {
	return fmt.Sprintf("%s/contacts", s.rootUrl)
}

func (s *server) companiesWebhookURL() string {
	return fmt.Sprintf("%s/companies", s.rootUrl)
}

func (s *server) membersWebhookURL() string {
	return fmt.Sprintf("%s/members", s.rootUrl)
}
