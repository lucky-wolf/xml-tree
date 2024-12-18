package xmltree

type ParentChildElement struct {
	Parent *XMLElement
	Child  *XMLElement
}

type Finder func(element *XMLElement) bool

// find the given element which has that value
// this algo only finds the first one (if there are multiple)
// returns the parent and the element (parent is nil if th root element is the matching target)
func (tree *XMLTree) Find(tag, value string) (parent, element *XMLElement) {
	litmus := func(element *XMLElement) bool { return element.Matches(tag, value) }
	return tree.FindUsing(litmus)
}

// find the first parent which has all of these key-value pairs
func (tree *XMLTree) FindElementWithAll(properties ...SearchProperty) (element *XMLElement) {
	litmus := func(element *XMLElement) bool { return element.MatchesAll(properties...) }
	element, _ = tree.FindUsing(litmus)
	return
}

// finds the first element which your finder function responds true to
// this algo only finds the first one (if there are multiple)
// returns the parent and the element (parent is nil if th root element is the matching target)
func (tree *XMLTree) FindUsing(finder Finder) (parent, element *XMLElement) {

	// use a simple bfs algo
	var queue []*ParentChildElement
	pop := func() *ParentChildElement {
		if len(queue) == 0 {
			return nil
		}
		e := queue[0]
		queue = queue[1:]
		return e
	}

	// start with our (root) element(s) (which have no parent)
	for _, e := range tree.Elements.Elements() {
		queue = append(queue, &ParentChildElement{nil, e})
	}

	// process each node until we find it, or there are no more nodes
	for pc := pop(); pc != nil; pc = pop() {

		// check if we've found our target
		if finder(pc.Child) {
			parent = pc.Parent
			element = pc.Child
			return
		}

		// queue up this child's children
		for _, e := range pc.Child.Elements() {
			queue = append(queue, &ParentChildElement{pc.Child, e})
		}
	}

	return
}
