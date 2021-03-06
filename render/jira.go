package render

import (
	"fmt"
	"strings"
	"tix/ticket/body"
)

type JiraBodyRenderer struct{}

func NewJiraBodyRenderer() *JiraBodyRenderer {
	return &JiraBodyRenderer{}
}

func (j JiraBodyRenderer) RenderSegment(bodySegment body.Segment) string {
	switch segment := bodySegment.(type) {
	case *body.BlockQuoteSegment:
		return j.renderBlockQuoteItem()
	case *body.BulletListItemSegment:
		return j.renderBulletListItem(segment)
	case *body.CodeBlockSegment:
		return j.renderCodeBlock(segment)
	case *body.CodeSpanSegment:
		return j.renderCodeSpan(segment)
	case *body.EmphasisSegment:
		return j.renderEmphasis(segment)
	case *body.LinkSegment:
		return j.renderLink(segment)
	case *body.ListEndSegment:
		return j.renderListEnd(segment)
	case *body.ListStartSegment:
		return segment.Value()
	case *body.LineBreakSegment:
		return segment.Value()
	case *body.OrderedListItemSegment:
		return j.renderOrderedListItem(segment)
	case *body.StrongEmphasisSegment:
		return j.renderStrongEmphasis(segment)
	case *body.TextBlockSegment:
		return segment.Value()
	case *body.TextSegment:
		return segment.Value()
	case *body.ThematicBreakSegment:
		return j.renderThematicBreak()
	default:
		return segment.Value()
	}
}

func (j JiraBodyRenderer) renderBlockQuoteItem() string {
	return "{quote}"
}

func (j JiraBodyRenderer) renderBulletListItem(segment *body.BulletListItemSegment) string {
	var builder strings.Builder
	level := segment.Attributes().Level

	for ii := 0; ii < level; ii++ {
		builder.WriteString(segment.Value())
	}
	builder.WriteString(" ")

	return builder.String()
}

func (j JiraBodyRenderer) renderCodeBlock(segment *body.CodeBlockSegment) string {
	var builder strings.Builder
	lang := segment.Attributes().Language

	if len(lang) > 0 {
		marker := fmt.Sprintf("{code:%s}\n", lang)
		builder.WriteString(marker)
	} else {
		builder.WriteString("{code}\n")
	}
	builder.WriteString(segment.Value())
	builder.WriteString("{code}")

	return builder.String()
}

func (j JiraBodyRenderer) renderCodeSpan(segment *body.CodeSpanSegment) string {
	suffix := j.codeSpanSuffix(segment)
	return fmt.Sprintf("{{%s}}%s", segment.Value(), suffix)
}

// Jira has a hard time formatting spans which are followed immediately with a non-empty string (ex - {{span}}text).
// This function will add a space suffix in these cases to prevent super funky Jira formatting.
func (j JiraBodyRenderer) codeSpanSuffix(segment *body.CodeSpanSegment) string {
	textSegment, ok := segment.Next().(*body.TextSegment)
	if ok && !strings.HasPrefix(textSegment.Value(), " ") {
		return " "
	}
	return ""
}

func (j JiraBodyRenderer) renderEmphasis(segment *body.EmphasisSegment) string {
	return fmt.Sprintf("_%s_", segment.Value())
}

func (j JiraBodyRenderer) renderLink(segment *body.LinkSegment) string {
	url := segment.Attributes().Url
	return fmt.Sprintf("[%s|%s]", segment.Value(), url)
}

func (j JiraBodyRenderer) renderOrderedListItem(segment *body.OrderedListItemSegment) string {
	var builder strings.Builder
	level := segment.Attributes().Level

	for ii := 0; ii < level; ii++ {
		builder.WriteString("#")
	}
	builder.WriteString(" ")

	return builder.String()
}

func (j JiraBodyRenderer) renderListEnd(segment *body.ListEndSegment) string {
	if segment.Attributes().Level == 1 {
		return "\n"
	} else {
		return ""
	}
}

func (j JiraBodyRenderer) renderStrongEmphasis(segment *body.StrongEmphasisSegment) string {
	return fmt.Sprintf("*%s*", segment.Value())
}

func (j JiraBodyRenderer) renderThematicBreak() string {
	return "----"
}
