package models

type CreatePipelineResponse struct {
	BeforeSha   string      `json:"before_sha"`
	CommittedAt interface{} `json:"committed_at"`
	Coverage    interface{} `json:"coverage"`
	CreatedAt   string      `json:"created_at"`
	Duration    interface{} `json:"duration"`
	FinishedAt  interface{} `json:"finished_at"`
	ID          int64       `json:"id"`
	Ref         string      `json:"ref"`
	Sha         string      `json:"sha"`
	StartedAt   interface{} `json:"started_at"`
	Status      string      `json:"status"`
	Tag         bool        `json:"tag"`
	UpdatedAt   string      `json:"updated_at"`
	User        struct {
		AvatarURL string `json:"avatar_url"`
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		Username  string `json:"username"`
		WebURL    string `json:"web_url"`
	} `json:"user"`
	WebURL     string      `json:"web_url"`
	YamlErrors interface{} `json:"yaml_errors"`
}
