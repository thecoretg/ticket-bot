package addigy

import "fmt"

func (c *Client) SearchDevices(p *DeviceSearchParams, pageLimit int) ([]Device, error) {
	var d []Device
	if p.Page == 0 {
		p.Page = c.defaultPage
	}

	if p.PerPage == 0 {
		p.PerPage = c.defaultPerPage
	}

	for {

		resp, err := Post[DevicesResp](c, "devices", p)
		if err != nil {
			return nil, fmt.Errorf("getting devices page %d: %w", p.Page, err)
		}

		d = append(d, resp.Items...)
		if p.Page >= resp.PageCount {
			break
		}

		p.Page++

		if pageLimit != 0 && p.Page > pageLimit {
			break
		}
	}

	return d, nil
}
