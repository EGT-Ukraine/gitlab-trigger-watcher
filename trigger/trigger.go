package trigger

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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

var (
	ErrNoServerResponse = errors.New("no response from server")
)

type Trigger struct {
	client                         *http.Client
	schema                         Schema
	host, privateToken, token, ref string
	projectID                      int
	variables                      []string
}

func New(tlsInsecureSkipVerify bool, schema Schema, host, privateToken, token, ref string, projectID int, variables []string) *Trigger {
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

	return &Trigger{client: client, schema: schema, host: host, privateToken: privateToken, token: token, ref: ref, projectID: projectID, variables: variables}
}

func (p Trigger) RunPipeline() (*models.CreatePipelineResponse, error) {
	tpl, err := p.createPipelineTpl()
	if err != nil {
		return nil, errors.Wrap(err, "run pipeline failed")
	}

	u, err := p.urlByTemplate(tpl)
	if err != nil {
		return nil, err
	}

	vars := p.urlVariables()
	if p.token != "" {
		vars.Add("token", p.token)
	}
	if p.ref != "" {
		vars.Add("ref", p.ref)
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

	if resp == nil || resp.Body == nil {
		return nil, ErrNoServerResponse
	}
	defer resp.Body.Close()

	createPipelineModel := new(models.CreatePipelineResponse)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response data failed on running the pipeline request")
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(data))
	}

	if err := json.Unmarshal(data, createPipelineModel); err != nil {
		return nil, errors.Wrap(err, "unmarshal response failed on running the pipeline request")
	}

	return createPipelineModel, nil
}

func (p Trigger) PollForCompletion(pipelineID int64) (*models.PipelineStatusResponse, error) {
	tpl, err := p.pollPipelineTpl()
	if err != nil {
		return nil, errors.Wrap(err, "poll pipeline failed")
	}

	u, err := p.urlByTemplate(tpl)
	if err != nil {
		return nil, err
	}

	headers := p.defaultHeaders()
	headers["PRIVATE-TOKEN"] = []string{p.privateToken}

	req, err := http.NewRequest(http.MethodGet, u, strings.NewReader(p.urlVariables().Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "create request failed")
	}
	req.Header = headers

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "make http call on polling the pipeline request failed")
	}

	if resp == nil || resp.Body == nil {
		return nil, ErrNoServerResponse
	}
	defer resp.Body.Close()

	var pipelineStatusesModel []*models.PipelineStatusResponse
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response data failed on polling the pipeline request")
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(data))
	}

	if err := json.Unmarshal(data, &pipelineStatusesModel); err != nil {
		return nil, errors.Wrap(err, "unmarshal response failed on polling the pipeline request")
	}

	for _, pipelineStatus := range pipelineStatusesModel {
		if pipelineStatus.ID == pipelineID {
			return pipelineStatus, nil
		}
	}

	return nil, errors.Errorf("no such pipelineID found: %d", pipelineID)
}

func (p Trigger) defaultHeaders() map[string][]string {
	return map[string][]string{
		"Content-Type": {"multipart/form-data"},
	}
}

func (p Trigger) schemaName(schema Schema) string {
	switch schema {
	case HTTPS:
		return "https"
	case HTTP:
		return "http"
	default:
		return "https"
	}
}

// urlVariables is a method to convert default string slice to url.Values representation
func (p Trigger) urlVariables() url.Values {
	values := url.Values{}
	for _, variable := range p.variables {
		kv := strings.SplitN(variable, ":", 2)
		if len(kv) == 2 {
			values.Add(fmt.Sprintf("variables[%s]", kv[0]), kv[1])
		}
	}

	return values
}

func (p Trigger) createPipelineTpl() (*template.Template, error) {
	urlTpl := "{{.Schema}}://{{.Host}}/api/v4/projects/{{.ProjectID}}/trigger/pipeline"
	tpl, err := template.New("tpl").Parse(urlTpl)
	if err != nil {
		return nil, errors.New("failed to parse pipeline url template")
	}

	return tpl, nil
}

func (p Trigger) pollPipelineTpl() (*template.Template, error) {
	urlTpl := "{{.Schema}}://{{.Host}}/api/v4/projects/{{.ProjectID}}/pipelines"
	tpl, err := template.New("tpl").Parse(urlTpl)
	if err != nil {
		return nil, errors.New("failed to parse pipeline url template")
	}

	return tpl, nil
}

func (p *Trigger) urlByTemplate(tpl *template.Template) (string, error) {
	if p.host == "" {
		p.host = DefaultHost
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "tpl", &struct {
		Schema    string
		Host      string
		ProjectID int
	}{Schema: p.schemaName(p.schema), Host: p.host, ProjectID: p.projectID}); err != nil {
		return "", errors.Wrap(err, "failed to execute pipeline url template")
	}

	return buf.String(), nil
}
