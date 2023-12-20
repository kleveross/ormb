package oci

import (
	"reflect"
	"testing"
)

func TestParseReference(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Reference
		wantErr bool
	}{
		{
			name:    "Test single digit tag (from issue #209).",
			args:    args{s: "foo:1"},
			want:    &Reference{Tag: "1", Repo: "foo"},
			wantErr: false,
		},
		{
			name:    "Test mulitple digit-only tag.",
			args:    args{s: "registry.example.com/project/repo:1612367"},
			want:    &Reference{Tag: "1612367", Repo: "registry.example.com/project/repo"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseReference(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseReference() = %v, want %v", got, tt.want)
			}
		})
	}
}
