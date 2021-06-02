package sliceutils

import "testing"

func TestContainsStr(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "find existing elem",
			args: args{
				s:   []string{"a", "b", "c"},
				str: "a",
			},
			want: true,
		},
		{
			name: "searching non existing elem",
			args: args{
				s:   []string{"a", "b", "c"},
				str: "d",
			},
			want: false,
		},
		{
			name: "searching empty elem",
			args: args{
				s:   []string{"a", "b", "c"},
				str: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsStr(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("ContainsStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
