package psa

import "fmt"

func ticketConfigurationsEndpoint(ticketID int) string {
	return fmt.Sprintf("%s/configurations", ticketIdEndpoint(ticketID))
}

func configurationEndpoint(configID int) string {
	return fmt.Sprintf("company/configurations/%d", configID)
}

func (c *Client) ListTicketConfigurations(ticketID int, params map[string]string) ([]TicketConfiguration, error) {
	return GetMany[TicketConfiguration](c, ticketConfigurationsEndpoint(ticketID), params)
}

func (c *Client) GetConfiguration(id int, params map[string]string) (*Configuration, error) {
	return GetOne[Configuration](c, configurationEndpoint(id), params)
}
