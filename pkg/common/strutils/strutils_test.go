package strutils

import "testing"

func TestAllSet(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "values are set",
			args: args{
				vals: []string{"foo", "bar"},
			},
			want: true,
		},
		{
			name: "not all values are set",
			args: args{
				vals: []string{"foo", ""},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllSet(tt.args.vals...); got != tt.want {
				t.Errorf("AllSet() error = %v, wantErr %v", got, tt.want)
			}
		})
	}
}
