package models

type PipelineStatus string

const (
	Running PipelineStatus = "running"
	Pending                = "pending"
	Failed                 = "failed"
	Success                = "success"
)

type PipelineStatusResponse struct {
	ID     int64          `json:"id"`
	Ref    string         `json:"ref"`
	Sha    string         `json:"sha"`
	Status PipelineStatus `json:"status"`
	WebURL string         `json:"web_url"`
}
