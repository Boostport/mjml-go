package mjml

import (
	"reflect"
	"testing"
)

func TestHTMLMiniferOptions(t *testing.T) {

	options := NewHTMLMinifierOptions().
		CaseSensitive(true).
		CollapseBooleanAttributes(true).
		CollapseInlineTagWhitespace(true).
		CollapseWhitespace(true).
		ConservativeCollapse(true).
		ContinueOnParseError(true).
		CustomAttrAssign([]string{"<div flex?=\"{{mode != cover}}\"></div>"}).
		CustomAttrCollapse("ng-class").
		CustomAttrSurround([]string{"<input {{#if value}}checked=\"checked\"{{/if}}>"}).
		DecodeEntities(true).
		HTML5(true).
		IgnoreCustomComments([]string{"^!"}).
		IgnoreCustomFragments([]string{"<%[\\s\\S]*?%>"}).
		IncludeAutoGeneratedTags(true).
		KeepClosingSlash(true).
		MaxLineLength(50).
		MinifyCSS(true).
		MinifyURLs(true).
		PreserveLineBreaks(true).
		PreventAttributesEscaping(true).
		ProcessConditionalComments(true).
		ProcessScripts([]string{"text/ng-template"}).
		QuoteCharacter(HTMLMinifierDoubleQuote).
		RemoveAttributeQuotes(true).
		RemoveComments(true).
		RemoveEmptyAttributes(true).
		RemoveEmptyElements(true).
		RemoveOptionalTags(true).
		RemoveRedundantAttributes(true).
		RemoveScriptTypeAttributes(true).
		RemoveStyleLinkTypeAttributes(true).
		RemoveTagWhitespace(true).
		SortAttributes(true).
		SortClassName(true).
		TrimCustomFragments(true).
		UseShortDoctype(true)

	expected := map[string]interface{}{
		"caseSensitive":                 true,
		"collapseBooleanAttributes":     true,
		"collapseInlineTagWhitespace":   true,
		"collapseWhitespace":            true,
		"conservativeCollapse":          true,
		"continueOnParseError":          true,
		"customAttrAssign":              []string{"<div flex?=\"{{mode != cover}}\"></div>"},
		"customAttrCollapse":            "ng-class",
		"customAttrSurround":            []string{"<input {{#if value}}checked=\"checked\"{{/if}}>"},
		"decodeEntities":                true,
		"html5":                         true,
		"ignoreCustomComments":          []string{"^!"},
		"ignoreCustomFragments":         []string{"<%[\\s\\S]*?%>"},
		"includeAutoGeneratedTags":      true,
		"keepClosingSlash":              true,
		"maxLineLength":                 uint(50),
		"minifyCSS":                     true,
		"minifyURLs":                    true,
		"preserveLineBreaks":            true,
		"preventAttributesEscaping":     true,
		"processConditionalComments":    true,
		"processScripts":                []string{"text/ng-template"},
		"quoteCharacter":                HTMLMinifierDoubleQuote,
		"removeAttributeQuotes":         true,
		"removeComments":                true,
		"removeEmptyAttributes":         true,
		"removeEmptyElements":           true,
		"removeOptionalTags":            true,
		"removeRedundantAttributes":     true,
		"removeScriptTypeAttributes":    true,
		"removeStyleLinkTypeAttributes": true,
		"removeTagWhitespace":           true,
		"sortAttributes":                true,
		"sortClassName":                 true,
		"trimCustomFragments":           true,
		"useShortDoctype":               true,
	}

	htmlMinifierOptions, ok := options.(*htmlMinifierOptions)

	if !ok {
		t.Fatal("Options is not a *htmlMinifierOptions")
	}

	if !reflect.DeepEqual(htmlMinifierOptions.data, expected) {
		t.Error("HTMLMinifierOptions does not match expected data")
	}
}
