package mjml

import "fmt"

type ValidationLevel string

const (
	Strict ValidationLevel = "strict"
	Soft   ValidationLevel = "soft"
	Skip   ValidationLevel = "skip"
)

type JuiceTag struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type options struct {
	data map[string]interface{}
}

type Fonts map[string]string

// ToHTMLOption provides options to customize the compilation process
// Detailed explanations of each option is available here: https://github.com/mjmlio/mjml#inside-nodejs
type ToHTMLOption func(options)

func WithBeautify(beautify bool) ToHTMLOption {
	return func(o options) {
		o.data["beautify"] = beautify
	}
}

func WithBeautifyOptions(bOptions BeautifyOptions) ToHTMLOption {
	beautifyOptions, ok := bOptions.(*beautifyOptions)

	if !ok {
		panic(fmt.Errorf("unsupported BeautifyOptions implementation: %#v", beautifyOptions))
	}

	return func(o options) {
		o.data["beautifyOptions"] = beautifyOptions.data
	}
}

func WithFonts(fonts Fonts) ToHTMLOption {
	return func(o options) {
		o.data["fonts"] = fonts
	}
}

func WithJuiceOptions(jOptions JuiceOptions) ToHTMLOption {
	juiceOptions, ok := jOptions.(*juiceOptions)

	if !ok {
		panic(fmt.Errorf("unsupported JuiceOptions implementation: %#v", juiceOptions))
	}

	return func(o options) {
		o.data["juiceOptions"] = juiceOptions.data
	}
}

func WithJuicePreserveTags(preserveTags map[string]JuiceTag) ToHTMLOption {
	return func(o options) {
		o.data["juicePreserveTags"] = preserveTags
	}
}

func WithKeepComments(keepComments bool) ToHTMLOption {
	return func(o options) {
		o.data["keepComments"] = keepComments
	}
}

func WithMinify(minify bool) ToHTMLOption {
	return func(o options) {
		o.data["minify"] = minify
	}
}

func WithMinifyOptions(minifyOptions HTMLMinifierOptions) ToHTMLOption {
	htmlMinifierOptions, ok := minifyOptions.(*htmlMinifierOptions)

	if !ok {
		panic(fmt.Errorf("unsupported HTMLMinifierOptions implementation: %#v", htmlMinifierOptions))
	}

	return func(o options) {
		o.data["minifyOptions"] = htmlMinifierOptions.data
	}
}

func WithPreprocessors(preprocessors []string) ToHTMLOption {
	return func(o options) {
		o.data["preprocessors"] = preprocessors
	}
}

func WithValidationLevel(validationLevel ValidationLevel) ToHTMLOption {
	return func(o options) {
		o.data["validationLevel"] = validationLevel
	}
}
