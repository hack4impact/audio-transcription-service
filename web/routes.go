package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/hack4impact/transcribe4all/tasks"
	"github.com/hack4impact/transcribe4all/transcription"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type transcriptionJobData struct {
	AudioURL       string   `json:"audioURL"`
	EmailAddresses []string `json:"emailAddresses"`
	SearchWords    []string `json:"searchWords"`
}

type flash struct {
	Title string
	Body  string
}

var routes = []route{
	route{
		"add_job",
		"POST",
		"/add_job",
		initiateTranscriptionJobHandler,
	},
	route{
		"add_job_json",
		"POST",
		"/add_job_json",
		initiateTranscriptionJobHandlerJSON,
	},
	route{
		"health",
		"GET",
		"/health",
		healthHandler,
	},
	route{
		"job_status",
		"GET",
		"/job_status/{id}",
		jobStatusHandler,
	},
	route{
		"form",
		"GET",
		"/",
		formHandler,
	},
}

// TODO: Flashes is a complete hack. Use real sessions.
var flashes = []flash{}

// initiateTranscriptionJobHandlerJSON takes a POST request containing a json object,
// decodes it into a transcriptionJobData struct, and starts a transcription task.
func initiateTranscriptionJobHandlerJSON(w http.ResponseWriter, r *http.Request) {
	jsonData := new(transcriptionJobData)

	// unmarshal from the response body directly into our struct
	if err := json.NewDecoder(r.Body).Decode(jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	executer := tasks.DefaultTaskExecuter
	executer.QueueTask(transcription.MakeIBMTaskFunction(jsonData.AudioURL, jsonData.EmailAddresses, jsonData.SearchWords))
}

// initiateTranscriptionJobHandler takes a POST request from a form,
// decodes it into a transcriptionJobData struct, and starts a transcription task.
func initiateTranscriptionJobHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r.FormValue("url"))
	executer := tasks.DefaultTaskExecuter
	id := executer.QueueTask(transcription.MakeIBMTaskFunction(r.FormValue("url"), r.Form["emails"], r.Form["words"]))

	flashes = append(flashes, flash{
		Title: "Task Started!",
		Body:  fmt.Sprintf("Task %s was successfully started. The results will be emailed to you upon completion.", id),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

// healthHandler returns a 200 response to the client if the server is healthy.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK :)")
}

// jobStatusHandler returns the status of a task with given id.
func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]

	executer := tasks.DefaultTaskExecuter
	status := executer.GetTaskStatus(id)
	io.WriteString(w, status.String())
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/form.html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, flashes)
	if err != nil {
		log.Fatal(err)
	}
	flashes = []flash{}
}
