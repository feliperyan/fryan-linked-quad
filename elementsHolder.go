package simplequad

import "fmt"

type Element struct {
	id int32
	Box
}

type ElementSlice struct {
	Elements   []Element
	ElementIds map[int32]int32 // elementID -> Elements Array Index
}

type ElementPointer struct {
	ElementIndex int32
	NextHolder   int32
}

type ElementsHolder struct {
	Pointers    []ElementPointer
	FreeIndexes *Stack[int32]
}

func NewElementsHolder(initialCap int) *ElementsHolder {
	return &ElementsHolder{
		Pointers:    make([]ElementPointer, 0, initialCap),
		FreeIndexes: NewStack[int32](),
	}
}

func (eh *ElementsHolder) Add(e Element, nextHolder int32, elements *ElementSlice) (int32, error) {

	ep := ElementPointer{
		ElementIndex: -1,
		NextHolder:   nextHolder,
	}

	// Has this element already been inserted?
	whichIndex, isElementAlreadyInserted := elements.ElementIds[e.id]

	// Don't re-insert the element if it's already there
	// this cover the cases where elements span across node boundaries
	// in other words, 1 element may have N pointers
	if isElementAlreadyInserted {
		ep.ElementIndex = whichIndex

	} else {
		elements.Elements = append(elements.Elements, e)
		whichIndex = int32(len(elements.Elements) - 1)
		ep.ElementIndex = whichIndex
		elements.ElementIds[e.id] = whichIndex
	}

	// Re-use a previously freed spot
	if !eh.FreeIndexes.IsEmpty() {
		i, err := eh.FreeIndexes.Pop()
		if err != nil {
			return -1, fmt.Errorf("Elements Holder error: %v", err)
		}

		eh.Pointers[i] = ep
		return i, nil
	}

	// No free spots so append to the end
	eh.Pointers = append(eh.Pointers, ep)
	return int32(len(eh.Pointers) - 1), nil
}

func (eh *ElementsHolder) Remove(e Element, firstPointer int32, elements *ElementSlice) (int32, error) {

	if firstPointer < 0 || int(firstPointer) >= len(eh.Pointers) {
		return -1, fmt.Errorf("firstPointer is out of bounds")
	}

	indexInLinkedList := firstPointer
	for indexInLinkedList != -1 {
		element := elements.Elements[eh.Pointers[indexInLinkedList].ElementIndex]

		if element.id == e.id {
			// remove from
			delete(elements.ElementIds, e.id)
			eh.FreeIndexes.Push(indexInLinkedList)

		}

		indexInLinkedList = eh.Pointers[indexInLinkedList].NextHolder
	}

	return -1, nil
}
