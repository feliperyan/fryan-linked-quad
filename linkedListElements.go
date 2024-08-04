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

	// If there's another element already, make sure its previous points to this element being inserted now
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

	currentNodeHead := int32(-1)
	someOtherNodesHead := int32(-1)

	for index != -1 {

		if element.id == em.Elements[index].id {
			thisPrev := em.Elements[index].Previous
			thisNext := em.Elements[index].Next

			// Remove element from the linked list be re-assigning the
			// next and previous elements' pointers so they "skip" the removed
			// element
			if thisPrev >= 0 {
				em.Elements[thisPrev].Next = thisNext
			}
			if thisNext >= 0 {
				em.Elements[thisNext].Previous = thisPrev

				// Just removed the head Element so return the new index
				if thisPrev == -1 {
					currentNodeHead = thisNext
				}
			}

			// Slice clean up, overwrite the deleted element with the last element in the slice
			// then reduce slice length by 1. This is a O(1) operation for deleting items from a slice

			// Get last element
			lastElement := em.Elements[len(em.Elements)-1]

			// Overwrite removed Element
			em.Elements[index] = lastElement

			// Whoever was pointing to the index
			// of the lastElement needs to be corrected.
			if lastElement.Previous >= 0 {
				em.Elements[lastElement.Previous].Next = index

			} else { // Previous = -1 meaning this was the Element referenced by a Node
				someOtherNodesHead = index
			}

			if lastElement.Next >= 0 {
				em.Elements[lastElement.Next].Previous = index
			}

			// Reduce length of slice
			em.Elements = em.Elements[:len(em.Elements)-1]

			// Special case when the element being remove is the only one
			// in the slice
			if someOtherNodesHead == 0 && len(em.Elements) == 0 {
				someOtherNodesHead = -1
			}

			break
		}

		index = em.Elements[index].Next
	}

	return currentNodeHead, someOtherNodesHead, nil
}
