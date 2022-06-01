package mjml

import (
	"reflect"
	"testing"
)

func TestBeautifyOptions(t *testing.T) {

	options := NewBeautifyOptions().
		IndentSize(2).
		IndentChar(" ").
		IndentWithTabs(true).
		Eol("\n").
		EndWithNewline(true).
		PreserveNewlines(true).
		MaxPreserveNewlines(10).
		IndentInnerHtml(true).
		BraceStyle(BeautifyBraceStyleCollapse).
		IndentScripts(BeautifyIndentScriptsSeparate).
		WrapLineLength(1).
		WrapAttributes(BeautifyWrapAttributesAuto).
		WrapAttributesIndentSize(10).
		Inline([]string{"a", "span"}).
		Unformatted([]string{"a"}).
		ContentUnformatted([]string{"pre"}).
		ExtraLiners([]string{"head", "body"}).
		UnformattedContentDelimiter(" ").
		IndentEmptyLines(true).
		Templating([]BeautifyTemplating{BeautifyTemplatingAuto})

	expected := map[string]interface{}{
		"indent_size":                   uint(2),
		"indent_char":                   " ",
		"indent_with_tabs":              true,
		"eol":                           "\n",
		"end_with_newline":              true,
		"preserve_newlines":             true,
		"max_preserve_newlines":         uint(10),
		"indent_inner_html":             true,
		"brace_style":                   BeautifyBraceStyleCollapse,
		"indent_scripts":                BeautifyIndentScriptsSeparate,
		"wrap_line_length":              uint(1),
		"wrap_attributes":               BeautifyWrapAttributesAuto,
		"wrap_attributes_indent_size":   uint(10),
		"inline":                        []string{"a", "span"},
		"unformatted":                   []string{"a"},
		"content_unformatted":           []string{"pre"},
		"extra_liners":                  []string{"head", "body"},
		"unformatted_content_delimiter": " ",
		"indent_empty_lines":            true,
		"templating":                    []BeautifyTemplating{BeautifyTemplatingAuto},
	}

	beautifierOptions, ok := options.(*beautifyOptions)

	if !ok {
		t.Fatal("Options is not a *beautifierOptions")
	}

	if !reflect.DeepEqual(beautifierOptions.data, expected) {
		t.Error("BeautifierOptions does not match expected data")
	}
}
