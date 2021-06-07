package api

import (
	"fmt"
)

// Job represents a Buildkite Agent API Job
type Job struct {
	ID                 string            `json:"id,omitempty"`
	Endpoint           string            `json:"endpoint"`
	State              string            `json:"state,omitempty"`
	Env                map[string]string `json:"env,omitempty"`
	ChunksMaxSizeBytes int               `json:"chunks_max_size_bytes,omitempty"`
	ExitStatus         string            `json:"exit_status,omitempty"`
	Signal             string            `json:"signal,omitempty"`
	SignalReason       string            `json:"signal_reason,omitempty"`
	StartedAt          string            `json:"started_at,omitempty"`
	FinishedAt         string            `json:"finished_at,omitempty"`
	RunnableAt         string            `json:"runnable_at,omitempty"`
	ChunksFailedCount  int               `json:"chunks_failed_count,omitempty"`
}

type JobState struct {
	State string `json:"state,omitempty"`
}

type jobStartRequest struct {
	StartedAt string `json:"started_at,omitempty"`
}

type jobFinishRequest struct {
	ExitStatus        string `json:"exit_status,omitempty"`
	Signal            string `json:"signal,omitempty"`
	SignalReason      string `json:"signal_reason,omitempty"`
	FinishedAt        string `json:"finished_at,omitempty"`
	ChunksFailedCount int    `json:"chunks_failed_count"`
}

// GetJobState returns the state of a given job
func (c *Client) GetJobState(id string) (*JobState, *Response, error) {
	u := fmt.Sprintf("jobs/%s", id)

	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	s := new(JobState)
	resp, err := c.doRequest(req, s)
	if err != nil {
		return nil, resp, err
	}

	return s, resp, err
}

// Acquires a job using its ID
func (c *Client) AcquireJob(id string) (*Job, *Response, error) {
	u := fmt.Sprintf("jobs/%s/acquire", id)

	req, err := c.newRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}

	j := new(Job)
	resp, err := c.doRequest(req, j)
	if err != nil {
		return nil, resp, err
	}

	return j, resp, err
}

// AcceptJob accepts the passed in job. Returns the job with its finalized set of
// environment variables (when a job is accepted, the agents environment is
// applied to the job)
func (c *Client) AcceptJob(job *Job) (*Job, *Response, error) {
	u := fmt.Sprintf("jobs/%s/accept", job.ID)

	req, err := c.newRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}

	j := new(Job)
	resp, err := c.doRequest(req, j)
	if err != nil {
		return nil, resp, err
	}

	return j, resp, err
}

// StartJob starts the passed in job
func (c *Client) StartJob(job *Job) (*Response, error) {
	u := fmt.Sprintf("jobs/%s/start", job.ID)

	req, err := c.newRequest("PUT", u, &jobStartRequest{
		StartedAt: job.StartedAt,
	})
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, nil)
}

// FinishJob finishes the passed in job
func (c *Client) FinishJob(job *Job) (*Response, error) {
	u := fmt.Sprintf("jobs/%s/finish", job.ID)

	req, err := c.newRequest("PUT", u, &jobFinishRequest{
		FinishedAt:        job.FinishedAt,
		ExitStatus:        job.ExitStatus,
		Signal:            job.Signal,
		SignalReason:      job.SignalReason,
		ChunksFailedCount: job.ChunksFailedCount,
	})
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, nil)
}
