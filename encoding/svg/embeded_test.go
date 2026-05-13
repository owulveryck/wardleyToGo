package svgmap

import "os"

func Example_embededJS() {
	set := jsData{
		AllLinks: []string{"edge_a_b", "edge_b_c"},
	}
	_ = jsTmpl.Execute(os.Stdout, set)
}
