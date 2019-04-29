package imaging

import (
	"io"

	"github.com/disintegration/imaging"
	"github.com/kfreiman/imageservice/pkg/processor"
)

type imagingProcessor struct {
}

// NewProcessor ..
func NewProcessor() processor.Processor {
	return &imagingProcessor{}
}

// Modify .
func (p *imagingProcessor) Modify(income io.Reader, output io.Writer, mods processor.Modifiers) error {
	img, err := imaging.Decode(income)
	if err != nil {
		return err
	}
	converted := img
	switch mods.ResizeType {
	case processor.ResizeFit:
		converted = imaging.Fit(img, mods.Width, mods.Height, imaging.Lanczos)
	case processor.ResizeFill:
		converted = imaging.Fill(img, mods.Width, mods.Height, imaging.Center, imaging.Lanczos)
	case processor.ResizeCrop:
		converted = imaging.CropCenter(img, mods.Width, mods.Height)
	}

	return imaging.Encode(output, converted, imaging.JPEG)
}
