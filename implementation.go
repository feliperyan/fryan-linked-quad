package simplequad

import (
	"fmt"
)

// Box is a simple 2d rectangle. X,Y mark the top left corner of the box
type Box struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func Collides(box1 Box, box2 Box) bool {
	// Check for non-intersection on X axis
	if box1.X+box1.Width < box2.X ||
		box2.X+box2.Width < box1.X {
		return false
	}

	// Check for non-intersection on Y axis
	if box1.Y+box1.Height < box2.Y ||
		box2.Y+box2.Height < box1.Y {
		return false
	}

	// Boxes are colliding on both axes
	return true
}

type QuadRoot struct {
	MaxDepth    int8
	MaxElements int32
	Width       int32
	Height      int32
	CleanMode   int8
	Nodes       []QuadNode
	//ElementPointers *ElementsHolder
	//Elements        *ElementSlice
	ElementManager *ElementManager
}

func NewQuadRoot(width int32, height int32, maxDepth int8, maxElements int32, cleanMode int8) *QuadRoot {
	q := &QuadRoot{
		MaxDepth:       maxDepth,
		MaxElements:    maxElements,
		CleanMode:      cleanMode,
		Nodes:          make([]QuadNode, 0, 100),
		Width:          width,
		Height:         height,
		ElementManager: &ElementManager{make([]DoublyLinkedListElement, 0)},
		//ElementPointers: NewElementsHolder(int(maxElements)),
		//Elements: &ElementSlice{
		//	Elements:   make([]Element, 0, 1000),
		//	ElementIds: make(map[int32]int32),
		//},
	}

	q.Nodes = append(q.Nodes, QuadNode{
		Quadrant:           Box{X: 0, Y: 0, Width: float32(width), Height: float32(height)},
		Depth:              0,
		Root:               -1,
		NumberElements:     0,
		NextQuadrant:       -1,
		FirstElementHolder: -1,
	})

	return q
}

type QuadNode struct {
	Quadrant           Box
	Depth              uint8
	Root               int32
	NumberElements     int32 // if -1 this node is a branch and nextNode points to the 1st node a level down
	NextQuadrant       int32 // The 1st node a level down
	FirstElementHolder int32
}

func NewQuadNode(quadrant Box, d uint8, r int32, numels int32, nextquad int32, firstElhold int32) *QuadNode {
	return &QuadNode{
		Quadrant:           quadrant,
		Depth:              d,
		Root:               r,
		NumberElements:     numels,
		NextQuadrant:       nextquad,
		FirstElementHolder: firstElhold,
	}
}

func (qn *QuadNode) isLeaf() bool {
	return qn.NextQuadrant == -1
}

/*
------ Public Functions ------
*/

func (q *QuadRoot) Insert(e Element) error {

	// Find correct node to insert, may be more than one
	// if shape spans multiple nodes.
	nodes := q.getCollidingLeaves(e.Box)
	if len(nodes) < 1 {
		return fmt.Errorf("element's position and/or size are invalid")
	}

	for _, p := range nodes {
		nodePos := p

		if q.Nodes[nodePos].NumberElements < q.MaxElements {

			//index, err := q.ElementPointers.Add(e, q.Nodes[nodePos].FirstElementHolder, q.Elements)

			index, err := q.ElementManager.Add(e, q.Nodes[nodePos].FirstElementHolder)
			if err != nil {
				return fmt.Errorf("Failed to insert into node %v with error: %v", q.Nodes[nodePos], err)
			}

			// Newly inserted holder is the new head of the linked list
			q.Nodes[nodePos].FirstElementHolder = index
			q.Nodes[nodePos].NumberElements += 1

		} else {
			// NumberElements > MaxElements --> Must divide and then insert
			divideNode(q, nodePos)
		}

	}
	return nil
}

func (q *QuadRoot) Remove(e Element) error {
	nodes := q.getCollidingLeaves(e.Box)

	for _, nodeIndex := range nodes {
		newHead, mustFix, err := q.ElementManager.Remove(e, q.Nodes[nodeIndex].FirstElementHolder)
		if err != nil {
			return err
		}
		if newHead >= 0 {
			q.Nodes[nodeIndex].FirstElementHolder = newHead
		}
		q.Nodes[nodeIndex].NumberElements = q.Nodes[nodeIndex].NumberElements - 1

		if mustFix >= 0 {
			// Get all nodes for the element which used to be at the end of the Elements slice
			// but now is at "mustFix". One of the nodes must point to the old last element index. Correct it.
			for _, i := range q.getCollidingLeaves(q.ElementManager.Elements[mustFix].Box) {
				if q.Nodes[i].FirstElementHolder == int32(len(q.ElementManager.Elements)) {
					q.Nodes[i].FirstElementHolder = mustFix
					break
				}
			}
		}
	}

	return nil
}

