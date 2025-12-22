package messenger

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/thecoretg/ticketbot/internal/models"
)

func createNotifierRuleList(rules []models.NotifierRuleFull) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
		"contentType": "application/vnd.microsoft.card.adaptive",
		"content": {
			"type": "AdaptiveCard",
			"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
			"version": "1.3",
			"body": [
				{
					"type": "TextBlock",
					"text": "Current Notifier Rules",
					"wrap": true,
					"weight": "Bolder"
				},
				{
					"type": "FactSet",
					"facts": [
						%s
					]
				}
			]
		}
	}`, rulesToCardEntries(rules)))
}

// sweet, man-made horrors beyond my comprehension

func createNotifierRulePayload(boards []models.Board, recips []models.WebexRecipient) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
	"contentType": "application/vnd.microsoft.card.adaptive",
	"content": {
		"type": "AdaptiveCard",
		"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
		"version": "1.3",
		"body": [
			{
				"type": "Input.ChoiceSet",
				"choices": [
					%s
				],
				"placeholder": "Pick a Connectwise board",
				"id": "cw_board",
				"label": "Connectwise Board",
				"isRequired": true,
				"errorMessage": "Connectwise board is required",
				"spacing": "none"
			},
			{
				"type": "Input.ChoiceSet",
				"choices": [
					%s
				],
				"placeholder": "Pick a Webex recipient",
				"id": "webex_recipient",
				"label": "Webex Recipient",
				"isRequired": true,
				"errormessage": "Webex recipient is required"
			},
			{
				"type": "TextBlock",
				"text": "By creating a notifier rule, you are enabling notifications for new tickets in your chosen ticket board to be sent to the webex recipient.",
				"wrap": true
			},
			{
				"type": "TextBlock",
				"text": "This also enables all users getting updated ticket notifications for tickets in this board that they are a resource of.",
				"wrap": true
			},
			{
				"type": "ActionSet",
				"actions": [
					{
						"type": "Action.Submit",
						"title": "Submit"
					}
				]
			}
		]
	}
}`, boardsToCardChoices(boards), recipientsToCardChoices(recips)))
}

func rulesToCardEntries(rules []models.NotifierRuleFull) string {
	var entries []string
	for _, r := range rules {
		title := strconv.Itoa(r.ID)
		board := r.BoardName
		rec := fmt.Sprintf("%s (%s)", r.RecipientName, r.RecipientType)
		val := fmt.Sprintf("%s > %s", board, rec)
		e := fmt.Sprintf(`{ "title": %s, "value": %s }`, strconv.Quote(title), strconv.Quote(val))
		entries = append(entries, e)
	}

	return strings.Join(entries, ",")
}

func boardsToCardChoices(boards []models.Board) string {
	var choices []string
	for _, b := range boards {
		s := fmt.Sprintf(`{ "title": %s, "value": %s }`,
			strconv.Quote(b.Name),
			strconv.Quote(strconv.Itoa(b.ID)))
		choices = append(choices, s)
	}

	return strings.Join(choices, ",")
}

func recipientsToCardChoices(recips []models.WebexRecipient) string {
	var choices []string
	for _, b := range recips {
		nameAndType := fmt.Sprintf("%s (%s)", b.Name, b.Type)
		s := fmt.Sprintf(`{ "title": %s, "value": %s }`,
			strconv.Quote(nameAndType),
			strconv.Quote(strconv.Itoa(b.ID)))
		choices = append(choices, s)
	}

	return strings.Join(choices, ",")
}
