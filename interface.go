package simplequad

type IQuadtree interface {
	Insert(b Element) error
	Remove(b Element) error
	Get(id uint32) (Element, error)
	Search(X, Y, Width, Height float32) ([]Element, error)
}
