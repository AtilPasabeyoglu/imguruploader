package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io/ioutil"
	"sync"
	//Gorilla mux library for routing
	//For more info please visit "http://www.gorillatoolkit.org/pkg/mux"
	"github.com/gorilla/mux"
	//Google library for creating uuid
	"github.com/google/uuid"
)

//This handler functions uploads the given request to imgur via running uploadTask function async.
func CreateJob(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	//jobID is given
	jobID := uuid.New().String();
	//Job Creation time
	createdTimestamp := time.Now().UTC().Format(time.RFC3339)

	//This is for checking the jobID for requesting
	fmt.Println(jobID)

	//Appending the job details
	job=jobDetails{jobID:jobID,imgDetail:imageDetails,CreatedTimeStamp:createdTimestamp,FinishedTimeStamp:"null"}


	type urlList struct {

		Urls []string `json:"urls"`

	}

	var urlsJson urlList
	body, err := ioutil.ReadAll(r.Body)
	if isError(err) { }

	error := json.Unmarshal([]byte(body), &urlsJson)
	if error != nil {
		fmt.Printf("err was %v", error)
	}

	jobStatus.pending=urlsJson.Urls
	wg:=&sync.WaitGroup{}
	for i := range urlsJson.Urls{
		wg.Add(1)
		go uploadTask(urlsJson.Urls[i],wg)
	}

	wg.Wait()

	finishedTimestamp := time.Now().UTC().Format(time.RFC3339)
	job.FinishedTimeStamp=finishedTimestamp


	result,_ := json.Marshal(jobID)

	w.Write(result)


}


//This function returns the Job Status with completed, pending, failure upload links.
//Also shows if the job is finished or still in progress.
func GetTasks(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	vars:=mux.Vars(r)
	jobID:=vars["jobID"]

	finished := false
	status := "pending"
	_=finished
	_=status

	if len(jobStatus.pending)==0{
		finished=true
		status="finished"
	}else {
		status = "in-progress"
	}

	//This struct contains information about links given by the request if they are completed or not.
	type uploaded struct {

		Pending []string `json:"pending"`
		Completed []string `json:"completed"`
		Failure []string `json:"failure"`

	}

	//Variable created from uploaded struct to keep track
	uploadStatus:=uploaded{

		Pending:jobStatus.pending,
		Completed:jobStatus.completed,
		Failure:jobStatus.failure,

	}

	//This is the struct for final job status.
	type resultStr struct {

		ID string `json:"jobID"`
		Created string `json:"created"`
		Finished string `json:"finished"`
		Status string `json:"status"`
		uploaded

	}
	//Variable created from resultStr struct
	var result resultStr

	//Checks if the ID is correct. If not prints out a warning message.
	if job.jobID==jobID{
		result = resultStr{jobID,job.CreatedTimeStamp,job.FinishedTimeStamp,status,uploadStatus }
	}else {
		fmt.Println("No jobs with that ID!")
	}

	//Result is encoded to json
	resultJson,_ := json.Marshal(result)

	w.Write(resultJson)


}

//This handler function returns the imgur links for the uploaded URLs given by the request
func GetLinks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	uploadedLinks := make([]string, 0)

	for _,url := range imageDetails{
		if url.Status=="Success!"{
			uploadedLinks = append(uploadedLinks, url.ImgurURL)
		}
	}

	result,_ := json.Marshal(uploadedLinks)

	w.Write(result)

}