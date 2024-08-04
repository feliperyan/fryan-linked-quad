package simplequad

import (
	"testing"
)

func TestNewElementsHolder(t *testing.T) {
	eh := NewElementsHolder(10)
	if len(eh.Pointers) != 0 {
		t.Errorf("Expected empty Pointers slice, got %v", len(eh.Pointers))
	}
}

func TestAdd(t *testing.T) {
	elements := &ElementSlice{
		Elements:   make([]Element, 0),
		ElementIds: make(map[int32]int32),
	}
	eh := NewElementsHolder(10)
	e1 := Element{id: 1, Box: Box{}}
	e2 := Element{id: 2, Box: Box{}}
	e2_same := Element{id: 2, Box: Box{}}

	pos, err := eh.Add(e1, -1, elements)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	pos2, err := eh.Add(e2, -1, elements)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if pos != 0 || pos2 != 1 {
		t.Errorf("Positions were not 0 and 1")
	}

	if len(elements.Elements) != 2 {
		t.Errorf("Expected 2 elements to have been inserted but got %v", len(elements.Elements))
	}

	pos3, err := eh.Add(e2_same, -1, elements)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if eh.Pointers[pos2].ElementIndex != eh.Pointers[pos3].ElementIndex {
		t.Errorf("Element inserted twice doesn't point to the same index in Elements array")
	}

}

//func TestRemove(t *testing.T) {
//	elements := &ElementSlice{
//		Elements:   make([]Element, 0),
//		ElementIds: make(map[int32]int32),
//	}
//
//	eh := NewElementsHolder(10)
//	e1 := Element{id: 55, Box: Box{}}
//	e2 := Element{id: 99, Box: Box{}}
//
//	_, err := eh.Add(e1, -1, elements)
//	if err != nil {
//		t.Errorf("Got error: %v", err)
//	}
//	_, err = eh.Add(e1, -1, elements)
//	if err != nil {
//		t.Errorf("Got error: %v", err)
//	}
//
//}
