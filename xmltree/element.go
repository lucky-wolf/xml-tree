package xmltree

import (
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
)

type SearchProperty struct {
	Tag   string
	Value string
}

// create a new element from scratch (it has no initial value)
// warn: you must set the value after creating it - it's invalid as-is
func MakeElement(name string) (e *XMLElement) {
	e = &XMLElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
	}
	return
}

// create a new element from scratch with the given initial value
// warn: value must be a valid value (see SetValue)
func MakeElementWithValue(name string, value any) (e *XMLElement) {
	e = &XMLElement{
		StartElement: xml.StartElement{Name: xml.Name{Local: name}},
	}
	e.SetValue(value)
	return
}

// clone an element
func (e *XMLElement) Clone() *XMLElement {
	return &XMLElement{
		StartElement: e.StartElement.Copy(),
		XMLValue:     e.XMLValue.Clone(),
	}
}

// clone an element and give it the specified local name
func (e *XMLElement) CloneAs(name string) (child *XMLElement) {
	child = &XMLElement{
		StartElement: e.StartElement.Copy(),
		XMLValue:     e.XMLValue.Clone(),
	}
	child.Name.Local = name
	return
}

// returns the first matching element from the list of elements based on tag (name)
func (e *XMLElement) Child(tag string) *XMLElement {
	for _, e = range e.Elements() {
		if e.Name.Local == tag {
			return e
		}
	}
	return nil
}

// returns the numeic value of the first matching element from the list of elements based on tag (name)
func (e *XMLElement) Float64ValueOf(tag string) (v float64, err error) {
	c := e.Child(tag)
	if c != nil {
		return c.GetNumericValue()
	}
	err = fmt.Errorf("no child with tag %s found in %s", tag, e.Name.Local)
	return
}

// returns the string value of the first matching element from the list of elements based on tag (name)
func (e *XMLElement) StringValueOf(tag string) (v string, err error) {
	c := e.Child(tag)
	if c != nil {
		var ok bool
		v, ok = c.GetStringValue()
		if ok {
			return
		}
		err = fmt.Errorf("child %s in %s is not a string", tag, e.Name.Local)
		return
	}
	err = fmt.Errorf("no child with tag %s found in %s", tag, e.Name.Local)
	return
}

// returns the int64 value of the first matching element from the list of elements based on tag (name)
func (e *XMLElement) Int64ValueOf(tag string) (v int64, err error) {
	c := e.Child(tag)
	if c != nil {
		return c.GetInt64Value()
	}
	err = fmt.Errorf("no child with tag %s found in %s", tag, e.Name.Local)
	return
}

// returns the first matching element from the list of elements based on tag (name) after setting it to the given value
// if the element is not found, it will be created and set to the given value
// if the element is found, it will be set to the given value
func (e *XMLElement) SetAppendChild(tag string, value any) (c *XMLElement, err error) {
	c = e.Child(tag)
	if c == nil {
		c = MakeElement(tag)
		err = e.Append(c)
		if err != nil {
			return
		}
	}
	c.SetValue(value)
	return
}

// returns the true index of the given element (index == -1 if not found)
func (e *XMLElement) ChildIndex(tag string) (index int) {

	switch v := e.contents.(type) {
	case []any:
		// multiple child elements
		for i, e := range v {
			switch v := e.(type) {
			case *XMLElement:
				if v.Name.Local == tag {
					index = i
					return
				}
			}
		}
	}

	return -1
}

// returns the true index of the zeroth element
func (e *XMLElement) ZeroElementIndex() (index int) {

	switch v := e.contents.(type) {
	case []any:
		// multiple child elements
		for i, e := range v {
			switch e.(type) {
			case *XMLElement:
				index = i
				return
			}
		}
	}

	return -1
}

// true if we match the given tag & value
func (e *XMLElement) Matches(tag, value string) bool {
	if e.Name.Local != tag {
		return false
	}
	v, ok := e.GetStringValue()
	return ok && v == value
}

// returns true if the given element has a sub element with specified tag and value
func (e *XMLElement) MatchesAll(properties ...SearchProperty) bool {
	for i := range properties {
		if !e.HasChildWithValue(properties[i].Tag, properties[i].Value) {
			return false
		}
	}
	return true
}

