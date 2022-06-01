package mjml

// JuiceOptions is used to construct Juice options to be passed to the MJML compiler
// Detailed explanations of the options are here: https://github.com/Automattic/juice#options
type JuiceOptions interface {
	ApplyAttributesTableElements(bool) JuiceOptions
	ApplyHeightAttributes(bool) JuiceOptions
	ApplyStyleTags(bool) JuiceOptions
	ApplyWidthAttributes(bool) JuiceOptions
	ExtraCss(string) JuiceOptions
	InsertPreservedExtraCss(bool) JuiceOptions
	InlinePseudoElements(bool) JuiceOptions
	PreserveFontFaces(bool) JuiceOptions
	PreserveImportant(bool) JuiceOptions
	PreserveMediaQueries(bool) JuiceOptions
	PreserveKeyFrames(bool) JuiceOptions
	PreservePseudos(bool) JuiceOptions
	RemoveStyleTags(bool) JuiceOptions
	XmlMode(bool) JuiceOptions
}

type juiceOptions struct {
	data map[string]interface{}
}

func (o *juiceOptions) ApplyAttributesTableElements(b bool) JuiceOptions {
	ret := *o
	ret.data["applyAttributesTableElements"] = b
	return &ret
}

func (o *juiceOptions) ApplyHeightAttributes(b bool) JuiceOptions {
	ret := *o
	ret.data["applyHeightAttributes"] = b
	return &ret
}

func (o *juiceOptions) ApplyStyleTags(b bool) JuiceOptions {
	ret := *o
	ret.data["applyStyleTags"] = b
	return &ret
}

func (o *juiceOptions) ApplyWidthAttributes(b bool) JuiceOptions {
	ret := *o
	ret.data["applyWidthAttributes"] = b
	return &ret
}

func (o *juiceOptions) ExtraCss(s string) JuiceOptions {
	ret := *o
	ret.data["extraCss"] = s
	return &ret
}

func (o *juiceOptions) InsertPreservedExtraCss(b bool) JuiceOptions {
	ret := *o
	ret.data["insertPreservedExtraCss"] = b
	return &ret
}

func (o *juiceOptions) InlinePseudoElements(b bool) JuiceOptions {
	ret := *o
	ret.data["inlinePseudoElements"] = b
	return &ret
}

func (o *juiceOptions) PreserveFontFaces(b bool) JuiceOptions {
	ret := *o
	ret.data["preserveFontFaces"] = b
	return &ret
}

func (o *juiceOptions) PreserveImportant(b bool) JuiceOptions {
	ret := *o
	ret.data["preserveImportant"] = b
	return &ret
}

func (o *juiceOptions) PreserveMediaQueries(b bool) JuiceOptions {
	ret := *o
	ret.data["preserveMediaQueries"] = b
	return &ret
}

func (o *juiceOptions) PreserveKeyFrames(b bool) JuiceOptions {
	ret := *o
	ret.data["preserveKeyFrames"] = b
	return &ret
}

func (o *juiceOptions) PreservePseudos(b bool) JuiceOptions {
	ret := *o
	ret.data["preservePseudos"] = b
	return &ret
}

func (o *juiceOptions) RemoveStyleTags(b bool) JuiceOptions {
	ret := *o
	ret.data["removeStyleTags"] = b
	return &ret
}

func (o *juiceOptions) XmlMode(b bool) JuiceOptions {
	ret := *o
	ret.data["xmlMode"] = b
	return &ret
}

func NewJuiceOptions() JuiceOptions {
	return &juiceOptions{
		data: map[string]interface{}{},
	}
}
