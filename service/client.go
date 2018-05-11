package service

import (
	"github.com/molsbee/go-cfclient"
	"net/url"
	"fmt"
)

type Client struct {
	*cfclient.Client
}

func NewClient(endpoint, username, password string) (*Client, error) {
	cfClient, err := cfclient.NewClient(&cfclient.Config{
		ApiAddress: endpoint,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	return &Client{cfClient}, nil
}

func (c *Client) ListAppsByOrgGuid(guid string) ([]cfclient.App, error) {
	return c.ListAppsByQuery(url.Values{
		"q": {fmt.Sprintf("organization_guid:%s", guid)},
	})
}

func (c *Client) ListServiceInstancesByOrgGuid(guid string) ([]cfclient.Service, error) {
	return c.ListServicesByQuery(url.Values{
		"q": {fmt.Sprintf("organization_guid:%s", guid)},
	})
}

func (c *Client) ListAppServiceBindings(appGuid string) ([]cfclient.ServiceBinding, error) {
	return c.ListServiceBindingsByQuery(url.Values{
		"q": {fmt.Sprintf("app_guid:%s", appGuid)},
	})
}