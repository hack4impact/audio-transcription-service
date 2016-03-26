package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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
}

var routes = []route{
	route{
		"hello",
		"GET",
		"/hello/{name}",
		helloHandler,
	},
	route{
		"add_job",
		"POST",
		"/add_job",
		initiateTranscriptionJobHandler,
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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

// initiateTranscriptionJobHandlerJSON takes a POST request containing a json object,
// decodes it into a transcriptionJobData struct, and starts a transcription task.
// func initiateTranscriptionJobHandlerJSON(w http.ResponseWriter, r *http.Request) {
// 	var jsonData transcriptionJobData
//
// 	// unmarshal from the response body directly into our struct
// 	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	executer := tasks.DefaultTaskExecuter
// 	id := executer.QueueTask(transcription.MakeTaskFunction(jsonData.AudioURL, jsonData.EmailAddresses))
//
// 	log.Print(w, "Accepted task %d!", id)
// }

// initiateTranscriptionJobHandler takes a POST request from a form,
// decodes it into a transcriptionJobData struct, and starts a transcription task.
func initiateTranscriptionJobHandler(w http.ResponseWriter, r *http.Request) {
	formData := transcriptionJobData{
		AudioURL:       r.FormValue("AudioURL"),
		EmailAddresses: r.Form["EmailAddresses"],
	}
	executer := tasks.DefaultTaskExecuter
	id := executer.QueueTask(transcription.MakeTaskFunction(formData.AudioURL, formData.EmailAddresses))

	log.Print(w, "Accepted task %d!", id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

// jobStatusHandler returns the status of a task with given id.
func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]

	executer := tasks.DefaultTaskExecuter
	status := executer.GetTaskStatus(id)
	w.Write([]byte(status.String()))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	const tpl = `
	<!DOCTYPE html>
	<html>
		<head>
		  <title></title>
		</head>
	<body>
	  <form action="/add_job" method="POST">
	    <div>URL:<input type="url" name="AudioURL" required></div>
			<div>Email Addresses:<input type="email" name="EmailAddresses" multiple required></div>
			<div><input type="submit" value="Submit"></div>
	  </form>
	</form>
	</body>
	</html>
	`
	t, _ := template.New("webpage").Parse(tpl)
	_ = t.Execute(w, transcriptionJobData{})
}
