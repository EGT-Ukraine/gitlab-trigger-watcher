package pipeline

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pipelineURL(t *testing.T) {
	type args struct {
		schema    Schema
		host      string
		projectID int
	}
	tests := []struct {
		name string
		p    *Pipeline
		args args
		want string
	}{
		{
			name: "success with HTTPS",
			p:    new(Pipeline),
			args: struct {
				schema    Schema
				host      string
				projectID int
			}{schema: HTTPS, host: "gitlab.com", projectID: 123},
			want: "https://gitlab.com/api/v4/projects/123/trigger/pipeline",
		}, {
			name: "success with HTTP and another host",
			p:    new(Pipeline),
			args: struct {
				schema    Schema
				host      string
				projectID int
			}{schema: HTTP, host: "gitlab1.com", projectID: 123},
			want: "http://gitlab1.com/api/v4/projects/123/trigger/pipeline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := tt.p.pipelineURL(tt.args.schema, tt.args.host, tt.args.projectID)
			assert.Nil(t, err)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, u)
		})
	}
}

func TestPipeline_variables(t *testing.T) {
	type args struct {
		variables []string
	}
	tests := []struct {
		name    string
		p       *Pipeline
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "success",
			p:    new(Pipeline),
			args: struct {
				variables []string
			}{
				variables: []string{
					"variable1:value1",
					"variable2:value2",
				},
			},
			want: map[string][]string{
				"variables[variable1]": {
					"value1",
				},
				"variables[variable2]": {
					"value2",
				},
			},
		}, {
			name: "fail",
			p:    new(Pipeline),
			args: struct {
				variables []string
			}{
				variables: []string{
					"variable1:value1",
					"variable2:value2",
				},
			},
			want: map[string][]string{
				"variables[variable3]": {
					"value3",
				},
				"variables[variable4]": {
					"value4",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pipeline{}
			if !tt.wantErr {
				assert.Equal(t, p.variables(tt.args.variables), tt.want)
				return
			}
			assert.NotEqual(t, p.variables(tt.args.variables), tt.want)
		})
	}
}
