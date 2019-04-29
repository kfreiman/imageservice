package handler

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kfreiman/imageservice/pkg/processor"

	"net/http"
	"net/url"
	"strconv"
)

// Extractor analyzes the HTTP request and determines how to modify the image.
// Extractor's implementation defines the external API design.
// It can use GET or POST params, HTTP-headers, or even user-agent.
type Extractor interface {
	Extract(r *http.Request) (processor.Modifiers, error)
}

var resizeTypes = map[string]processor.ResizeType{
	"fit":  processor.ResizeFit,
	"fill": processor.ResizeFill,
	"crop": processor.ResizeCrop,
}

type extractor struct {
}

// NewExtractor returns Extractor's implementation based on GET params
func NewExtractor() Extractor {
	return &extractor{}
}

// Extract parses the query and puts the result to processor.Modifiers
func (e *extractor) Extract(r *http.Request) (processor.Modifiers, error) {
	var mods processor.Modifiers

	err := e.applyModifiers(&mods, r.URL.Query())
	if err != nil {
		return mods, &BadRequestError{err}
	}

	// even if extraction is ok, values should be validated
	err = e.validate(mods)
	if err != nil {
		return mods, &ValidationError{err}
	}
	return mods, nil
}

func (e *extractor) applyModifiers(modifiers *processor.Modifiers, params url.Values) error {
	for name, values := range params {
		if len(values) > 1 {
			return fmt.Errorf("Invalid %s arguments: %v", name, values)
		}

		if err := e.applyModifier(modifiers, name, values[0]); err != nil {
			return err
		}
	}
	return nil
}

func (e *extractor) applyModifier(mods *processor.Modifiers, name string, value string) error {
	switch name {
	case "width":
		if i, err := strconv.Atoi(value); err == nil {
			mods.Width = i
			return nil
		}
	case "height":
		if i, err := strconv.Atoi(value); err == nil {
			mods.Height = i
			return nil
		}
	case "resizing_type":
		if r, ok := resizeTypes[value]; ok {
			mods.ResizeType = r
			return nil
		}
	default:
		return nil // ignore unknown params
	}

	return fmt.Errorf("Invalid %s: %s", name, value)
}

func (e *extractor) validate(mods processor.Modifiers) error {
	err := validation.ValidateStruct(&mods,
		validation.Field(&mods.Width, validation.Required, validation.Min(1)),
		validation.Field(&mods.Height, validation.Required, validation.Min(1)),
		validation.Field(&mods.ResizeType,
			validation.In(
				processor.ResizeFill,
				processor.ResizeFit,
				processor.ResizeCrop,
			)),
	)

	if err != nil {
		return err
	}
	return nil
}
