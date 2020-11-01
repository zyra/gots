package post

import "github.com/zyra/gots/example/pkg/image"

const DefaultBlogContent = "Hello world"

type Blog struct {
	*BasePost `json:",inline"`

	// Images used in blog post
	Images []image.Image `json:"images"`

	// Blog post content
	Content string `json:"content"`
}
