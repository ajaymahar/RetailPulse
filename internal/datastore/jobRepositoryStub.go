package datastore

import (
	"errors"

	"github.com/ajaymahar/RetailPulse/internal"
)

// JobRepositoryStub is basic implementaiton to store the job data
type JobRepositoryStub struct {
	job []internal.Job
}

func NewJobRepositoryStub() *JobRepositoryStub {
	return &JobRepositoryStub{}
}

func (jr *JobRepositoryStub) Save(job internal.Job) (internal.Job, error) {
	// Get Previous jobid, this can be improved later on
	// TODO: find a batter way to get next id
	// it's just a dummy implementation
	var pID int
	if l := len(jr.job); l != 0 {

		pID = jr.job[len(jr.job)-1].JobID
	}
	job.JobID = pID + 1
	job.Status = internal.Ongoing
	jr.job = append(jr.job, job)
	return job, nil
}

func (jr *JobRepositoryStub) FindAll() ([]internal.Job, error) {
	return jr.job, nil
}

func (jr *JobRepositoryStub) Find(jobID int) (*internal.Job, error) {
	for _, j := range jr.job {
		if j.JobID == jobID {
			return &j, nil
		}
	}
	return nil, errors.New("job not found")
}
