package graph

import (
	"bytes"
	"strconv"
)

type node struct {
	parent   *node
	index    int
	Name     string `json:"name"`
	children map[string]*node
	Value    int `json:"value"`
	depth    int
}

func (n node) String() string {
	var b bytes.Buffer
	for i := 0; i < n.depth; i++ {
		b.WriteString("\t")
	}
	b.WriteString("-> ")
	b.WriteString(strconv.Itoa(n.index))
	b.WriteString(" ")
	b.WriteString(n.Name)
	b.WriteString(" (")
	b.WriteString(strconv.Itoa(n.Value))
	b.WriteString(")")

	return b.String()
}
