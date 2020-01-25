package md

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tix/ticket/body"
)

func TestListItemSegmentParser_Parse(t *testing.T) {
	text := `
- Item 1
`
	parser := NewListItemSegmentParser(false, 1, "-", 0)
	state, rootNode := setupTextParser(text)
	state.StartTicket()
	node := rootNode.FirstChild()

	err := parser.Parse(state, node.FirstChild())

	expectedBody := []body.Segment{
		body.NewBulletListItemSegment(1, "-"),
		body.NewLineBreakSegment(),
	}
	ticketBody := state.CurrentTicket().Body
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, ticketBody)
}
