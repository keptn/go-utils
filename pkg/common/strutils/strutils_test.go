package strutils

import "testing"

func TestAllSet(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "values are set",
			args: args{
				vals: []string{"foo", "bar"},
			},
			wantErr: false,
		},
		{
			name: "not all values are set",
			args: args{
				vals: []string{"foo", ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AllSet(tt.args.vals...); (err != nil) != tt.wantErr {
				t.Errorf("AllSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
