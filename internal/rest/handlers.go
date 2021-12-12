package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ajaymahar/RetailPulse/internal"
	"github.com/ajaymahar/RetailPulse/internal/service"
	"github.com/gorilla/mux"
)

type JobHandler struct {
	svc service.JobService
}

// Factory method
func NewJobHandler(svc service.JobService) *JobHandler {
	return &JobHandler{
		svc: svc,
	}
}

// CreateJobResponse
type CreateJobResponse struct {
	// internal.Job `json:"jobs"`
	JobID int `json:"job_id"`
}

// GetJobResponse
type GetJobResponse struct {
	Status string `json:"status,omitempty"`
	JobID  int    `json:"job_id"`
}

// FailedJobResponse
type FailedJobResponse struct {
	GetJobResponse
	Error error `json:"error"`
}

// EmptyResponse
type EmptyResponse struct{}

func (jh *JobHandler) Register(r *mux.Router) {

	subR := r.PathPrefix("/api").Subrouter()
	subR.HandleFunc("/submit", jh.ValidatePayload(jh.createJob)).Methods("POST")

	subR.HandleFunc("/status", jh.getJobStatus).Queries("job_id", "{job_id}").Methods("GET")

}

//NOTE: assumed that each job will have different id
// not checking for dubplicate jobs if user submits the same jobs again
// will accept the requst and process
func (jh *JobHandler) createJob(rw http.ResponseWriter, r *http.Request) {
	//TODO: refactor/clean up needed
	var j internal.Job

	// deserialize the jobs
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		log.Println(err)
		renderErrorResponse(rw, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	job, err := jh.svc.CreateJob(j)
	if err != nil {
		log.Println(err)
		renderErrorResponse(rw, err.Error(), http.StatusBadRequest)
		return

	}
	renderResponse(rw, &CreateJobResponse{
		JobID: job.JobID,
	}, http.StatusCreated)

}

// Validate payload data in this middle ware
// NOTE: this can be done in more effecient way via validator library
// gopkg.in/go-playground/validator.v9
func (jh *JobHandler) ValidatePayload(handler http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		var err error

		defer func() {
			if err != nil {
				renderErrorResponse(rw, err.Error(), http.StatusBadRequest)
			}
		}()
		var j internal.Job

		// deserialize the jobs
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err = errors.New("invalid json data")
			return
		}

		if !json.Valid(body) {
			err = errors.New("invalid json")
		}
		if err = json.Unmarshal(body, &j); err != nil {
			err = errors.New("invalid json data")
			return
		}

		// NOTE: putting back the body data for sub-sequent call
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if err != nil {
			log.Println(err)
			return
		}

		// Validate fields
		if j.Count != len(j.Visits) || j.Count <= 0 || len(j.Visits) == 0 {
			err = fmt.Errorf("invalid payload count: <%v> and visits: <%v> data missmatch", j.Count, len(j.Visits))
			return
		}

		// validate visit fields
		for _, f := range j.Visits {
			if f.StoreID == "" || f.VisitTime == "" || len(f.ImageURL) == 0 {
				err = errors.New("invalid payload data please check store_id, visit_time, or image_url fields")
				return
			}
		}
		handler.ServeHTTP(rw, r)
	}
}

func (jh *JobHandler) getJobStatus(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement me
	jobID, err := strconv.Atoi(r.URL.Query()["job_id"][0])
	if err != nil {
		log.Println(err)
		renderErrorResponse(rw, "job_id must be <int> type", http.StatusBadRequest)
		return
	}

	job, err := jh.svc.GetStatus(jobID)
	if err != nil {
		log.Println(err)
		// NOTE:
		// status code should be status not found
		// used StatusBadRequest as per given instructions in assignment
		// renderErrorResponse(rw, err.Error(), http.StatusBadRequest)
		renderResponse(rw, EmptyResponse{}, http.StatusBadRequest)

		return
	}
	if job.Error != nil {
		log.Println(job.Error.Error())

		renderResponse(rw, &internal.JobError{
			Status: job.Status,
			JobID:  job.JobID,
			Err: []internal.StoreError{
				{
					StoreID: job.StoreID,
					// SErr:    job.Error.Error(),
					SErr: "invalid image url ",
				},
			},
		}, http.StatusOK)
		return
	}

	renderResponse(rw, &GetJobResponse{
		Status: job.Status,
		JobID:  job.JobID,
	}, http.StatusOK)
}
