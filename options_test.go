package mjml

import (
	"reflect"
	"testing"
)

func TestOptions(t *testing.T) {

	beautifyOptions := NewBeautifyOptions().IndentEmptyLines(true).WrapLineLength(10)
	juiceOptions := NewJuiceOptions().PreserveKeyFrames(true).PreservePseudos(true)
	htmlMinifierOptions := NewHTMLMinifierOptions().HTML5(true).MinifyURLs(true)

	o := options{data: map[string]interface{}{}}

	optionFunctions := []ToHTMLOption{
		WithBeautify(true),
		WithBeautifyOptions(beautifyOptions),
		WithFonts(Fonts{"test": "test"}),
		WithJuiceOptions(juiceOptions),
		WithJuicePreserveTags(
			map[string]JuiceTag{
				"myTag": {
					Start: "<#",
					End:   "</#",
				},
			},
		),
		WithKeepComments(true),
		WithMinify(true),
		WithMinifyOptions(htmlMinifierOptions),
		WithPreprocessors([]string{"(xml) => xml"}),
		WithValidationLevel(Strict),
	}

	for _, f := range optionFunctions {
		f(o)
	}

	expected := map[string]interface{}{
		"beautify": true,
		"beautifyOptions": map[string]interface{}{
			"indent_empty_lines": true,
			"wrap_line_length":   uint(10),
		},
		"fonts": Fonts{"test": "test"},
		"juiceOptions": map[string]interface{}{
			"preserveKeyFrames": true,
			"preservePseudos":   true,
		},
		"juicePreserveTags": map[string]JuiceTag{
			"myTag": {
				Start: "<#",
				End:   "</#",
			},
		},
		"keepComments": true,
		"minify":       true,
		"minifyOptions": map[string]interface{}{
			"html5":      true,
			"minifyURLs": true,
		},
		"preprocessors":   []string{"(xml) => xml"},
		"validationLevel": Strict,
	}

	if !reflect.DeepEqual(o.data, expected) {
		t.Error("Options does not match expected data")
	}
}
