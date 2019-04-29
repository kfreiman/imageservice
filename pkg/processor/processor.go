package processor

import (
	"io"
)

// Processor describes processor abilities
type Processor interface {
	Modify(income io.Reader, output io.Writer, modifiers Modifiers) error
}

// ResizeType .
type ResizeType int

// ResizeType
const (
	ResizeFill ResizeType = iota
	ResizeFit
	ResizeCrop
)

// Modifiers defines how to modify original image
type Modifiers struct {
	Width      int
	Height     int
	ResizeType ResizeType
	// Quality, Watermark, Gravity etc...
}
