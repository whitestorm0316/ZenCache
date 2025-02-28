package cache

import (
	"testing"
)

// MockGetter is a simple mock implementation of the Getter interface.
type MockGetter struct{}

func (m MockGetter) Get(key string) ([]byte, error) {
	return []byte("mock value"), nil
}

func TestEngine_AddGroup(t *testing.T) {
	e := NewEngine()
	name := "testGroup"
	getter := MockGetter{}
	maxBytes := int64(100)

	// Add the first group
	e.AddGroup(name, getter, maxBytes)
	g1 := e.GetGroup(name)
	if g1 == nil {
		t.Error("expected group to be added, but it was nil")
	}

	// Add another group with the same name
	newGetter := MockGetter{}
	e.AddGroup(name, newGetter, maxBytes)
	g2 := e.GetGroup(name)
	if g2 == g1 {
		t.Error("expected new group to replace the old one, but they are the same")
	}
}

func TestEngine_GetGroup(t *testing.T) {
	e := NewEngine()
	name := "testGroup"
	getter := MockGetter{}
	maxBytes := int64(100)

	// Retrieve non-existent group
	g := e.GetGroup(name)
	if g != nil {
		t.Error("expected nil for non-existent group, but got a group")
	}

	// Add a group and retrieve it
	e.AddGroup(name, getter, maxBytes)
	g = e.GetGroup(name)
	if g == nil {
		t.Error("expected group to be retrieved, but got nil")
	}
}

// Additional tests for concurrency can be added following a similar pattern
