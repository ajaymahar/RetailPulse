package datastore

import (
	"github.com/ajaymahar/RetailPulse/internal"
)

// JobRepository is a port for repository to abstract the implementation
type JobRepository interface {
	Save(internal.Job) (internal.Job, error)
	Find(int) (*JobStatus, error)
}
