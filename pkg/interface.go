package actorsdk

import (
	"io"
	"net/http"
)

type ResourceRequest string

const (
	ResourceCheckout ResourceRequest = "resource/checkout"
	ResourcePull     ResourceRequest = "resource/pull"
	ResourcePush     ResourceRequest = "resource/push"
)

type ActordSDK interface {
	FindActors(appID, fsmPath string, filters []FindFilter) ([]Actor, error)
	GetActorByRef(appID, fsmPath, ref string) (Actor, error)
	Resource(requestType ResourceRequest, branch, repo string) error
	ResourceRead(branch, repo, path string) (io.ReadCloser, error)
	ResourceWriteString(branch, repo, path, body string) error
	ResourceDirList(branch, repo, path string) ([]DirItem, error)

	// internal things
	SendRequestResp(method string, data map[string]interface{}) (*http.Response, error)
	SendRequest(method string, data map[string]interface{}, response interface{}) error
}

func NewActorSDK(endpoint, orgID, deploymentID string) (ActordSDK, error) {
	return newActorSDK(endpoint, orgID, deploymentID)
}

type FindFilter struct {
	Index string `json:"index"`
	Fun   string `json:"fun"`
	Value string `json:"value"`
}

type Actor struct {
	ID                string                 `json:"id"`
	Ref               string                 `json:"ref"`
	CurrentNodeID     string                 `json:"current_node_id"`
	CurrentLogicID    string                 `json:"current_logic_id"`
	AwaitEvent        bool                   `json:"await_event"`
	AllowedEventTypes []LogicAwaitEventType  `json:"allowed_event_types,omitempty"`
	CreatedAt         int64                  `json:"created_at"`
	UpdatedAt         int64                  `json:"updated_at"`
	Data              map[string]interface{} `json:"data"`
}

type LogicAwaitEventType struct {
	Typename        string  `json:"typename"`
	Label           *string `json:"label"`
	InitialFillFrom *string `json:"initial_fill_from"`
}

type LogicAwaitEvent struct {
	ID                string                `json:"id"`
	Logic             string                `json:"logic"`
	AllowedEventTypes []LogicAwaitEventType `json:"allowed_event_types"`
}

type DirItem struct {
	IsDir    bool   `json:"isDir"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Mode     uint32 `json:"mode"`
	ModeTime int64  `json:"modeTime"`
}
