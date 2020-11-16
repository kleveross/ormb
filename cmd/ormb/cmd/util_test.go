package cmd

import "testing"

func Test_convertRef(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "two level",
			args: args{
				ref: "goharbor.io/release:v1",
			},
			want: "goharbor.io/release/release:v1",
		},
		{
			name: "three level",
			args: args{
				ref: "goharbor.io/release/release:v1",
			},
			want: "goharbor.io/release/release:v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertRef(tt.args.ref); got != tt.want {
				t.Errorf("convertRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
