package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"html/template"
)

type Project struct {
	Name        string
	Repository  string
	Twitter     string
	Website     string
	Description string
}

type Projects struct {
	Items []Project
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	// 解析模板
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Println("Parse template failed, err%v", err)
		return
	}
	// 渲染模板
	var projects []Project
	client := &http.Client{Timeout: 10 * time.Second}
	request, _ := http.NewRequest("GET", getEnv("LIST_PROJECTS_ENDPOINT", "http://127.0.0.1:8001/projects"), nil)
	request.Header.Set("Content-type", "application/json")
	resp, err := client.Do(request)
	if err == nil && resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&projects)
	}

	ps := Projects{
		Items: projects,
	}

	err = t.Execute(w, ps)
	if err != nil {
		log.Println("render template failed, err%v", err)
		return
	}
}

var log = logrus.New()

func main() {
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	router := mux.NewRouter()
	router.HandleFunc("/", sayHello).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./templates/"))))

	var handler http.Handler = router
	handler = &logHandler{log: log, next: handler} // add logging

	log.Infof("starting frontend on 0.0.0.0:8000")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", handler))
}
