package wtg

import (
	"testing"

	"github.com/owulveryck/wardleyToGo/components/wardley"
)

func Test_computeEvolutionPosition(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{
			"error nil string",
			args{s: "bad"},
			0,
			0,
			true,
		}, { //		{
			"too many |",
			args{s: "|...|...|...|...|...|"},
			0,
			0,
			true,
		}, { //		{
			"two x",
			args{s: "|.x.|.x.|...|..."},
			0,
			0,
			true,
		}, { //		{
			"missing a |",
			args{s: "|...|...|...|..."},
			0,
			0,
			true,
		}, {
			"empty",
			args{s: "|...|...|...|...|"},
			0,
			0,
			true,
		}, {
			"simple 1",
			args{s: "|.x.|...|...|...|"},
			9,
			0,
			false,
		}, {
			"simple 1.1",
			args{s: "|x|...|...|...|"},
			9,
			0,
			false,
		}, {
			"simple 2.0",
			args{s: "||.x.|...|...|"},
			29,
			0,
			false,
		}, {
			"simple 2.1",
			args{s: "||x|...|...|"},
			29,
			0,
			false,
		}, {
			"evolution with no x",
			args{s: "||>|...|...|"},
			0,
			0,
			true,
		}, {
			"two evolutions with",
			args{s: "|x|>|.>.|...|"},
			0,
			0,
			true,
		}, {
			"negative evolutions",
			args{s: "|||.>.|.x.|"},
			0,
			0,
			true,
		}, {
			"ok with evolutions",
			args{s: "|x|>|...|...|"},
			9,
			29,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := computeEvolutionPosition(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeEvolutionPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("computeEvolutionPosition() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("computeEvolutionPosition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestParser_computeY(t *testing.T) {
	t.Run("simple two nodes", computeYsimpleTwoNodes)
	t.Run("simple two nodes with evolution", computeYsimpleTwoNodesWithEvolution)
	t.Run("simple three nodes", computeYsimpleThreeNodes)
	t.Run("three nodes with different visibility", computeYThreeNodes)
	t.Run("four nodes with different visibility", computeYFourNodes)
}

func computeYsimpleTwoNodesWithEvolution(t *testing.T) {
	p := NewParser()
	c1 := wardley.NewComponent(1)
	c2 := wardley.NewComponent(2)
	c22 := wardley.NewEvolvedComponent(22)
	p.WMap.AddComponent(c1)
	p.WMap.AddComponent(c2)
	p.WMap.AddComponent(c22)
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c2,
		Visibility: 1,
	})
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:    c2,
		T:    c22,
		Type: wardley.EvolvedComponentEdge,
	})
	SetCoords(*p.WMap, false)
	if c1.GetPosition().Y != 3 {
		t.Errorf("expected position of c1 to be 3, but is %v", c1.GetPosition().Y)
	}
	if c2.GetPosition().Y != 97 {
		t.Errorf("expected position of c2 to be 97, but is %v", c2.GetPosition().Y)
	}
	if c22.GetPosition().Y != 97 {
		t.Errorf("expected position of c22 to be 97, but is %v", c22.GetPosition().Y)
	}
}
func computeYsimpleThreeNodes(t *testing.T) {
	p := NewParser()
	c1 := wardley.NewComponent(1)
	c2 := wardley.NewComponent(2)
	c3 := wardley.NewComponent(3)
	p.WMap.AddComponent(c1)
	p.WMap.AddComponent(c2)
	p.WMap.AddComponent(c3)
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c2,
		Visibility: 1,
	})
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c2,
		T:          c3,
		Visibility: 1,
	})
	SetCoords(*p.WMap, false)
	if c1.GetPosition().Y != 3 {
		t.Errorf("expected position of c1 to be 3, but is %v", c1.GetPosition().Y)
	}
	if c2.GetPosition().Y != 50 {
		t.Errorf("expected position of c2 to be 50, but is %v", c2.GetPosition().Y)
	}
	if c3.GetPosition().Y != 97 {
		t.Errorf("expected position of c3 to be 97, but is %v", c3.GetPosition().Y)
	}
}
func computeYsimpleTwoNodes(t *testing.T) {
	p := NewParser()
	c1 := wardley.NewComponent(1)
	c2 := wardley.NewComponent(2)
	p.WMap.AddComponent(c1)
	p.WMap.AddComponent(c2)
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c2,
		Visibility: 1,
	})
	SetCoords(*p.WMap, false)
	if c1.GetPosition().Y != 3 {
		t.Errorf("expected position of c1 to be 3, but is %v", c1.GetPosition().Y)
	}
	if c2.GetPosition().Y != 97 {
		t.Errorf("expected position of c2 to be 97, but is %v", c2.GetPosition().Y)
	}
}
func computeYThreeNodes(t *testing.T) {
	p := NewParser()
	c1 := wardley.NewComponent(1)
	c2 := wardley.NewComponent(2)
	c3 := wardley.NewComponent(3)
	p.WMap.AddComponent(c1)
	p.WMap.AddComponent(c2)
	p.WMap.AddComponent(c3)
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c2,
		Visibility: 1,
	})
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c3,
		Visibility: 3,
	})
	SetCoords(*p.WMap, false)
	if c1.GetPosition().Y != 3 {
		t.Errorf("expected position of c1 to be 3, but is %v", c1.GetPosition().Y)
	}
	if c2.GetPosition().Y != 34 {
		t.Errorf("expected position of c2 to be 34, but is %v", c2.GetPosition().Y)
	}
	if c3.GetPosition().Y != 96 {
		t.Errorf("expected position of c3 to be 96, but is %v", c3.GetPosition().Y)
	}
}
func computeYFourNodes(t *testing.T) {
	p := NewParser()
	c1 := wardley.NewComponent(1)
	c2 := wardley.NewComponent(2)
	c3 := wardley.NewComponent(3)
	c4 := wardley.NewComponent(4)
	p.WMap.AddComponent(c1)
	p.WMap.AddComponent(c2)
	p.WMap.AddComponent(c3)
	p.WMap.AddComponent(c4)
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c2,
		Visibility: 1,
	})
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c1,
		T:          c3,
		Visibility: 3,
	})
	p.WMap.SetCollaboration(&wardley.Collaboration{
		F:          c3,
		T:          c4,
		Visibility: 1,
	})
	SetCoords(*p.WMap, false)
	if c1.GetPosition().Y != 3 {
		t.Errorf("expected position of c1 to be 3, but is %v", c1.GetPosition().Y)
	}
	if c2.GetPosition().Y != 26 {
		t.Errorf("expected position of c2 to be 26, but is %v", c2.GetPosition().Y)
	}
	if c3.GetPosition().Y != 72 {
		t.Errorf("expected position of c3 to be 72, but is %v", c3.GetPosition().Y)
	}
	if c4.GetPosition().Y != 95 {
		t.Errorf("expected position of c4 to be 95, but is %v", c4.GetPosition().Y)
	}
}
