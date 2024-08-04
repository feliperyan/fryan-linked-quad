package simplequad

import "testing"

func TestCollides(t *testing.T) {
	type args struct {
		box1 Box
		box2 Box
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "box2 right of box 1",
			args: args{
				box1: Box{10, 10, 10, 10},
				box2: Box{21, 10, 10, 10},
			},
			want: false,
		},
		{
			name: "box2 left of box 1",
			args: args{
				box1: Box{21, 10, 10, 10},
				box2: Box{10, 10, 10, 10},
			},
			want: false,
		},
		{
			name: "box2 above box 1",
			args: args{
				box1: Box{10, 21, 10, 10},
				box2: Box{10, 10, 10, 10},
			},
			want: false,
		},
		{
			name: "box2 below box 1",
			args: args{
				box1: Box{10, 10, 10, 10},
				box2: Box{10, 21, 10, 10},
			},
			want: false,
		},
		{
			name: "box2 inside box 1",
			args: args{
				box1: Box{10, 10, 10, 10},
				box2: Box{15, 15, 3, 3},
			},
			want: true,
		},
		{
			name: "touches on the right, no overlap",
			args: args{
				box1: Box{10, 10, 10, 10},
				box2: Box{20, 10, 10, 10},
			},
			want: true,
		},
		{
			name: "box1 overlaps smaller box2",
			args: args{
				box1: Box{10, 10, 10, 10},
				box2: Box{15, 15, 10, 10},
			},
			want: true,
		},
		{
			name: "box1 above box2 no overlap since y grows downwards",
			args: args{
				box1: Box{10, 10, 10, 10},
				box2: Box{10, 21, 10, 10},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Collides(tt.args.box1, tt.args.box2); got != tt.want {
				t.Errorf("Collides() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TL TR BL* BR -> BLTL BLTR BLBL BLBR
func TreeOne() *QuadRoot {
	tree := NewQuadRoot(16, 16, 3, 3, 0)

	tree.Nodes[0].NextQuadrant = 1
	tree.Nodes[0].NumberElements = -1

	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{0, 0, 8, 8}, 1, 0,
		0, -1, -1))
	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{8, 0, 8, 8}, 1, 0,
		0, -1, -1))
	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{0, 8, 8, 8}, 1, 0,
		-1, 5, -1)) // branch
	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{8, 8, 8, 8}, 1, 0,
		0, -1, -1))

	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{0, 8, 4, 4}, 2, 3,
		0, -1, -1))
	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{4, 8, 4, 4}, 2, 3,
		0, -1, -1))
	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{0, 12, 4, 4}, 2, 3,
		0, -1, -1))
	tree.Nodes = append(tree.Nodes, *NewQuadNode(Box{4, 12, 4, 4}, 2, 3,
		0, -1, -1))

	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{8, 0, 8, 8},
	//	Depth:          1,
	//	Root:           0,
	//	NumberElements: 0,
	//	NextQuadrant:   -1,
	//})
	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{0, 8, 8, 8},
	//	Depth:          1,
	//	Root:           0,
	//	NumberElements: -1,
	//	NextQuadrant:   5,
	//}) // branch
	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{8, 8, 8, 8},
	//	Depth:          1,
	//	Root:           0,
	//	NumberElements: -1,
	//	NextQuadrant:   -1,
	//})
	//
	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{0, 8, 4, 4},
	//	Depth:          2,
	//	Root:           3,
	//	NumberElements: 0,
	//	NextQuadrant:   -1,
	//})
	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{4, 8, 4, 4},
	//	Depth:          2,
	//	Root:           3,
	//	NumberElements: 0,
	//	NextQuadrant:   -1,
	//})
	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{0, 12, 4, 4},
	//	Depth:          2,
	//	Root:           3,
	//	NumberElements: 0,
	//	NextQuadrant:   -1,
	//})
	//tree.Nodes = append(tree.Nodes, QuadNode{
	//	Quadrant:       Box{4, 12, 4, 4},
	//	Depth:          2,
	//	Root:           3,
	//	NumberElements: 0,
	//	NextQuadrant:   -1,
	//})

	return tree
}

