package account

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Account struct {
	// Account ID
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" gots:"type:string"`

	// Account holder name
	Name string `json:"name"`

	// Email address
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