// returns the first matching element from the list of elements based on regex of tag-name
func (e *XMLElement) Matching(r *regexp.Regexp) (children []*XMLElement) {

	// scan our elements for those that match
	for _, e = range e.Elements() {
		if r.MatchString(e.Name.Local) {
			children = append(children, e)
		}
	}

	return
}

// returns the first matching child element whose tag and value equal the find tag and value
func (e *XMLElement) FindRecurse(tag string, value string) *XMLElement {

	// breadth first: check our contents for a match (non-recursive)
	if e.HasChildWithValue(tag, value) {
		// if we have one, we are the parent of this tag + value
		return e
	}

	// depth: now check each of our children for ownership of this tag+value item
	for _, child := range e.Elements() {
		child = child.FindRecurse(tag, value)
		if child != nil {
			// this child owns the target, so it is the result
			return child
		}
	}

	// could not find this tag + value parent
	return nil
}

// returns the first matching child element whose tag and value equal the find tag and value
func (e *XMLElement) ChildWithValue(tag string, value string) *XMLElement {
	c := e.Child(tag)
	if c != nil && c.StringValue() == value {
		return c
	}
	return nil
}

// returns the first matching child element whose tag and value equal the find tag and value
func (e *XMLElement) ChildWithValueOneOf(tag string, values ...string) *XMLElement {
	c := e.Child(tag)
	if c == nil {
		return nil
	}

	for _, value := range values {
		if c.StringValue() == value {
			return c
		}
	}

	return nil
}

// returns true if the given element has a sub element with specified tag and value
func (e *XMLElement) HasChildWithValue(tag string, value string) bool {
	return e.ChildWithValue(tag, value) != nil
}

// returns true if the given element has a sub element with specified tag and value
func (e *XMLElement) HasChildWithValueOneOf(tag string, values ...string) bool {
	return e.ChildWithValueOneOf(tag, values...) != nil
}

// like Find, but looks for anything with the value that matches as a prefix
func (e *XMLElement) HasPrefix(tag string, prefix string) bool {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueStartsWith(prefix) {
			return true
		}
	}
	return false
}

// like Find, but looks for anything with the value that matches as a suffix
func (e *XMLElement) HasSuffix(tag string, suffix string) bool {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueEndsWith(suffix) {
			return true
		}
	}
	return false
}

// visits each child with the given visitor function (aborts on error)
func (e *XMLElement) VisitChildren(visit func(*XMLElement) error) (err error) {

	for _, e := range e.Elements() {
		err = visit(e)
		if err != nil {
			return
		}
	}

	return
}

var ErrTargetNodeNotFound = errors.New("target doesn't have a node to copy to")

// copies the specified child from the given element to the ourself (replacing any we already have)
func (e *XMLElement) CopyByTag(tag string, from *XMLElement) (err error) {

	// get source
	source := from.Child(tag)
	if source == nil {
		err = fmt.Errorf("%s doesn't have a %s to copy from", from.Child("Name").StringValue(), tag)
		return
	}

	// get target
	target := e.Child(tag)
	if target == nil {
		err = fmt.Errorf("%s doesn't have a %s to copy to", e.Child("Name").StringValue(), tag)
		return
	}

	// clone source to target
	target.SetContents(source.CloneContents())

	return
}

// copies the specified child from the given element to the ourself (replacing any we already have)
// and visits the new elements
func (e *XMLElement) CopyAndVisitByTag(tag string, from *XMLElement, visit func(*XMLElement) error) (err error) {

	// copy by tag
	err = e.CopyByTag(tag, from)
	if err != nil {
		return
	}

	// visit the new elements
	err = e.Child(tag).VisitChildren(visit)
	if err != nil {
		return
	}

	return
}

// sets existing element value if present, returns that element or nil
func (e *XMLElement) SetChildValue(tag string, value any) (c *XMLElement, err error) {
	c = e.Child(tag)
	if c != nil {
		c.SetValue(value)
	} else if value != "" && value != "0" && value != 0 && value != 0.0 {
		err = fmt.Errorf("failed to set %s to %v", tag, value)
		return
	}
	return
}

