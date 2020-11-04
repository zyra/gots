package image

// Image type
type ImageType string

const (
	// Raw image
	ImageTypeRaw ImageType = "raw"

	// Processed image
	ImageTypeProcessed ImageType = "processed"
)

type Fooo int

const (
	Fee Fooo = iota + 1
	Barr1
)
