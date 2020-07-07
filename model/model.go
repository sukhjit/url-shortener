package model

// Shortener struct
type Shortener struct {
	Slug   string `json:"slug"`
	URL    string `json:"url"`
	Visits int    `json:"visits"`
}
