package models

import "time"

type Project struct {
	ID        uint
	Name      string `validate:"nonzero"`
	OwnerId   uint   `validate:"nonzero"`
	CreatedAt time.Time
}
