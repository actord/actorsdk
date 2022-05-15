package actorsdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

type ExecuteFunc func(r *http.Request, inputDecoder *json.Decoder) (interface{}, error)

type FaaS struct {
	Port           string
	ActordEndpoint string
	OrgID          string
	DeploymentID   string

	sdk ActordSDK

	inputT  reflect.Type
	outputT reflect.Type

	revision string

	executeFunc ExecuteFunc
	endpoints   map[string]http.HandlerFunc
}

type FaaSResponse struct {
	Response interface{} `json:"response"`
	Error    string      `json:"error,omitempty"`
}

func NewFaaS(revision string, executeFunc ExecuteFunc, endpoints map[string]http.HandlerFunc) (*FaaS, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	actordEndpoint := os.Getenv("ACTORD_ENDPOINT")
	orgID := os.Getenv("ACTORD_ORG_ID")
	deploymentID := os.Getenv("ACTORD_DEPLOYMENT_ID")

	f := &FaaS{
		Port:           port,
		ActordEndpoint: actordEndpoint,
		OrgID:          orgID,
		DeploymentID:   deploymentID,
		revision:       revision,
		executeFunc:    executeFunc,
		endpoints:      endpoints,
	}

	return f, nil
}

func (f *FaaS) SDK() ActordSDK {
	if f.sdk == nil {
		sdk, err := NewActorSDK(f.ActordEndpoint, f.OrgID, f.DeploymentID)
		if err != nil {
			panic(err)
		}
		f.sdk = sdk
	}
	return f.sdk
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h(rw, r) // serve the original request

		duration := time.Since(start)

		log.Printf("%s %s duration=%dmsec\n", method, uri, duration.Milliseconds())
	}
	return logFn
}

func (f *FaaS) Listen() error {
	http.HandleFunc("/execute", WithLogging(f.handleExecute))
	http.HandleFunc("/healthcheck", WithLogging(f.handleHealthcheck))

	if f.endpoints != nil {
		for endpoint, handlerFunc := range f.endpoints {
			http.HandleFunc(endpoint, WithLogging(handlerFunc))
		}
	}

	log.Printf("Listen and serve at :%s, revision: %s", f.Port, f.revision)
	fmt.Println("FAAS READY.")
	return http.ListenAndServe(fmt.Sprintf(":%s", f.Port), nil)
}

func (f *FaaS) handleExecute(w http.ResponseWriter, r *http.Request) {

	inputDecoder := json.NewDecoder(r.Body)

	output, err := f.executeFunc(r, inputDecoder)
	response := FaaSResponse{
		Response: output,
	}
	if err != nil {
		statusErr, ok := err.(StatusError)
		if !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Error = statusErr.Error()
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (f *FaaS) handleHealthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ready"))
}
