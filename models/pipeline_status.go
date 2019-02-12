package models

type PipelineStatusID string

const (
	Running PipelineStatusID = "running"
	Pending                  = "pending"
	Failed                   = "failed"
	Success                  = "success"
)

type PipelineStatus struct {
	ID     int64            `json:"id"`
	Ref    string           `json:"ref"`
	Sha    string           `json:"sha"`
	Status PipelineStatusID `json:"status"`
	WebURL string           `json:"web_url"`
}
