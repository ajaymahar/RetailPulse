package internal

import "fmt"

// define custom error to response back to the client
// if job status is failed
type JobError struct {
	Status string       `json:"status"`
	JobID  int          `json:"job_id"`
	Err    []StoreError `json:"error,omitempty"`
}

type StoreError struct {
	StoreID string `json:"store_id"`
	SErr    string `json:"error"`
}

// implement the Error interface
func (je *JobError) Error() string {
	return fmt.Sprintf("%v \n %v", je.Err[0].StoreID, je.Err[0].SErr)
}
