package mjml

import (
	"reflect"
	"testing"
)

func TestJuiceOptions(t *testing.T) {

	options := NewJuiceOptions().
		ApplyAttributesTableElements(true).
		ApplyHeightAttributes(true).
		ApplyStyleTags(true).
		ApplyWidthAttributes(true).
		ExtraCss("somestring").
		InsertPreservedExtraCss(true).
		InlinePseudoElements(true).
		PreserveFontFaces(true).
		PreserveImportant(true).
		PreserveMediaQueries(true).
		PreserveKeyFrames(true).
		PreservePseudos(true).
		RemoveStyleTags(true).
		XmlMode(true)

	expected := map[string]interface{}{
		"applyAttributesTableElements": true,
		"applyHeightAttributes":        true,
		"applyStyleTags":               true,
		"applyWidthAttributes":         true,
		"extraCss":                     "somestring",
		"insertPreservedExtraCss":      true,
		"inlinePseudoElements":         true,
		"preserveFontFaces":            true,
		"preserveImportant":            true,
		"preserveMediaQueries":         true,
		"preserveKeyFrames":            true,
		"preservePseudos":              true,
		"removeStyleTags":              true,
		"xmlMode":                      true,
	}

	juiceOptions, ok := options.(*juiceOptions)

	if !ok {
		t.Fatal("Options is not a *juiceOptions")
	}

	if !reflect.DeepEqual(juiceOptions.data, expected) {
		t.Error("JuiceOptions does not match expected data")
	}
}
