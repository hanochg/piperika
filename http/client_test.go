package http

import (
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"reflect"
	"testing"
)

func TestNewPipelineHttp(t *testing.T) {
	type args struct {
		c *components.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *pipelineHttpClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPipelineHttp(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPipelineHttp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPipelineHttp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
