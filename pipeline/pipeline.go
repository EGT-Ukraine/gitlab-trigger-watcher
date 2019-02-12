package pipeline

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/EGT-Ukraine/gitlab-trigger-watcher/models"
	"github.com/pkg/errors"
)

const (
	DefaultHost       = "gitlab.com"
	httpClientTimeout = 10 * time.Second
)

type Schema uint8

const (
	HTTPS Schema = iota + 1
	HTTP
)

type Pipeline struct {
	client *http.Client
}

func New(tlsInsecureSkipVerify bool) *Pipeline {
	httpTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: httpClientTimeout,
		}).DialContext,
		TLSHandshakeTimeout: httpClientTimeout,
		DisableKeepAlives:   true,
		DisableCompression:  true,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
	}

	client := &http.Client{
		Transport: httpTransport,
		Timeout:   httpClientTimeout,
	}

	return &Pipeline{client: client}
}

func (p Pipeline) Run(schema Schema, host, token, ref string, projectID int, variables []string) (*models.CreatePipeline, error) {
	u, err := p.pipelineURL(schema, host, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "run pipeline failed")
	}

	vars := p.variables(variables)
	if token != "" {
		vars.Add("token", token)
	}
	if ref != "" {
		vars.Add("ref", ref)
	}

	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(vars.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "create request failed")
	}
	req.Header = p.defaultHeaders()

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "make http call on running the pipeline request failed")
	}

	var createPipelineModel *models.CreatePipeline
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()

		createPipelineModel = new(models.CreatePipeline)
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "read response data failed on running the pipeline request")
		}

		if err := json.Unmarshal(data, createPipelineModel); err != nil {
			return nil, errors.Wrap(err, "unmarshal response failed on running the pipeline request")
		}

		if createPipelineModel.ID == 0 {
			log.Printf("DEBUG: %s\n", string(data))
		}
	}

	return createPipelineModel, nil
}

func (p Pipeline) defaultHeaders() map[string][]string {
	return map[string][]string{
		"Content-Type": {"multipart/form-data"},
	}
}

func (p Pipeline) schemaName(schema Schema) string {
	switch schema {
	case HTTPS:
		return "https"
	case HTTP:
		return "http"
	default:
		return "https"
	}
}

// variable is a method to convert default string slice to url.Values representation
func (p Pipeline) variables(variables []string) url.Values {
	values := url.Values{}
	for _, variable := range variables {
		kv := strings.Split(variable, ":")
		if len(kv) == 2 {
			values.Add(fmt.Sprintf("variables[%s]", kv[0]), kv[1])
		}
	}

	return values
}

func (p Pipeline) pipelineURL(schema Schema, host string, projectID int) (string, error) {
	urlTpl := "{{.Schema}}://{{.Host}}/api/v4/projects/{{.ProjectID}}/trigger/pipeline"
	tpl, err := template.New("run-pipeline").Parse(urlTpl)
	if err != nil {
		return "", errors.New("failed to parse pipeline url template")
	}

	if host == "" {
		host = DefaultHost
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "run-pipeline", &struct {
		Schema    string
		Host      string
		ProjectID int
	}{Schema: p.schemaName(schema), Host: host, ProjectID: projectID}); err != nil {
		return "", errors.Wrap(err, "failed to execute pipeline url template")
	}

	return buf.String(), nil
}