// sets existing element value if present, or creates it (appends to parent)
func (e *XMLElement) EnsureChildValue(tag string, value any) (c *XMLElement, err error) {

	c = e.Child(tag)

	if c != nil {
		c.SetValue(value)
		return
	}

	// note: we could elide the node if we knew for a fact that the default value is the same as what we're setting it to here...
	c = MakeElementWithValue(tag, value)
	err = e.Append(c)

	return
}

// updates it to be scaled by the given input
func (e *XMLElement) ScaleChildBy(tag string, scale float64) (err error) {

	// scaling by 1 is a noop
	if scale == 1.0 {
		return
	}

	// child must exist for this to be possible
	child := e.Child(tag)
	if child == nil {
		err = fmt.Errorf("no child %s to scale by %f", tag, scale)
		return
	}

	// do it
	return child.ScaleBy(scale)
}

// updates it to be scaled by the given input
func (e *XMLElement) AdjustChildBy(tag string, adjustment float64) (err error) {

	// offset by 0 is a noop
	if adjustment == 0 {
		return
	}

	// child must exist for this to be possible
	child := e.Child(tag)
	if child == nil {
		err = fmt.Errorf("no child %s to adjust by %f", tag, adjustment)
		return
	}

	// do it
	return child.AdjustValue(adjustment)
}

// sets one value to be that of another (both must be value types)
func (e *XMLElement) SetChildToSibling(child, sibling string) (err error) {

	// no sibling = nothing to do
	sib := e.Child(sibling)
	if sib == nil {
		err = fmt.Errorf("no sibling %s to set %s to", sibling, child)
		return
	}

	// this will ignore the set if the value is effectively a noop & the child doesn't exist
	// it will panic however if the sibling has a non-zero non-blank value but the child fails to exist
	_, err = e.SetChildValue(child, sib.StringValue())
	if err != nil {
		return
	}
	return
}

// updates it to be scaled by the given input
func (e *XMLElement) ScaleChildToSiblingBy(tag, siblingtag string, scale float64) (err error) {

	// sibling must exist for this to be possible
	sibling := e.Child(siblingtag)
	if sibling == nil {
		err = fmt.Errorf("no sibling %s to scale %s by", siblingtag, tag)
		return
	}

	// child must exist for this to be possible
	child := e.Child(tag)
	if child == nil {
		err = fmt.Errorf("no child %s to scale by %s", tag, siblingtag)
		return
	}

	child.SetValue(sibling.FloatValue() * scale)
	return
}

// updates it to be offset by the given input
func (e *XMLElement) AdjustChildToSiblingBy(child, sibling string, adj float64) (err error) {
	sib := e.Child(sibling)
	if sib == nil {
		err = fmt.Errorf("no sibling %s to adjust %s by", sibling, child)
		return
	}
	e.Child(child).SetValue(sib.FloatValue() + adj)
	return
}

// removes the specified child from the given element if present
func (e *XMLElement) RemoveByTag(attr string) (err error) {
	index := e.ChildIndex(attr)
	if index == -1 {
		return
	}
	return e.RemoveSpan(index, 1)
}

// removes the child at the specified index (must be a valid index from ChildIndex)
// returns an error if the index is out of bounds
// warn: YOU MUST GIVE US INDEXES using ChildIndex, not from Elements()
func (e *XMLElement) RemoveChildAt(index int) (err error) {
	if index == -1 {
		return fmt.Errorf("invalid child index: -1")
	}
	return e.RemoveSpan(index, 1)
}

// simple way to verify that this element is of the given kind
func (e *XMLElement) Is(kind string) bool {
	return e.Name.Local == kind
}

// simple way to verify that this element is of the given kind
func (e *XMLElement) MustBe(kind string) (err error) {
	if !e.Is(kind) {
		err = fmt.Errorf("invalid element: expected %s but found %s", kind, e.Name.Local)
	}
	return
}

// returns the specified attribute value (ok is true if this attribute was found)
func (e *XMLElement) Attribute(name string) (value string, ok bool) {
	for _, attr := range e.Attr {
		if attr.Name.Local == name {
			value = attr.Value
			ok = true
			return
		}
	}
	return
}
