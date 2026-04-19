package svgmap

import "os"

func Example_embededJS() {
	set := jsData{
		AllLinks: []string{"edge_a_b", "edge_b_c"},
		G: []jsGraphEntry{
			{Key: "a", Edges: []string{"edge_a_b", "edge_b_c"}},
			{Key: "d", Edges: []string{"edge_d_e", "edge_e_f"}},
		},
	}
	_ = jsTmpl.Execute(os.Stdout, set)
}
