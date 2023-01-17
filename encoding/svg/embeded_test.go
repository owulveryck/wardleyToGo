package svgmap

import "os"

func ExampleEmbededJS() {
	set := jsData{
		AllLinks: []string{"edge_a_b", "edge_b_c"},
		G: map[string][]string{
			"a": []string{
				"edge_a_b",
				"edge_b_c",
			},
			"d": []string{
				"edge_d_e",
				"edge_e_f",
			},
		},
	}
	jsTmpl.Execute(os.Stdout, set)
}
