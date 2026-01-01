package xmltree

import (
	"fmt"
	"log"

	"github.com/lucky-wolf/xml-tree/etc"
)

// returns XMLElements only
// warn: this slice doesn't correspond exactly with our contents (e.g. comments are omitted)
func (v *XMLValue) Elements() (elements []*XMLElement) {
	// it is legal to call on a nil value (we simply have no child elements)
	if v == nil {
		return
	}

	// we are either a single or multiple elements
	switch v := v.contents.(type) {
	case *XMLElement:
		// single element
		elements = append(elements, v)
	case []any:
		// multiple child elements
		for _, e := range v {
			switch v := e.(type) {
			case *XMLElement:
				elements = append(elements, v)
			}
		}
	default:
		// we have no child elements
	}

	return
}

type ElementWithIndex struct {
	Element   *XMLElement
	TrueIndex int
}

// returns XMLElements paired with their true indexes in the contents array
// this is essential because Elements() filters out comments/directives,
// so its indexes don't correspond 1:1 with the actual contents array
// each pair is [element, trueIndex]
func (v *XMLValue) ElementsWithIndexes() (pairs []ElementWithIndex) {
	// it is legal to call on a nil value (we simply have no child elements)
	if v == nil {
		return
	}

	// we are either a single or multiple elements
	switch v := v.contents.(type) {
	case *XMLElement:
		// single element at index 0
		pairs = append(pairs, ElementWithIndex{Element: v, TrueIndex: 0})
	case []any:
		// iterate through contents and pair elements with their true index
		for i, e := range v {
			if elem, ok := e.(*XMLElement); ok {
				pairs = append(pairs, ElementWithIndex{Element: elem, TrueIndex: i})
			}
		}
	default:
		// we have no child elements
	}

	return
}

func (v XMLValue) Clone() XMLValue {
	// v is already a shallow copy, just do a deep copy on the contents
	v.contents = CloneContents(v.contents)
	return v
}

func (v *XMLValue) SetContents(contents any) {
	switch contents.(type) {
	case []any:
	case *XMLElement:
	case *XMLComment:
	case *XMLDirective:
	case *XMLProcInst:
	case string:
	default:
		err := fmt.Errorf("invalid content type: %T", contents)
		panic(err)
	}

	v.contents = contents
}

func (v *XMLValue) CloneContents() any {
	return CloneContents(v.contents)
}

func CloneContents(contents any) any {
	switch t := contents.(type) {

	case []any:
		// multiple child contents
		contents := []any{}
		for _, e := range t {
			contents = append(contents, CloneContents(e))
		}
		return contents

	case *XMLElement:
		return &XMLElement{StartElement: t.StartElement.Copy(), XMLValue: t.XMLValue.Clone()}
	case *XMLComment:
		return &XMLComment{Comment: t.Copy()}
	case *XMLDirective:
		return &XMLDirective{Directive: t.Copy()}
	case *XMLProcInst:
		return &XMLProcInst{ProcInst: t.Copy()}

	case string:
		return t
	}

	err := fmt.Errorf("cannot clone: invalid contents: %T", contents)
	log.Fatal(err)
	panic(err)
}

// ensures we have count copies of our first element
func (v *XMLValue) SetElementCountByCopyingFirstElementAsNeeded(count int) (err error) {
	elements := v.Elements()
	l := len(elements)
	switch {

	case l == 0:
		return fmt.Errorf("no elements found")

	case l < count:
		fmt.Printf("info: extending by %d elements", count-l)
		for i := l; i < count; i++ {
			err = v.Append(elements[0].Clone())
			if err != nil {
				return
			}
		}

	case l > count:
		fmt.Printf("info: truncating by %d elements", l-count)
		err = v.Truncate(count)
		if err != nil {
			return
		}
	}

	return
}

