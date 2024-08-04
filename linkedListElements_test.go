package simplequad

import "testing"

func TestAddingElements(t *testing.T) {
	em := ElementManager{Elements: make([]DoublyLinkedListElement, 0)}

	e1 := Element{
		id:  1,
		Box: Box{},
	}
	e2 := Element{
		id:  2,
		Box: Box{},
	}

	pos, err := em.Add(e1, -1)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if pos != 0 {
		t.Errorf("Expected pos == 0, got: %v", pos)
	}

	_, err = em.Add(e2, pos)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	if len(em.Elements) != 2 {
		t.Errorf("Expected 2 Elements, got %v", len(em.Elements))
	}
}

func TestRemovingElement(t *testing.T) {
	em := ElementManager{Elements: make([]DoublyLinkedListElement, 0)}

	e1 := Element{
		id:  1,
		Box: Box{},
	}
	e2 := Element{
		id:  2,
		Box: Box{},
	}

	pos, err := em.Add(e1, -1)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	head, err := em.Add(e2, pos)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	newHead, fixHead, err := em.Remove(e2, head)

	if newHead != 0 {
		t.Errorf("expected newHead == 0, got %v", newHead)
	}
	if fixHead != 1 {
		t.Errorf("expected fixHead == 1, got %v", fixHead)
	}

	em = ElementManager{Elements: make([]DoublyLinkedListElement, 0)}
	pos, err = em.Add(e1, -1)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	head, err = em.Add(e2, pos)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	e3 := Element{
		id:  3,
		Box: Box{},
	}
	_, err = em.Add(e3, -1)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}

	newHead, fixHead, err = em.Remove(e2, head)
	if fixHead != 1 {
		t.Errorf("expected fixHead == 1, got %v", fixHead)
	}

	if len(em.Elements) != 2 {
		t.Errorf("Expected length to be 2, got %v", len(em.Elements))
	}
}
