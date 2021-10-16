package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Count int
	Head  *ListItem
	Tail  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.Count
}

func (l *list) Front() *ListItem {
	return l.Head
}

func (l *list) Back() *ListItem {
	return l.Tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{v, l.Head, nil}
	if l.Count == 0 {
		l.Tail = newListItem
	} else {
		newListItem.Next.Prev = newListItem
	}
	l.Head = newListItem
	l.Count++
	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{v, nil, l.Tail}
	if l.Count == 0 {
		l.Head = newListItem
	} else {
		newListItem.Prev.Next = newListItem
	}
	l.Tail = newListItem
	l.Count++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil {
		l.Tail = i.Prev
		l.Tail.Next = nil
	}
	if i.Prev == nil {
		l.Head = i.Next
		l.Head.Prev = nil
	}
	if i.Next != nil && i.Prev != nil {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.Count--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.Head == i {
		return
	}
	if l.Tail == i {
		i.Next = l.Head.Next
		l.Tail = i.Prev
		l.Tail.Next = nil
	} else {
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}

	l.Head.Prev = i
	l.Head.Next = i.Next
	i.Next = l.Head
	i.Prev = nil
	l.Head = i

}