func (q *QuadRoot) Get(id uint32) (Element, error) {
	//TODO implement me
	panic("implement me")
}

func (q *QuadRoot) Search(X, Y, Width, Height float32) ([]Element, error) {
	//TODO implement me
	panic("implement me")
}

/*
------ Private Functions ------
*/

func (q *QuadRoot) getCollidingLeaves(b Box) []int {

	branchPositions := make([]int32, 0, 16)
	branchPositions = append(branchPositions, 0)

	nodesFound := make([]int, 0, 16)

	for len(branchPositions) > 0 {
		// get last element and remove it
		node := &q.Nodes[branchPositions[len(branchPositions)-1]]
		pos := branchPositions[len(branchPositions)-1]
		branchPositions = branchPositions[:len(branchPositions)-1]

		// test for root node
		if node.Root == -1 {
			if node.NextQuadrant == -1 && Collides(node.Quadrant, b) {
				return append(nodesFound, int(pos))
			}

			if Collides(node.Quadrant, b) && node.NextQuadrant >= 0 {
				branchPositions = append(branchPositions, node.NextQuadrant)
				continue
			}
		}

		// Nodes always added as 4, contiguously
		for i := 0; i < 4; i++ {
			thisNode := q.Nodes[int(pos)+i]
			if !Collides(thisNode.Quadrant, b) {
				continue
			}
			if !thisNode.isLeaf() {
				branchPositions = append(branchPositions, thisNode.NextQuadrant)
				continue
			}
			nodesFound = append(nodesFound, int(pos)+i)
		}
	}

	return nodesFound
}

func divideNode(q *QuadRoot, pos int) {
	node := &q.Nodes[pos]

	q0 := QuadNode{
		Quadrant:       Box{X: node.Quadrant.X, Y: node.Quadrant.Y, Width: node.Quadrant.Width / 2, Height: node.Quadrant.Height / 2},
		Depth:          node.Depth + 1,
		Root:           int32(pos),
		NumberElements: 0,
		NextQuadrant:   -1,
	}
	q1 := QuadNode{
		Quadrant:       Box{X: node.Quadrant.X + node.Quadrant.Width/2, Y: node.Quadrant.Y, Width: node.Quadrant.Width / 2, Height: node.Quadrant.Height / 2},
		Depth:          node.Depth + 1,
		Root:           int32(pos),
		NumberElements: 0,
		NextQuadrant:   -1,
	}
	q2 := QuadNode{
		Quadrant:       Box{X: node.Quadrant.X, Y: node.Quadrant.Y + node.Quadrant.Height/2, Width: node.Quadrant.Width / 2, Height: node.Quadrant.Height / 2},
		Depth:          node.Depth + 1,
		Root:           int32(pos),
		NumberElements: 0,
		NextQuadrant:   -1,
	}
	q3 := QuadNode{
		Quadrant:       Box{X: node.Quadrant.X + node.Quadrant.Width/2, Y: node.Quadrant.Y + node.Quadrant.Height/2, Width: node.Quadrant.Width / 2, Height: node.Quadrant.Height / 2},
		Depth:          node.Depth + 1,
		Root:           int32(pos),
		NumberElements: 0,
		NextQuadrant:   -1,
	}

	q.Nodes = append(q.Nodes, q0)
	q.Nodes = append(q.Nodes, q1)
	q.Nodes = append(q.Nodes, q2)
	q.Nodes = append(q.Nodes, q3)

	node.NumberElements = -1
	node.NextQuadrant = int32(len(q.Nodes) - 3)

	//TODO: now transfer all the elements...

}

func combineNodes(q *QuadRoot, rootPos int) error {
	root := &q.Nodes[rootPos]
	if root.isLeaf() {
		return fmt.Errorf("node is not a branch")
	}

	firstQuadrantOfFour := root.NextQuadrant

	// TODO: EFFICIENCY = this is grossly inneficient. Either figure out
	// the best solution is to have a data structure that marks items as deleted
	// then overrides the deleted items. That way I dont shift items around.
	q.Nodes = append(q.Nodes[:firstQuadrantOfFour], q.Nodes[firstQuadrantOfFour+4:]...)

	// TODO: Move elements to root note
	return nil
}
