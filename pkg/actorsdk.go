package actorsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type actorSDK struct {
	endpoint     string
	orgID        string
	deploymentID string
}

func newActorSDK(endpoint, orgID, deploymentID string) (*actorSDK, error) {
	sdk := &actorSDK{
		endpoint:     endpoint,
		orgID:        orgID,
		deploymentID: deploymentID,
	}

	return sdk, nil
}

func (sdk *actorSDK) FindActors(fsmID string, filters []FindFilter) ([]Actor, error) {
	var response struct {
		Actors []Actor `json:"actors"`
	}
	err := sdk.sendRequest("find_actors", map[string]interface{}{
		"fsm_id":  fsmID,
		"filters": filters,
	}, &response)
	if err != nil {
		panic(err)
	}
	return response.Actors, nil
}

func (sdk *actorSDK) GetActorByRef(fsmID, ref string) (Actor, error) {
	var response struct {
		Actor Actor  `json:"actor"`
		Error string `json:"error"`
	}
	err := sdk.sendRequest("get_actor", map[string]interface{}{
		"fsm_id":    fsmID,
		"actor_ref": ref,
	}, &response)
	if err != nil {
		return response.Actor, err
	}
	if response.Error != "" {
		return response.Actor, newErrorFromSDKResponse(response.Error)
	}

	return response.Actor, nil
}

func (sdk *actorSDK) Resource(requestType ResourceRequest, branch, repo string) error {
	var response struct {
		Error *string `json:"error"`
	}
	err := sdk.sendRequest(string(requestType), map[string]interface{}{
		"repo":   repo,
		"branch": branch,
	}, &response)
	if err != nil {
		log.Println("sdk error")
		return err
	}
	if response.Error != nil {
		log.Println("error from actord")
		return errors.New(*response.Error)
	}
	return nil
}

func (sdk *actorSDK) sendRequest(method string, data map[string]interface{}, response interface{}) error {
	url := fmt.Sprintf("%s/%s", sdk.endpoint, method)

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("ACTORD_ORG_ID", sdk.orgID)
	req.Header.Add("ACTORD_DEPLOYMENT_ID", sdk.deploymentID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("bad response(statusCode: %d): %s", resp.StatusCode, string(responseBody))
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}
	return nil
}
