package trigger

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
		args args
		want string
	}{
		{
			name: "success with HTTPS",
			args: struct {
				schema    Schema
				host      string
				projectID int
			}{schema: HTTPS, host: "gitlab.com", projectID: 123},
			want: "https://gitlab.com/api/v4/projects/123/trigger/pipeline",
		}, {
			name: "success with HTTP and another host",
			args: struct {
				schema    Schema
				host      string
				projectID int
			}{schema: HTTP, host: "gitlab1.com", projectID: 123},
			want: "http://gitlab1.com/api/v4/projects/123/trigger/pipeline",
		},
	}
	for _, tt := range tests {
		p := new(Trigger)
		p.schema = tt.args.schema
		p.host = tt.args.host
		p.projectID = tt.args.projectID

		t.Run(tt.name, func(t *testing.T) {
			tpl, err := p.createPipelineTpl()
			assert.Nil(t, err)

			u, err := p.urlByTemplate(tpl)
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
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "success",
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
			name: "success with double separator",
			args: struct {
				variables []string
			}{
				variables: []string{
					"variable1:value1:prefix",
					"variable2:value2:prefix",
				},
			},
			want: map[string][]string{
				"variables[variable1]": {
					"value1:prefix",
				},
				"variables[variable2]": {
					"value2:prefix",
				},
			},
		}, {
			name: "fail",
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
			p := &Trigger{variables: tt.args.variables}
			if !tt.wantErr {
				assert.Equal(t, tt.want, p.urlVariables())
				return
			}
			assert.NotEqual(t, tt.want, p.urlVariables())
		})
	}
}
