package notifications

type Subscriber struct {
	SubscriberID string                 `json:"subscriberId"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	Email        string                 `json:"email"`
	Avatar       string                 `json:"avatar"`
	Data         map[string]interface{} `json:"data"`
	Phone        string                 `json:"phone"`
}

type Trigger struct {
	EventID string                 `json:"eventId"`
	Data    map[string]interface{} `json:"data"`
	To      map[string]interface{} `json:"to"`
}

type TriggerTopic struct {
	TopicKey string `json:"topic_key"`
	EventID  string `json:"eventId"`
	Title    string `json:"title"`
	Logo     string `json:"logo"`
}