func TestGetOverlappingQuadrants(t *testing.T) {

	tree := TreeOne()

	nodes := tree.getCollidingLeaves(Box{13, 13, 2, 2})

	got := len(nodes)
	wants := 1
	if wants != got {
		t.Errorf("getCollidingLeaves() = %v, wanted %v", got, wants)
	}

}

func TestInsertDepth2(t *testing.T) {
	tree := TreeOne()
	b := Element{
		id:  0,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}

	// bl -> bltr
	tree.Insert(b)

	wants := int32(1)
	got := tree.Nodes[6].NumberElements

	if wants != got {
		t.Errorf("Insert() = %v, wanted %v", got, wants)
	}
}

func TestInsertDepth2AndDivide(t *testing.T) {
	tree := TreeOne()
	b := Element{
		id:  0,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}

	tree.Nodes[6].NumberElements = 4

	// bl -> bltr
	tree.Insert(b)

	wants := 13
	got := len(tree.Nodes)

	if wants != got {
		t.Errorf("Insert() = %v, wanted %v", got, wants)
	}
}

// Can I insert the first element
func TestInsertFirstElement(t *testing.T) {

	tree := NewQuadRoot(16, 16, 3, 3, 0)
	b := Element{
		id:  0,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}

	err := tree.Insert(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	want := 1
	//got := len(tree.Elements.Elements)
	got := len(tree.ElementManager.Elements)

	if want != got {
		t.Errorf("Wanted %v , got %v", want, got)
	}

}

func TestInsertTwoElements(t *testing.T) {

	tree := NewQuadRoot(16, 16, 3, 3, 0)
	a := Element{
		id:  0,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}
	b := Element{
		id:  0,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}
	//
	//err := tree.Insert(a)
	//if err != nil {
	//	t.Errorf("error: %v", err)
	//}
	//err = tree.Insert(b)
	//if err != nil {
	//	t.Errorf("error: %v", err)
	//}
	//
	//want := 1
	////got := len(tree.Elements.Elements)
	//got := len(tree.ElementManager.Elements)
	//if want != got {
	//	t.Errorf("Wanted %v , got %v", want, got)
	//}

	// Same tree but now diff element IDs
	tree = NewQuadRoot(16, 16, 3, 3, 0)
	b.id = 99
	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = tree.Insert(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	want := 2
	//got = len(tree.Elements.Elements)
	got := len(tree.ElementManager.Elements)
	if want != got {
		t.Errorf("Wanted %v , got %v", want, got)
	}

}

func TestInsertTwoElementsIntoChildNode(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()
	a := Element{
		id:  0,
		Box: Box{X: 1, Y: 1, Width: 1, Height: 1},
	}
	b := Element{
		id:  1,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Insert(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	want := 2
	//got := len(tree.Elements.Elements)
	got := len(tree.ElementManager.Elements)
	if want != got {
		t.Errorf("Wanted %v , got %v", want, got)
	}

}

func TestInsertMultipleElementsIntoChildNode(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()

	// TL
	a := Element{
		id:  0,
		Box: Box{X: 1, Y: 1, Width: 1, Height: 1},
	}
	d := Element{
		id:  3,
		Box: Box{X: 1, Y: 1, Width: 1, Height: 1},
	}

	// BL -> BLTR
	b := Element{
		id:  1,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}
	c := Element{
		id:  2,
		Box: Box{X: 5, Y: 9, Width: 1, Height: 1},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Insert(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Insert(c)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Insert(d)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	want := 4
	//got := len(tree.Elements.Elements)
	got := len(tree.ElementManager.Elements)
	if want != got {
		t.Errorf("Wanted %v , got %v", want, got)
	}

}

func TestInsertAcrossNodes(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()

	// TL and TR
	a := Element{
		id:  0,
		Box: Box{X: 7, Y: 1, Width: 3, Height: 1},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	want := 2
	//got := len(tree.ElementPointers.Pointers)
	got := len(tree.ElementManager.Elements)
	if want != got {
		t.Errorf("Wanted %v , got %v", want, got)
	}

	//isTrue := tree.ElementPointers.Pointers[0].ElementIndex == tree.ElementPointers.Pointers[1].ElementIndex
	isTrue := tree.ElementManager.Elements[0].id == tree.ElementManager.Elements[1].id
	if !isTrue {
		t.Errorf("Wanted %v , got %v", true, isTrue)
	}

}

func TestInsertAcrossSubNodes(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()

	// TL and TR
	a := Element{
		id:  0,
		Box: Box{X: 1, Y: 7, Width: 1, Height: 6},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	want := 3
	//got := len(tree.ElementPointers.Pointers)
	got := len(tree.ElementManager.Elements)
	if want != got {
		t.Errorf("Wanted %v , got %v", want, got)
	}

}

func TestRemoveFromSingleNode(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()

	// TL and TR
	a := Element{
		id:  0,
		Box: Box{X: 1, Y: 1, Width: 1, Height: 1},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Remove(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(tree.ElementManager.Elements) != 0 {
		t.Errorf("Expected lengh 0, got %v", len(tree.ElementManager.Elements))
	}
}

func TestRemoveFromSingleNode2(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()

	// TL and TR
	a := Element{
		id:  0,
		Box: Box{X: 1, Y: 1, Width: 1, Height: 1},
	}
	b := Element{
		id:  1,
		Box: Box{X: 9, Y: 1, Width: 1, Height: 1},
	}
	c := Element{
		id:  2,
		Box: Box{X: 9, Y: 1, Width: 1, Height: 1},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = tree.Insert(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = tree.Insert(c)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Remove(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if tree.Nodes[2].FirstElementHolder != 0 {
		t.Errorf("Exected FirstHolder for TR to be 0, got %v", tree.Nodes[2].FirstElementHolder)
	}
}

func TestRemoveFromMultipleNodes(t *testing.T) {

	// TL TR BL* BR -> BLTL BLTR BLBL BLBR
	tree := TreeOne()

	// TL and TR
	a := Element{
		id:  0,
		Box: Box{X: 1, Y: 1, Width: 1, Height: 1},
	}
	b := Element{
		id:  1,
		Box: Box{X: 7, Y: 1, Width: 3, Height: 1},
	}
	c := Element{
		id:  2,
		Box: Box{X: 9, Y: 1, Width: 1, Height: 1},
	}

	err := tree.Insert(a)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = tree.Insert(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	err = tree.Insert(c)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	err = tree.Remove(b)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if tree.Nodes[1].NumberElements != 1 || tree.Nodes[2].NumberElements != 1 {
		t.Errorf("Exected TL and TR NumberOfElements to be 1, got %v and %v",
			tree.Nodes[1].NumberElements, tree.Nodes[2].NumberElements)
	}
}

// Can the tree subdivide the root quadrant once N > MaxElements
func TestInsertElementCauseDivisionDepth0(t *testing.T) {}

// Can handle elements that span quadrants
func TestInsertOverlapsQuadrants(t *testing.T) {}

// Can subdivide a quadrant once Pointers >= MaxElements
func TestInsertElementCauseDivisionDepth1(t *testing.T) {}

// Can stop subdividing if Depth >= MaxDepth
func TestInsertElementCauseDivisionDepthMax(t *testing.T) {}

// Can handle elements that span quadrants and subdivides quadrant
func TestInsertOverlapsQuadrantsCausesSubdivision(t *testing.T) {}

// Test simple collapse
func TestCollapseQuadrant(t *testing.T) {
	tree := TreeOne()

	combineNodes(tree, 3)

	wants := 5
	got := len(tree.Nodes)

	if wants != got {
		t.Errorf("combineNodes() = %v, wanted %v", got, wants)
	}
}
