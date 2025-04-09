package cmdutils

import "testing"

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

	if result := findLengthOfOutput(output); result != expected {
		t.Errorf("Expected %d, got %d", expected, result)
	}
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
		if result := getMax(tt.a, tt.b); result != tt.expected {
			t.Errorf("GetMax(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestCapitalizeString(t *testing.T) {
	input := "some-long-name"
	expected := "Some Long Name"

	if result := capitalizeString(input); result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestGetNameFromNodes(t *testing.T) {
	node := mockNodeInfo{name: "TestNode"}

	name := getNameFromNodes(node)
	if name != "TestNode" {
		t.Errorf("Expected 'TestNode', got '%s'", name)
	}

	ptr := &node
	name = getNameFromNodes(ptr)
	if name != "TestNode" {
		t.Errorf("Expected 'TestNode' (pointer), got '%s'", name)
	}

	name = getNameFromNodes(123)
	if name != "" {
		t.Errorf("Expected empty string for non-matching type, got '%s'", name)
	}
}
