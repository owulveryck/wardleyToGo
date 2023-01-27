package svg

import (
	"reflect"
	"testing"
)

func Test_splitString(t *testing.T) {
	type args struct {
		s   string
		max int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"empty",
			args{"", 1},
			[]string{""},
		},
		{
			"a",
			args{"a", 1},
			[]string{"a"},
		},
		{
			"aa",
			args{"aa", 1},
			[]string{"aa"},
		},
		{
			"aa  ",
			args{"aa  ", 1},
			[]string{"aa"},
		},
		{
			"aa  ",
			args{"aa  ", 2},
			[]string{"aa"},
		},
		{
			"one   two    three   four  ",
			args{"one   two    three   four  ", 2},
			[]string{"one", "two", "three", "four"},
		},
		{
			"one   two    three   four  ",
			args{"one   two    three   four  ", 2},
			[]string{"one", "two", "three", "four"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitString(tt.args.s, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitString() = %v, want %v", got, tt.want)
			}
		})
	}
}
