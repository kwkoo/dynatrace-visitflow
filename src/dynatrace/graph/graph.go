package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type link struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}

// Graph structure contains all nodes.
type Graph struct {
	nodes []*node
}

// NewGraph creates a new graph structure.
func NewGraph() Graph {
	root := node{
		parent:   nil,
		Name:     "Root",
		children: make(map[string]*node, 0),
	}
	return Graph{nodes: []*node{&root}}
}

// GetRootValue gets the value of the root node.
func (g Graph) GetRootValue() int {
	if len(g.nodes) == 0 {
		return 0
	}
	return g.nodes[0].Value
}

// NewCopy creates a new copy of the graph
func (g *Graph) NewCopy() Graph {
	g.updateNodeIndexes()
	newcopy := Graph{nodes: make([]*node, 0)}
	for i, srcnode := range g.nodes {
		newcopy.nodes = append(newcopy.nodes, &node{
			index:    i,
			Name:     srcnode.Name,
			children: make(map[string]*node),
			Value:    srcnode.Value,
			depth:    srcnode.depth,
		})
	}
	for i, src := range g.nodes {
		if src.parent != nil {
			newcopy.nodes[i].parent = newcopy.nodes[src.parent.index]
		}
		for name, child := range src.children {
			newcopy.nodes[i].children[name] = newcopy.nodes[child.index]
		}
	}
	return newcopy
}

// AddNodes traverses the graph and adds values to nodes, creating new nodes if necessary.
func (g *Graph) AddNodes(names []string) {
	if len(names) == 0 {
		return
	}
	current := g.nodes[0]
	current.Value++
	for _, name := range names {
		child, ok := current.children[name]
		if !ok {
			child = &node{
				parent:   current,
				Name:     name,
				children: make(map[string]*node, 0),
				depth:    current.depth + 1,
			}
			g.nodes = append(g.nodes, child)
			current.children[name] = child
		}
		child.Value++
		current = child
	}
}

// GetTrimmedCopy creates a copy of the graph and trims it according to the given criteria.
func (g *Graph) GetTrimmedCopy(maxchildren, linkvaluethreshold int) Graph {
	newcopy := g.NewCopy()
	newcopy.trimNodes(maxchildren, linkvaluethreshold)
	return newcopy
}

// DumpNodes lists out nodes in a plain format.
func (g Graph) DumpNodes(w io.Writer) {
	if len(g.nodes) == 0 {
		return
	}
	g.updateNodeIndexes()
	q := newqueue()
	q.push(g.nodes[0])
	for !q.isEmpty() {
		n := q.pop()
		fmt.Fprintln(w, n.String())
		for _, c := range n.children {
			q.push(c)
		}
	}
}

// StreamNodes streams nodes in JSON.
func (g Graph) StreamNodes(w io.Writer) {
	if len(g.nodes) == 0 {
		return
	}
	g.updateNodeIndexes()
	enc := json.NewEncoder(w)
	enc.Encode(g.nodes)
}

// StreamLinks streams links in JSON. It is assumed that StreamNodes is called right before StreamLinks - StreamLinks does not update node indexes.
func (g Graph) StreamLinks(w io.Writer) {
	if len(g.nodes) == 0 {
		return
	}

	links := make([]link, 0)
	for _, src := range g.nodes {
		srcIndex := src.index
		for _, dst := range src.children {
			links = append(links, link{
				Source: srcIndex,
				Target: dst.index,
				Value:  dst.Value,
			})
		}
	}
	enc := json.NewEncoder(w)
	enc.Encode(links)
}

func (g *Graph) updateNodeIndexes() {
	for i, n := range g.nodes {
		n.index = i
	}
}

// Do a depth-first search traverse, leaving out nodes which don't fit.
func (g *Graph) trimNodes(maxchildren, linkvaluethreshold int) {
	if len(g.nodes) == 0 {
		return
	}
	q := newqueue()
	root := g.nodes[0]
	q.push(root)

	// we need to rebuild g.nodes
	newnodes := make([]*node, 0)
	var counter int

	for !q.isEmpty() {
		n := q.pop()

		// those that go in the queue make the cut
		n.index = counter
		counter++
		newnodes = append(newnodes, n)

		// delete children who don't make the linkvaluethreshold
		for childsname, child := range n.children {
			if child.Value < linkvaluethreshold {
				delete(n.children, childsname)
			}
		}

		// if there are too many children, sort according to value and keep the top few
		if len(n.children) > maxchildren {
			slice := make([]*node, 0)
			for _, child := range n.children {
				slice = append(slice, child)
			}
			// sort in descending order
			sort.SliceStable(slice, func(i, j int) bool { return (*slice[i]).Value > (*slice[j]).Value })
			for _, child := range slice[maxchildren:] {
				delete(n.children, child.Name)
			}
		}

		for _, child := range n.children {
			q.push(child)
		}
	}

	g.nodes = newnodes
}
