package models

type Address struct {
	Address string `json:"address"`
}

type RouteRequest struct {
	Origin        Address   `json:"origin"`
	Intermediates []Address `json:"intermediates,omitempty"`
	Destination   Address   `json:"destination"`
	TravelMode    string    `json:"travelMode"`
}

type RouteResponse struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	DistanceMeters int    `json:"distanceMeters"`
	Duration       string `json:"duration"`
	Legs           []Leg  `json:"legs"`
}

type Leg struct {
	DistanceMeters int    `json:"distanceMeters"`
	Duration       string `json:"duration"`
}

type RouteDistances struct {
	ToCurrentAddress int
	FullRoute        int
}
