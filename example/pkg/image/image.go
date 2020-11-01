package image

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Default image type
var DefaultImageType = ImageTypeRaw

// Image object
type Image struct {
	// Unique image ID
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	// Image owner's account ID
	Owner primitive.ObjectID `json:"owner,omitempty"`

	// Image type
	Type ImageType `json:"type,omitempty"`

	UpdatedAt time.Time `json:"-"`
}
