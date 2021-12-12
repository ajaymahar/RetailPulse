package datastore

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ajaymahar/RetailPulse/internal"
)

// local queue to store the jobs
var queue *jobQueue

// store the job status in memory
var jobStat map[int]JobStatus

var storeIDs map[string]bool

func init() {
	// to avoid nil poiner error at run time
	queue = NewjobQueue()

	// list of valid storeID
	// assumed we have the valid storeID which we can validate the requst coming from
	storeIDs = map[string]bool{
		"S00339218": true,
		"S01408764": true,
		"S00339219": true,
	}
	jobStat = make(map[int]JobStatus, 100)
	// launch workers
	launchWorker()
}

// jobQueue is local jobQueue to sotre jobs to process
type jobQueue struct {
	q chan internal.Job
}

// Factory method to create new queue
func NewjobQueue() *jobQueue {
	return &jobQueue{
		q: make(chan internal.Job, 100),
	}
}

// JobStatus will hold the info about all the jobs ongoin, failed, completed or error
type JobStatus struct {
	JobID     int
	StoreID   string
	Status    string
	Error     error
	Perimeter int
}

// NewJobStatus is factory function to initialize NewJobStatus
func NewJobStatus() *JobStatus {
	return &JobStatus{
		JobID:     0,
		StoreID:   "",
		Status:    "",
		Error:     nil,
		Perimeter: 0,
	}
}

// JobRepositoryStub is basic implementaiton to store the job data
type JobRepositoryStub struct {
	job   []internal.Job
	local *jobQueue
}

func NewJobRepositoryStub() *JobRepositoryStub {
	return &JobRepositoryStub{
		job:   []internal.Job{},
		local: queue,
	}
}
func (jr *JobRepositoryStub) Save(job internal.Job) (internal.Job, error) {

	if _, ok := storeIDs[job.Visits[0].StoreID]; !ok {
		return internal.Job{}, errors.New("store id not exist")
	}

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

	// adding job to the local queue
	jr.local.q <- job

	// update the status of job
	jobStat[job.JobID] = JobStatus{
		JobID:     job.JobID,
		StoreID:   job.Visits[0].StoreID,
		Status:    "ongoing",
		Error:     nil,
		Perimeter: 0,
	}
	return job, nil
}

func (jr *JobRepositoryStub) Find(jobID int) (*JobStatus, error) {
	js, ok := jobStat[jobID]
	if !ok {
		return nil, errors.New("job not found")
	}
	return &js, nil
}

// Save the status of each job if completed or failed
func saveStatus(jobID int, storeID string, per int, status string, err error) {

	// check if map is nil
	if jobStat == nil {
		jobStat = map[int]JobStatus{}
	}
	// save job status in memory
	jobStat[jobID] = JobStatus{
		JobID:     jobID,
		StoreID:   storeID,
		Status:    status,
		Error:     err,
		Perimeter: per,
	}
}

func launchWorker() {
	// launch 3 workers to work on jobs
	// each worker is required to perform the following jobs
	// 1. Get the jobs from the local queue
	// 2. get all the URLs from the given job
	// 3. download the images
	// 4. process the image
	// 5. store/record the status of job [ongoing, completed, failed] or error if any
	for i := 0; i < 3; i++ {
		go func() {
			var j internal.Job
			var perimeter int
			for j = range queue.q {
				fmt.Printf("working on job_id: %v", j.JobID)
				urls, err := getURLs(j)
				if err != nil {
					// if there is an error stop here and save the status and return
					saveStatus(j.JobID, j.Visits[0].StoreID, 0, "failed", fmt.Errorf("launchWorker: getURLs: %w", err))
					return
				}

				for _, url := range urls {

					if err := validateURL(url); err != nil {
						saveStatus(j.JobID, j.Visits[0].StoreID, 0, "failed", fmt.Errorf("launchWorker: validateURL: %w", err))
						return
					}
					// TODO: handle the error if there is any
					name, err := internal.GetFileName(url)
					if err != nil {
						// if there is an error stop here and save the status and return
						saveStatus(j.JobID, j.Visits[0].StoreID, 0, "failed", fmt.Errorf("launchWorker: GetFileName: %w", err))
						return
					}

					// TODO: handle the error if there is any
					err = internal.DownloadImage(url, name)
					if err != nil {
						// if there is an error stop here and save the status and return
						saveStatus(j.JobID, j.Visits[0].StoreID, 0, "failed", fmt.Errorf("launchWorker: DownloadImage: %w", err))
						return
					}

					// process the image perimeter 2* [Height+Width] of each image
					perimeter, err = processImage(name)
					if err != nil {
						// if there is an error stop here and save the status and return
						saveStatus(j.JobID, j.Visits[0].StoreID, 0, "failed", fmt.Errorf("launchWorker: processImage: %w", err))
						return
					}
				}
				// if everyting is ok mark the status completed
				saveStatus(j.JobID, j.Visits[0].StoreID, perimeter, "completed", nil)
			}

		}()
	}
}

// get all the urls from the given jobs
func getURLs(job internal.Job) ([]string, error) {
	if len(job.Visits) == 0 {
		return nil, fmt.Errorf("no urls in the job, job_id: %v", job.JobID)
	}
	urls := make([]string, 0, len(job.Visits))

	for _, s := range job.Visits {
		urls = append(urls, s.ImageURL...)
	}
	return urls, nil
}

// process the image
func processImage(imgName string) (int, error) {
	return internal.GetDimmensions(imgName)
}

// validate URL
func validateURL(u string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return fmt.Errorf("invalid url %w", err)
	}
	return nil
}
