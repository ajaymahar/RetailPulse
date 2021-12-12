package service

import (
	"fmt"

	"github.com/ajaymahar/RetailPulse/internal"
	"github.com/ajaymahar/RetailPulse/internal/datastore"
)

// DefaultJobService is concrete implementaion of service
type DefaultJobService struct {
	repo datastore.JobRepository
}

// NewDefaultJobService is factory method to create new NewDefaultJobService
func NewDefaultJobService(repo datastore.JobRepository) DefaultJobService {
	return DefaultJobService{repo}
}

// JobService is service port to intract with external entities like http handlers
type JobService interface {
	CreateJob(internal.Job) (internal.Job, error)
	GetStatus(int) (*datastore.JobStatus, error)
}

// CreateJob is method of DefaultJobService to pass the request to the repo
func (dSvc DefaultJobService) CreateJob(job internal.Job) (internal.Job, error) {
	job, err := dSvc.repo.Save(job)
	if err != nil {
		return internal.Job{}, fmt.Errorf("repo save: %w", err)
	}
	return job, nil
}

// GetStatus to get the job status
func (dSvc DefaultJobService) GetStatus(jobID int) (*datastore.JobStatus, error) {
	job, err := dSvc.repo.Find(jobID)
	if err != nil {
		return &datastore.JobStatus{
			JobID:     jobID,
			Status:    "",
			Error:     err,
			Perimeter: 0,
		}, fmt.Errorf("repo find: %w", err)
	}
	return job, nil
}
