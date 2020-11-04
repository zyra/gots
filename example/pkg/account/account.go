package account

import (
	prim "go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status int

const (
	StatusSuspended Status = iota + 1
	StatusActive
)

type Role int

const (
	RoleAnonymous Role = iota
	RoleUser
	RoleAdmin

	Hello = "world"
)

const RoleSomethingElse Role = 55

const RoleSomethingNotRole = 56

type Foo = int

const (
	FooFoo Foo = 0
	FooBar
)

type EmailStatus string

const EmailStatusVerified EmailStatus = "verified"
const EmailStatusUnverified EmailStatus = "unverified"

type Account struct {
	// Account ID
	ID prim.ObjectID `bson:"_id,omitempty" json:"id,omitempty" gots:"type:string"`

	// Account holder name
	Name string `json:"name"`

	// Email address
	Email string `json:"email"`

	// Email status
	EmailStatus EmailStatus `json:"emailStatus"`

	// Account role
	Role Role `json:"role"`

	// Test
	Foo Foo `json:"foo"`

	// Account status
	Status Status `json:"status"`

	// Updated at timestamp
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
