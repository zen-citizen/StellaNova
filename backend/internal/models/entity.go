package models

type EntitiesRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	City      string  `json:"city" binding:"required"`
}

type EntitiesResponse struct {
	Entities []Entity `json:"entities"`
}

type Entity struct {
	Name                string      `json:"name"`
	IsAvailable         bool        `json:"is_available"`
	NotAvailableMessage string      `json:"not_available_message"`
	Disclaimer          *string     `json:"disclaimer,omitempty"`
	Attributes          []Attribute `json:"attributes"`
}

type Attribute struct {
	Name    string   `json:"name"`
	Value   string   `json:"value"`
	Address *Address `json:"address,omitempty"`
	IsFound bool     `json:"-"`
}

type Address struct {
	Text string  `json:"text"`
	Link *string `json:"link,omitempty"`
}

func NewUnavailableEntity(name, message string, disclaimer *string) Entity {
	return Entity{
		Name:                name,
		IsAvailable:         false,
		NotAvailableMessage: message,
		Disclaimer:          disclaimer,
		Attributes:          []Attribute{},
	}
}

func NewAvailableEntity(name, message string, disclaimer *string, attributes []Attribute) Entity {
	return Entity{
		Name:                name,
		IsAvailable:         true,
		NotAvailableMessage: message,
		Disclaimer:          disclaimer,
		Attributes:          attributes,
	}
}
