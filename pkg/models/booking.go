package models

import "time"

type Booking struct {
	Services       []string  `json:"services"`
	MovingDate     Date `json:"moving_date"`
	FlexibleDate   bool      `json:"flexible_date"`
	CleaningDate   Date `json:"cleaning_date"`
	CurrentAddress Address   `json:"current_address"`
	NewAddress     Address   `json:"new_address"`
	Contact        Contact   `json:"contact"`
}

type Address struct {
	Adress        string `json:"address"`
	ResidenceType string `json:"residence_type"`
	LivingSpace   int    `json:"living_space"`
	Accessibility string `json:"accessibility"`
	Floor         int    `json:"floor"`
}

type Contact struct {
	Name      string `json:"name"`
	SSN       string `json:"ssn"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Rutavdrag bool   `json:"rutavdrag"`
	Message   string `json:"message"`
	Policy    bool   `json:"terms_and_conditions"`
}

type Date struct {
	time.Time
}

func (cd *Date) UnmarshalJSON(b []byte) error {
	// Remove quotes
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}