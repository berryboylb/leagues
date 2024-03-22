package teams

type TeamRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Country     string `json:"country" binding:"required,min=3"`
	State       string `json:"state" binding:"required,min=3"`
	FoundedYear int    `json:"founded_year" binding:"required"`
	Stadium     string `json:"stadium" binding:"required,min=3"`
	Sponsor     string `json:"sponsor" binding:"required,min=3"`
}

type TeamFilterRequest struct{}
