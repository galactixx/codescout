package codescout

import (
	"path/filepath"
	"testing"
)

func TestScoutFunction(t *testing.T) {
	path := filepath.Join("pkg", "codescout", "testdata", "scout_single.go")
	funcConfig := FunctionConfig{Name: "Greet"}
	funcNode, err := ScoutFunction(path, funcConfig)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	if funcNode.Node.Line != 20 {
		t.Errorf("got %v, want %v", funcNode.Node.Line, 20)
	}
}
