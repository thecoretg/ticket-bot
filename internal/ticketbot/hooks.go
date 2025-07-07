package ticketbot

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"tctg-automation/pkg/connectwise"
	"tctg-automation/pkg/webex"
)

func (s *server) addHooksGroup(r *gin.Engine) {
	hooks := r.Group("/hooks")
	cw := hooks.Group("/cw", requireValidCWSignature(), ErrorHandler(s.exitOnError))
	cw.POST("/tickets", s.processTicketPayload)
	cw.POST("/companies", s.processCompanyPayload)
	cw.POST("/contacts", s.processContactPayload)
	cw.POST("/members", s.processMemberPayload)

	//webex := hooks.Group("/webex", ErrorHandler(s.exitOnError))

}

func (s *server) initiateAllHooks(ctx context.Context) error {
	if err := s.initiateCWHooks(ctx); err != nil {
		return fmt.Errorf("initiating connectwise hooks: %w", err)
	}

	return nil
}

func (s *server) initiateCWHooks(ctx context.Context) error {
	cwHooks, err := s.cwClient.ListCallbacks(ctx, nil)
	if err != nil {
		return fmt.Errorf("listing callbacks: %w", err)
	}

	if err := s.processCwHook(ctx, s.ticketsWebhookURL(), "ticket", "owner", 1, cwHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	if err := s.processCwHook(ctx, s.contactsWebhookURL(), "contact", "owner", 1, cwHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	if err := s.processCwHook(ctx, s.companiesWebhookURL(), "company", "owner", 1, cwHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	if err := s.processCwHook(ctx, s.membersWebhookURL(), "member", "owner", 1, cwHooks); err != nil {
		return fmt.Errorf("processing tickets hook: %w", err)
	}

	return nil
}

func (s *server) processCwHook(ctx context.Context, url, entity, level string, objectID int, currentHooks []connectwise.Callback) error {
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
				slog.Debug("found existing connectwise webhook", "id", c.ID, "entity", entity, "level", level, "url", url)
				found = true
				continue
			} else {
				if err := s.cwClient.DeleteCallback(ctx, c.ID); err != nil {
					return fmt.Errorf("deleting webhook %d: %w", c.ID, err)
				}
				slog.Debug("deleted unused connectwise webhook", "id", c.ID, "url", c.URL)
			}
		}
	}

	if !found {
		if _, err := s.cwClient.PostCallback(ctx, hook); err != nil {
			return fmt.Errorf("posting webhook: %w", err)
		}
		slog.Info("added new connectwise hook", "url", url, "entity", entity, "level", level, "objectID", objectID)
	}
	return nil
}

func (s *server) processWebexHook(ctx context.Context, name, url, resource, event string, filter *string, currentHooks []webex.Webhook) error {
	hook := &webex.Webhook{
		Name:      name,
		TargetUrl: url,
		Resource:  resource,
		Event:     event,
	}
	if filter != nil {
		hook.Filter = *filter
	}

	found := false
	for _, c := range currentHooks {
		if c.TargetUrl == hook.TargetUrl {
			if c == *hook && !found {
				slog.Debug("found existing webex webhook", "name", c.Name, "url", c.TargetUrl, "resource", c.Resource, "event", c.Event, "filter", c.Event)
				found = true
				continue
			} else {
				if err := s.webexClient.DeleteWebhook(ctx, c.ID); err != nil {
					return fmt.Errorf("deleting webhook: %w", err)
				}
				slog.Debug("deleted unused webex webhook", "id", c.ID, "url", c.TargetUrl)
			}
		}
	}

	if !found {
		if _, err := s.webexClient.CreateWebhook(ctx, hook); err != nil {
			return fmt.Errorf("posting hook: %w", err)
		}
		slog.Info("added new hook", "url", url, "resource", resource, "event", event)
	}

	return nil
}

func (s *server) ticketsWebhookURL() string {
	return fmt.Sprintf("%s/hooks/cw/tickets", s.rootUrl)
}

func (s *server) contactsWebhookURL() string {
	return fmt.Sprintf("%s/hooks/cw/contacts", s.rootUrl)
}

func (s *server) companiesWebhookURL() string {
	return fmt.Sprintf("%s/hooks/cw/companies", s.rootUrl)
}

func (s *server) membersWebhookURL() string {
	return fmt.Sprintf("%s/hooks/cw/members", s.rootUrl)
}

func requireValidCWSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(fmt.Errorf("reading request body: %w", err))
			c.Next()
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		valid, err := connectwise.ValidateWebhook(c.Request)
		if err != nil || !valid {
			c.Error(fmt.Errorf("invalid ConnectWise webhook signature: %w", err))
			c.Next()
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for further processing
		c.Next()
	}
}
