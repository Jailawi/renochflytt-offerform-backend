package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID               primitive.ObjectID `bson:"_id," json:"id"`
	Services         []string           `bson:"services" json:"services"`
	MovingDate       Date               `bson:"moving_date" json:"moving_date"`
	IsFlexibleDate   bool               `bson:"is_flexible_date" json:"is_flexible_date"`
	EstimatedPrice   int                `bson:"estimated_price" json:"estimated_price"`
	CleaningDate     *Date              `bson:"cleaning_date,omitempty" json:"cleaning_date,omitempty"`
	CurrentResidence Residence          `bson:"current_address" json:"current_address"`
	NewResidence     Residence          `bson:"new_address" json:"new_address"`
	Contact          Contact            `bson:"contact" json:"contact"`
	EmailSent        bool               `bson:"email_sent" json:"email_sent"`
	CreatedAt        time.Time          `bson:"created_at" json:"-"`
}

type Residence struct {
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
