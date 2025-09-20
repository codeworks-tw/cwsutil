package cwsbase

type (
	// Stack is a generic LIFO (Last In, First Out) data structure
	Stack[T any] struct {
		// top points to the top node in the stack
		top    *node[T]
		// length stores the current number of items in the stack
		length int
	}
	// node represents an internal node in the stack
	node[T any] struct {
		// value holds the actual data
		value *T
		// prev points to the previous node in the stack
		prev  *node[T]
	}
)

// New creates a new empty stack of type T
func New[T any]() Stack[T] {
	return Stack[T]{nil, 0}
}

// Len returns the number of items currently in the stack
func (this *Stack[T]) Len() int {
	return this.length
}

// Peek returns a pointer to the top item on the stack without removing it
// Returns nil if the stack is empty
func (this *Stack[T]) Peek() *T {
	if this.length == 0 {
		return nil
	}
	return this.top.value
}

// Pop removes and returns a pointer to the top item of the stack
// Returns nil if the stack is empty
func (this *Stack[T]) Pop() *T {
	if this.length == 0 {
		return nil
	}

	n := this.top
	this.top = n.prev
	this.length--
	return n.value
}

// Push adds a new value to the top of the stack
func (this *Stack[T]) Push(value *T) {
	n := &node[T]{value, this.top}
	this.top = n
	this.length++
}
