package mobilevikings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.mobilevikings.be/v3"

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

type usageResponse struct {
	Results []Usage `json:"results"`
}

func (c *client) PhoneNumbers() ([]PhoneNumber, error) {
	response, err := c.doRequest("GET", "msisdns/")
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
	response, err := c.doRequest("GET", fmt.Sprintf("msisdns/%s/insights/", phoneNumber))
	if err != nil {
		return &Insights{}, err
	}
	unmarshalled := &Insights{}
	if err := json.Unmarshal(response, &unmarshalled); err != nil {
		return &Insights{}, err
	}
	return unmarshalled, nil
}

func (c *client) Usage(
	phoneNumber string,
	from time.Time,
	until time.Time,
) ([]Usage, error) {
	query := url.Values{}
	query.Set("from_date", from.Format("2006-01-02"))
	query.Set("until_date", until.Format("2006-01-02"))
	response, err := c.doRequest(
		"GET",
		fmt.Sprintf("msisdns/%s/usage/?%s", phoneNumber, query.Encode()),
	)
	if err != nil {
		return nil, err
	}
	unmarshalled := &usageResponse{}
	fmt.Println(string(response))
	if err := json.Unmarshal(response, &unmarshalled); err != nil {
		return nil, err
	}
	return unmarshalled.Results, nil
}

func (c *client) doRequest(method string, path string) ([]byte, error) {
	fullUrl := strings.Join([]string{c.baseURL, path}, "/")
	fmt.Println(fullUrl)
	request, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	request.Header.Set("Accept", "application/json")
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
