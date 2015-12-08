package mobilevikings

import (
	"strconv"
	"time"
)

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
	Type        string `json:"type"`
	PriceString string `json:"price"`
	Length      int64  `json:"lenght"`
	Number      string `json:"number"`
}

func (u Usage) Price() (int64, error) {
	matches := usageAmount.FindStringSubmatch(u.PriceString)
	i, err := strconv.ParseInt(matches[1], 10, 64)
	f, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, err
	}
	return i + (f * 100), nil
}

func NewClient(accessToken string) Client {
	return newClient(accessToken)
}
