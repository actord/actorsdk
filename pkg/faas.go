package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
)

type ExecuteFunc func(inputDecoder *json.Decoder) (interface{}, error)

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
}

func NewFaaS(revision string, executeFunc ExecuteFunc) (*FaaS, error) {
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

func (f *FaaS) Listen() error {
	http.HandleFunc("/execute", f.handleExecute)
	http.HandleFunc("/healthcheck", f.handleHealthcheck)

	log.Printf("Listen and serve at :%s, revision: %s", f.Port, f.revision)
	fmt.Println("FAAS READY.")
	return http.ListenAndServe(fmt.Sprintf(":%s", f.Port), nil)
}

func (f *FaaS) handleExecute(w http.ResponseWriter, r *http.Request) {

	inputDecoder := json.NewDecoder(r.Body)

	output, err := f.executeFunc(inputDecoder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}

func (f *FaaS) handleHealthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ready"))
}
