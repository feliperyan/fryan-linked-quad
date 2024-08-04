package simplequad

type DoublyLinkedListElement struct {
	Element
	Previous int32
	Next     int32
}

type ElementManager struct {
	Elements []DoublyLinkedListElement
}

func (em *ElementManager) Add(element Element, firstIndex int32) (int32, error) {

	em.Elements = append(em.Elements, DoublyLinkedListElement{element, -1, firstIndex})

	if firstIndex >= 0 {
		em.Elements[firstIndex].Previous = int32(len(em.Elements) - 1)
	}

	return int32(len(em.Elements) - 1), nil
}

// Remove returns a new head for this QuadrantNode if the Element removed was the Head
// returns the position of an Element that needs to be corrected in another QuadrantNode
// if the last item in the Slice happened to be the head for some QuadrantNode
func (em *ElementManager) Remove(element Element, firstIndex int32) (int32, int32, error) {

	index := firstIndex

	newHead := int32(-1)
	fixHead := int32(-1)

	for index != -1 {

		if element.id == em.Elements[index].id {
			thisPrev := em.Elements[index].Previous
			thisNext := em.Elements[index].Next

			// Link previous with Next therefore removing it from Linked List
			if thisPrev >= 0 {
				em.Elements[thisPrev].Next = thisNext
			}
			if thisNext >= 0 {
				em.Elements[thisNext].Previous = thisPrev

				// Just removed the head Element so return the new index
				if thisPrev == -1 {
					newHead = thisNext
				}
			}

			lastItem := em.Elements[len(em.Elements)-1]
			if lastItem.Previous == -1 { // reduce
				em.Elements[index] = lastItem

				if lastItem.Next >= 0 {
					em.Elements[lastItem.Next].Previous = index
				}

				fixHead = index
			}

			// Reduce length of slice
			em.Elements = em.Elements[:len(em.Elements)-1]

			break
		}

		index = em.Elements[index].Next
	}

	return newHead, fixHead, nil
}
