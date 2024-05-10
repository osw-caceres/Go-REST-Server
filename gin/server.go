package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/internal/taskstore"
	"github.com/gin-gonic/gin"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

// func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling task create at %s\n", req.URL.Path)

// 	// Types used internally in this handler to (de-)serialize the request and
// 	// response from/to JSON.
// 	type RequestTask struct {
// 		Text string    `json:"text"`
// 		Tags []string  `json:"tags"`
// 		Due  time.Time `json:"due"`
// 	}

// 	type ResponseId struct {
// 		Id int `json:"id"`
// 	}

// 	// Enforce a JSON Content-Type.
// 	contentType := req.Header.Get("Content-Type")
// 	mediatype, _, err := mime.ParseMediaType(contentType)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 	}
// 	if mediatype != "application/json" {
// 		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
// 		return
// 	}

// 	dec := json.NewDecoder(req.Body)
// 	dec.DisallowUnknownFields()
// 	var rt RequestTask
// 	if err := dec.Decode(&rt); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
// 	js, err := json.Marshal(ResponseId{Id: id})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

func (ts *taskServer) createTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindBodyWithJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

// func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling get all tasks at %s\n", req.URL.Path)

// 	allTasks := ts.store.GetAllTasks()
// 	js, err := json.Marshal(allTasks)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

func (ts *taskServer) getAllTasksHandler(c *gin.Context) {
	allTask := ts.store.GetAllTasks()
	c.JSON(http.StatusOK, allTask)
}

// func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling get task at %s\n", req.URL.Path)

// 	id, err := strconv.Atoi(req.PathValue("id"))
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusBadRequest)
// 		return
// 	}

// 	task, err := ts.store.GetTask(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	js, err := json.Marshal(task)
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusBadRequest)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(js)
// }

func (ts *taskServer) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
	}

	c.JSON(http.StatusOK, task)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *taskServer) DeleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all tasks at %s\n", req.URL.Path)

	ts.store.DeleteAllTasks()
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by tag at %s\n", req.URL.Path)

	tag := req.PathValue("tag")

	tasks := ts.store.GetTasksByTag(tag)
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by due date at %s\n", req.URL.Path)

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	year, errYear := strconv.Atoi(req.PathValue("year"))
	month, errMonth := strconv.Atoi(req.PathValue("month"))
	day, errDay := strconv.Atoi(req.PathValue("day"))

	if errYear != nil || errMonth != nil || errDay != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	js, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// func ping(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling ping at %s\n", req.URL.Path)
// }

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func main() {
	// mux := http.NewServeMux()
	// server := NewTaskServer()

	// mux.HandleFunc("GET /ping/", ping)
	// mux.HandleFunc("POST /task/", server.createTaskHandler)
	// mux.HandleFunc("GET /task/", server.getAllTasksHandler)
	// mux.HandleFunc("DELETE /task/", server.deleteTaskHandler)
	// mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	// mux.HandleFunc("DELETE /task/{id}/", server.deleteTaskHandler)
	// mux.HandleFunc("GET /tag/{tag}/", server.tagHandler)
	// mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)

	// log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))

	router := gin.Default()
	server := NewTaskServer()

	router.GET("/task/", server.getAllTasksHandler)
	router.GET("/task/:id", server.getTaskHandler)
	router.GET("/ping", ping)

}
