package util

import (
	"os"
	"path"
	"testing"

	ormbmodel "github.com/kleveross/ormb/pkg/model"
)

func TestWriteORMBFile(t *testing.T) {
	cwd, _ := os.Getwd()

	ormbfilePath := path.Join(cwd, "ormbfile.yaml")
	defer os.RemoveAll(ormbfilePath)

	type args struct {
		filePath string
		format   ormbmodel.Format
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "WriteORMBFile",
			args: args{
				filePath: ormbfilePath,
				format:   ormbmodel.FormatMXNetParams,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteORMBFile(tt.args.filePath, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("WriteORMBFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
