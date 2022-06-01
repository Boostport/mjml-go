package mjml

import (
	"fmt"
	"strings"
)

type Error struct {
	Message string `json:"message"`
	Details []struct {
		Line    int    `json:"line"`
		Message string `json:"message"`
		TagName string `json:"tagName"`
	} `json:"details"`
}

func (e Error) Error() string {

	var sb strings.Builder

	sb.WriteString(e.Message)

	numDetails := len(e.Details)

	if numDetails > 0 {
		sb.WriteString(":\n")
	}

	for i, detail := range e.Details {
		sb.WriteString(fmt.Sprintf("- Line %d of (%s) - %s", detail.Line, detail.TagName, detail.Message))

		if i != numDetails-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
