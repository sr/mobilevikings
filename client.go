package mobilevikings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultBaseURL = "https://api.mobilevikings.be/v3/%s/"

func newClient(accessToken string) *client {
	return &client{defaultBaseURL, accessToken, &http.Client{}}
}

type client struct {
	baseURL     string
	accessToken string
	client      *http.Client
}

type phoneNumbersResponse struct {
	Results []PhoneNumber `json:"results"`
}

func (c *client) PhoneNumbers() ([]PhoneNumber, error) {
	response, err := c.doRequest("GET", "msisdns")
	if err != nil {
		return nil, err
	}
	unmarshalled := &phoneNumbersResponse{}
	if err := json.Unmarshal(response, &unmarshalled); err != nil {
		return nil, err
	}
	return unmarshalled.Results, nil
}

func (c *client) Insights(phoneNumber string) (*Insights, error) {
	response, err := c.doRequest("GET", fmt.Sprintf("msisdns/%s/insights", phoneNumber))
	if err != nil {
		return &Insights{}, err
	}
	unmarshalled := &Insights{}
	if err := json.Unmarshal(response, &unmarshalled); err != nil {
		return &Insights{}, err
	}
	return unmarshalled, nil
}

func (c *client) doRequest(method string, path string) ([]byte, error) {
	fullUrl := fmt.Sprintf(c.baseURL, path)
	request, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}