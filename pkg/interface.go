package actorsdk

type ResourceRequest string

const (
	ResourceCheckout ResourceRequest = "resource/checkout"
	ResourcePull     ResourceRequest = "resource/pull"
	ResourcePush     ResourceRequest = "resource/push"
)

type ActordSDK interface {
	FindActors(fsmID string, filters []FindFilter) ([]Actor, error)
	GetActorByRef(fsmID, ref string) (Actor, error)
	Resource(requestType ResourceRequest, branch, repo string) error
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
