package wtg

import "testing"

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
