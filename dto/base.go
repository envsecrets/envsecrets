package dto

import "time"

// The base structure embeded in every structure that represents a table in our database.
// This structure shouldn't be of much use to those who only intend to use the secrets structure for data transformations.
type Base struct {

	// Primary identifier of an object of this structure.
	//
	// required: true
	ID string `json:"id,omitempty"`

	// The “normalised” UTC timestamp with time zone at which this object was created.
	//
	// In postgreSQL, this is a "Timestamptz" type field.
	//
	// For example, if your input string is: 2018-08-28T12:30:00+05:30 , when this timestamp is stored in the database, it will be stored as 2018-08-28T07:00:00.
	//
	// required: true
	CreatedAt time.Time `json:"created_at,omitempty"`

	// The “normalised” UTC timestamp with time zone at which this object was last updated.
	//
	// In postgreSQL, this is a "Timestamptz" type field.
	//
	// For example, if your input string is: 2018-08-28T12:30:00+05:30 , when this timestamp is stored in the database, it will be stored as 2018-08-28T07:00:00.
	//
	// required: true
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
