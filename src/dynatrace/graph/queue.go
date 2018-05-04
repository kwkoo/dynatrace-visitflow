package graph

type queue struct {
	q []*node
}

func newqueue() queue {
	return queue{q: make([]*node, 0)}
}

func (q *queue) isEmpty() bool {
	return len(q.q) == 0
}

func (q *queue) push(n *node) {
	q.q = append(q.q, n)
}

func (q *queue) pop() *node {
	if q.isEmpty() {
		return nil
	}
	lastIndex := len(q.q) - 1
	n := q.q[lastIndex]
	q.q = q.q[0:lastIndex]

	return n
}
