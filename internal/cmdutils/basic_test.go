package cmdutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockNodeInfo is a mock implementation of the NodeInfo interface
type mockNodeInfo struct {
	name string
}

func (m mockNodeInfo) Code() string   { return "" }
func (m mockNodeInfo) PrintNode()     {}
func (m mockNodeInfo) PrintComments() {}
func (m mockNodeInfo) Name() string   { return m.name }

func TestFindLengthOfOutput(t *testing.T) {
	output := "short\nmuch longer line\nmid"
	expected := len("much longer line")
	assert.Equal(t, findLengthOfOutput(output), expected)
}

func TestGetMax(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{3, 5, 5},
		{10, 2, 10},
		{7, 7, 7},
	}

	for _, tt := range tests {
		assert.Equal(t, getMax(tt.a, tt.b), tt.expected)
	}
}

func TestCapitalizeString(t *testing.T) {
	input := "some-long-name"
	expected := "Some Long Name"
	assert.Equal(t, capitalizeString(input), expected)
}

func TestGetNameFromNodes(t *testing.T) {
	node := mockNodeInfo{name: "TestNode"}

	name := getNameFromNodes(node)
	assert.Equal(t, "TestNode", name)

	ptr := &node
	name = getNameFromNodes(ptr)
	assert.Equal(t, "TestNode", name)

	name = getNameFromNodes(123)
	assert.Equal(t, "", name)
}
