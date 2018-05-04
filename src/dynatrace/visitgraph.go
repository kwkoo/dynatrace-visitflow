package dynatrace

import (
	"dynatrace/graph"
	"fmt"
	"io"
)

var actionsgraph graph.Graph

// InitGraph initializes the internal data structures.
func InitGraph() {
	actionsgraph = graph.NewGraph()
}

// AddVisit adds the user actions for the visit into the graph.
func AddVisit(v Visit, app string) {
	actionNames := make([]string, 0)
	for _, action := range v.UserActions {
		if len(app) > 0 && app != action.Application {
			continue
		}
		actionNames = append(actionNames, action.Name)
	}
	actionsgraph.AddNodes(actionNames)
}

// PrintGraph prints out the graph in plain format.
func PrintGraph(w io.Writer) {
	actionsgraph.DumpNodes(w)
}

// GenerateFlow generates the JSON for the flow diagram.
func GenerateFlow(w io.Writer, maxchildren, linkvaluethreshold int) {
	newgraph := actionsgraph.GetTrimmedCopy(maxchildren, linkvaluethreshold)
	fmt.Fprint(w, "{\"visitCount\":")
	fmt.Fprint(w, newgraph.GetRootValue())
	fmt.Fprint(w, ",\"nodes\":")
	newgraph.StreamNodes(w)
	fmt.Fprint(w, ",\"links\":")
	newgraph.StreamLinks(w)
	fmt.Fprintln(w, "}")
}
