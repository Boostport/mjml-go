package mjml

type BeautifyBraceStyle string

const (
	BeautifyBraceStyleCollapsePreserveInline BeautifyBraceStyle = "collapse-preserve-inline"
	BeautifyBraceStyleCollapse               BeautifyBraceStyle = "collapse"
	BeautifyBraceStyleExpand                 BeautifyBraceStyle = "expand"
	BeautifyBraceStyleEndExpand              BeautifyBraceStyle = "end-expand"
	BeautifyBraceStyleNone                   BeautifyBraceStyle = "none"
)

type BeautifyIndentScripts string

const (
	BeautifyIndentScriptsKeep     BeautifyIndentScripts = "keep"
	BeautifyIndentScriptsSeparate BeautifyIndentScripts = "separate"
	BeautifyIndentScriptsNormal   BeautifyIndentScripts = "normal"
)

type BeautifyWrapAttributes string

const (
	BeautifyWrapAttributesAuto                 BeautifyWrapAttributes = "auto"
	BeautifyWrapAttributesForce                BeautifyWrapAttributes = "force"
	BeautifyWrapAttributesForceAligned         BeautifyWrapAttributes = "force-aligned"
	BeautifyWrapAttributesForceExpandMultiline BeautifyWrapAttributes = "force-expand-multiline"
	BeautifyWrapAttributesAlignedMultiple      BeautifyWrapAttributes = "aligned-multiple"
	BeautifyWrapAttributesPreserve             BeautifyWrapAttributes = "preserve"
	BeautifyWrapAttributesPreserveAligned      BeautifyWrapAttributes = "preserved-aligned"
)

type BeautifyTemplating string

const (
	BeautifyTemplatingAuto       BeautifyTemplating = "auto"
	BeautifyTemplatingNone       BeautifyTemplating = "none"
	BeautifyTemplatingDjango     BeautifyTemplating = "django"
	BeautifyTemplatingERB        BeautifyTemplating = "erb"
	BeautifyTemplatingHandlebars BeautifyTemplating = "handlebars"
	BeautifyTemplatingPHP        BeautifyTemplating = "php"
	BeautifyTemplatingSmarty     BeautifyTemplating = "smarty"
)

// BeautifyOptions is used to construct Beautify options to be passed to the MJML compiler
// Detailed explanations of the options are here: https://github.com/beautify-web/js-beautify#css--html
type BeautifyOptions interface {
	IndentSize(uint) BeautifyOptions
	IndentChar(string) BeautifyOptions
	IndentWithTabs(bool) BeautifyOptions
	Eol(string) BeautifyOptions
	EndWithNewline(bool) BeautifyOptions
	PreserveNewlines(bool) BeautifyOptions
	MaxPreserveNewlines(uint) BeautifyOptions
	IndentInnerHtml(bool) BeautifyOptions
	BraceStyle(BeautifyBraceStyle) BeautifyOptions
	IndentScripts(BeautifyIndentScripts) BeautifyOptions
	WrapLineLength(uint) BeautifyOptions
	WrapAttributes(BeautifyWrapAttributes) BeautifyOptions
	WrapAttributesIndentSize(uint) BeautifyOptions
	Inline([]string) BeautifyOptions
	Unformatted([]string) BeautifyOptions
	ContentUnformatted([]string) BeautifyOptions
	ExtraLiners([]string) BeautifyOptions
	UnformattedContentDelimiter(string) BeautifyOptions
	IndentEmptyLines(bool) BeautifyOptions
	Templating([]BeautifyTemplating) BeautifyOptions
}

type beautifyOptions struct {
	data map[string]interface{}
}

func (o *beautifyOptions) IndentSize(indentSize uint) BeautifyOptions {
	ret := *o
	ret.data["indent_size"] = indentSize
	return &ret
}

func (o *beautifyOptions) IndentChar(character string) BeautifyOptions {
	ret := *o
	ret.data["indent_char"] = character
	return &ret
}

func (o *beautifyOptions) IndentWithTabs(b bool) BeautifyOptions {
	ret := *o
	ret.data["indent_with_tabs"] = b
	return &ret
}

func (o *beautifyOptions) Eol(string string) BeautifyOptions {
	ret := *o
	ret.data["eol"] = string
	return &ret
}

func (o *beautifyOptions) EndWithNewline(b bool) BeautifyOptions {
	ret := *o
	ret.data["end_with_newline"] = b
	return &ret
}

func (o *beautifyOptions) PreserveNewlines(b bool) BeautifyOptions {
	ret := *o
	ret.data["preserve_newlines"] = b
	return &ret
}

func (o *beautifyOptions) MaxPreserveNewlines(max uint) BeautifyOptions {
	ret := *o
	ret.data["max_preserve_newlines"] = max
	return &ret
}

func (o *beautifyOptions) IndentInnerHtml(b bool) BeautifyOptions {
	ret := *o
	ret.data["indent_inner_html"] = b
	return &ret
}

func (o *beautifyOptions) BraceStyle(braceStyle BeautifyBraceStyle) BeautifyOptions {
	ret := *o
	ret.data["brace_style"] = braceStyle
	return &ret
}

func (o *beautifyOptions) IndentScripts(indentScripts BeautifyIndentScripts) BeautifyOptions {
	ret := *o
	ret.data["indent_scripts"] = indentScripts
	return &ret
}

func (o *beautifyOptions) WrapLineLength(lineLength uint) BeautifyOptions {
	ret := *o
	ret.data["wrap_line_length"] = lineLength
	return &ret
}

func (o *beautifyOptions) WrapAttributes(wrapAttributes BeautifyWrapAttributes) BeautifyOptions {
	ret := *o
	ret.data["wrap_attributes"] = wrapAttributes
	return &ret
}

func (o *beautifyOptions) WrapAttributesIndentSize(indentSize uint) BeautifyOptions {
	ret := *o
	ret.data["wrap_attributes_indent_size"] = indentSize
	return &ret
}

func (o *beautifyOptions) Inline(tags []string) BeautifyOptions {
	ret := *o
	ret.data["inline"] = tags
	return &ret
}

func (o *beautifyOptions) Unformatted(tags []string) BeautifyOptions {
	ret := *o
	ret.data["unformatted"] = tags
	return &ret
}

func (o *beautifyOptions) ContentUnformatted(tags []string) BeautifyOptions {
	ret := *o
	ret.data["content_unformatted"] = tags
	return &ret
}

func (o *beautifyOptions) ExtraLiners(tags []string) BeautifyOptions {
	ret := *o
	ret.data["extra_liners"] = tags
	return &ret
}

func (o *beautifyOptions) UnformattedContentDelimiter(string string) BeautifyOptions {
	ret := *o
	ret.data["unformatted_content_delimiter"] = string
	return &ret
}

func (o *beautifyOptions) IndentEmptyLines(b bool) BeautifyOptions {
	ret := *o
	ret.data["indent_empty_lines"] = b
	return &ret
}

func (o *beautifyOptions) Templating(templating []BeautifyTemplating) BeautifyOptions {
	ret := *o
	ret.data["templating"] = templating
	return &ret
}

func NewBeautifyOptions() BeautifyOptions {
	return &beautifyOptions{
		data: map[string]interface{}{},
	}
}