// we must already be a []any or we error
// warn: you must clone the element if you need this to be a clone and not that exact instance
func (v *XMLValue) Append(e ...any) (err error) {
	switch t := v.contents.(type) {

	case *XMLElement:
		// expand from a single element to an array of any
		v.contents = append([]any{t}, e...)

	case []any:
		// append slice
		v.contents = append(t, e...)

	default:
		err = fmt.Errorf("xmlvalue must be *XMLElement or []any, not %t", v.contents)
	}

	return
}

// we must already be a []any or this is an error
// we'll truncate to count elements (but keep all other things, such as comments and directives)
func (v *XMLValue) Truncate(count int) (err error) {
	switch t := v.contents.(type) {

	case []any:

		// create a new slice to hold what we're keeping
		keep := []any{}

		// build a new collection including count elements and whatever else was there before
		for i, j := 0, 0; i < len(t); i++ {
			switch t[i].(type) {
			case *XMLElement:
				if j < count {
					// keep it
					keep = append(keep, t[i])
					j++
				}
			default:
				// keep non-elements always
				keep = append(keep, t[i])
			}
		}
		v.contents = keep

	default:
		err = fmt.Errorf("xmlvalue must be []any, not %t", v.contents)
	}

	return
}

// we must already be a []any or we error
// warn: YOU MUST GIVE US INDEXES using ChildIndex, not from Elements()
func (v *XMLValue) InsertCopyOf(index, copy int) (err error) {
	switch t := v.contents.(type) {

	case []any:
		v.contents = etc.InsertAt(t, index, t[copy])

	default:
		err = fmt.Errorf("xmlvalue must be []any")
	}

	return
}

// we must already be a []any or we error
// WARN! we'll use the element you hand us, if you need to copy it, use e.Clone() when calling us!
// warn: YOU MUST GIVE US INDEXES using ChildIndex, not from Elements()
func (v *XMLValue) InsertAt(index int, e *XMLElement) (err error) {
	switch t := v.contents.(type) {

	case []any:
		v.contents = etc.InsertAt(t, index, any(e))

	default:
		err = fmt.Errorf("xmlvalue must be []any")
	}

	return
}

// we must already be a []any or this is an error
// warn: YOU MUST GIVE US INDEXES using ChildIndex, not from Elements()
func (v *XMLValue) RemoveSpan(startIndex int, count int) (err error) {
	switch t := v.contents.(type) {

	case []any:
		if startIndex < 0 || count < 0 || startIndex+count > len(t) {
			err = fmt.Errorf("index or count out of bounds")
			return
		}

		// a zero count is a special case of "do nothing"
		if count != 0 {
			v.contents = etc.RemoveSpanInSitu(t, startIndex, count)
		}

	default:
		err = fmt.Errorf("xmlvalue must be []any")
	}

	return
}

// warn: this doesn't work because it fails to clone the cells
// todo: fix this to clone the cells
// // extends our collection by inserting a run of count copies of the element at index
// // we must already be a []any or this is an error
// func (v *XMLValue) ExtendAt(index int, count int) (err error) {

// 	switch t := v.contents.(type) {

// 	case []any:
// 		v.contents = etc.InsertRunAt(t, index, count)

// 	default:
// 		err = fmt.Errorf("xmlvalue must be []any")
// 	}

// 	return
// }

// reorders the given child to the specified index
// warn: we must already be a []any or we error
// warn: YOU MUST GIVE US INDEXES using ChildIndex, not from Elements()
func (v *XMLValue) Reorder(from, to int) (err error) {
	switch t := v.contents.(type) {

	case []any:
		if from != to {
			// todo: we could optimize the shifted cells in the array for minimum copying
			// however, this is a slice of any, which are pointers, so not a biggie
			// todo: what would be really cool would be a generalized algo that could figure out the minimum moves to achieve end results from the whole list of changes
			e := t[from]                         // subtle: we need to grab e BEFORE we remove it!
			t = etc.RemoveSpanInSitu(t, from, 1) // subtle: this shouldn't change memory locations
			v.contents = etc.InsertAt(t, to, e)  // subtle: this shouldn't change memory locations
		}

	default:
		err = fmt.Errorf("xmlvalue must be []any")
	}

	return
}
