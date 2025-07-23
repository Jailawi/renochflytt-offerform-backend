package models

import "time"

type Booking struct {
	Services       []string  `bson:"services" json:"services"`
	MovingDate     Date      `bson:"moving_date" json:"moving_date"`
	FlexibleDate   bool      `bson:"flexible_date" json:"flexible_date"`
	CleaningDate   *Date     `bson:"cleaning_date,omitempty" json:"cleaning_date,omitempty"`
	CurrentAddress Address   `bson:"current_address" json:"current_address"`
	NewAddress     Address   `bson:"new_address" json:"new_address"`
	Contact        Contact   `bson:"contact" json:"contact"`
	CreatedAt      time.Time `bson:"created_at" json:"-"`
}

type Address struct {
	Address       string `bson:"address" json:"address"`
	ResidenceType string `bson:"residence_type" json:"residence_type"`
	LivingArea    *int   `bson:"living_area,omitempty" json:"living_area,omitempty"`
	Accessibility string `bson:"accessibility" json:"accessibility"`
	Floor         *int   `bson:"floor,omitempty" json:"floor,omitempty"`
}

type Contact struct {
	Name      string  `bson:"name" json:"name"`
	SSN       string  `bson:"ssn" json:"ssn"`
	Email     string  `bson:"email" json:"email"`
	Phone     string  `bson:"phone" json:"phone"`
	Rutavdrag bool    `bson:"rutavdrag" json:"rutavdrag"`
	Message   *string `bson:"message" json:"message"`
	Consent   bool    `bson:"consent" json:"consent"`
}

type Date struct {
	time.Time `bson:"date" json:"date"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove quotes

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
