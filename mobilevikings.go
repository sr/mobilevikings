package mobilevikings

type Client interface {
	PhoneNumbers() ([]PhoneNumber, error)
	Insights(phoneNumber string) (*Insights, error)
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

func NewClient(accessToken string) Client {
	return newClient(accessToken)
}
