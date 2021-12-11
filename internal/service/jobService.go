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
	GetStatus(int) (*internal.Job, error)
	GetAllJobs() ([]internal.Job, error)
}

// CreateJob is method of DefaultJobService to pass the request to the repo
func (dSvc DefaultJobService) CreateJob(job internal.Job) (internal.Job, error) {
	job, err := dSvc.repo.Save(job)
	if err != nil {
		return internal.Job{}, fmt.Errorf("repo save: %w", err)
	}
	return job, nil
}

func (dSvc DefaultJobService) GetAllJobs() ([]internal.Job, error) {
	jobs, err := dSvc.repo.FindAll()
	if err != nil {
		return []internal.Job{}, fmt.Errorf("repo findall: %w", err)
	}
	return jobs, nil
}

func (dSvc DefaultJobService) GetStatus(jobID int) (*internal.Job, error) {
	job, err := dSvc.repo.Find(jobID)
	if err != nil {
		return &internal.Job{}, fmt.Errorf("repo find: %w", err)
	}
	return job, nil
}
