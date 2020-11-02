package post

import "time"

type Post interface {
	GetPostType() string
}

// Base post properties
type BasePost struct {
	// Post title
	Title string `json:"title"`

	// Post type
	Type string `json:"type"`
}

func (p *BasePost) GetPostType() string {
	return ""
}

// Key for a post array
type Posts []*Post

// Sample type that could be used to update a post via an API
type NewPostParams struct {
	// Post object
	Post Post `json:"post"`

	// Time to publish the post on
	PublishAt time.Time `json:"publishAt,omitempty"`

	// Whether to publish the post now
	PublishNow bool `json:"publishNow,omitempty"`
}
