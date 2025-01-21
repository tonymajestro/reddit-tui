package components

import "slices"

type PageType int

type FocusStack []PageType

const (
	Home PageType = iota
	Subreddit
	Comments
	Quit
	Empty
)

func (f *FocusStack) Push(p PageType) {
	*f = append(*f, p)
}

func (f *FocusStack) Pop() PageType {
	if len(*f) == 0 {
		return Empty
	}

	last := (*f)[len(*f)-1]
	*f = slices.Delete(*f, len(*f)-1, len(*f))
	return last
}

func (f FocusStack) Peek() PageType {
	if len(f) == 0 {
		return Empty
	}

	return (f)[len(f)-1]
}

func (f *FocusStack) Clear() {
	*f = nil
}
