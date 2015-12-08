package mobilevikings

import "time"

type Client interface {
	PhoneNumbers() ([]PhoneNumber, error)
	Insights(phoneNumber string) (*Insights, error)
	Usage(phoneNumber string, from time.Time, until time.Time) ([]Usage, error)
}

type PhoneNumber struct {
	ID    string `json:"msisdn"`
	Alias string `json:"alias"`
}

type Insights struct {
	VikingLife VikingLife `json:"viking_life"`
}

type VikingLife struct {
	DaysAsViking int `json:"days_as_a_viking"`
}

type Usage struct {
	Type    string `json:"type"`
	Price   string `json:"price"`
	Numbert string `json:"number"`
}

func NewClient(accessToken string) Client {
	return newClient(accessToken)
}
