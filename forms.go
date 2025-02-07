package forms

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
)

type cloudIdRes struct {
    CloudId string `json:"cloudId"`
}

type formIdRes struct {
    FormId string `json:"id"`
    FormTemplate []byte `json:"formTemplate"`
    Internal bool `json:"internal"`
    Submitted bool `json:"submitted"`
    Lock bool `json:"lock"`
    Name string `json:"name"`
    Updated string `json:"updated"`
}

type FormService struct {
    username string
    jira_api_token string
    jira_base_url string
    jira_api_base_url string
    http_client *http.Client
    cloud_id string
}

// Uses the fields of the struct to make a request to the _edge/tenant_info endpoint to retrieve the cloudId and set the corresponding field to the
// value
func (formservice *FormService) setCloudId() {

    resstring := formservice.request("_edge/tenant_info", false)

    var cir cloudIdRes

    json.Unmarshal(resstring, &cir)
    formservice.cloud_id = cir.CloudId

}

func (formService FormService) getIssueFormId(issueKey string) string {

    resstring := formService.request("forms/cloud/" + formService.cloud_id + "/issue/" + issueKey + "/form", true)

    var fir formIdRes
    json.Unmarshal(resstring[1:len(resstring) - 1], &fir)

    return fir.FormId

}

func (formService FormService) getIssueForm(issueKey string, formId string) {

    resstring := formService.request("forms/cloud/" + formService.cloud_id + "/issue/" + issueKey + "/form/" + formId, true)

    fmt.Println(string(resstring))

}

func (formservice FormService) request(uri string, api_request bool) []byte {

    base_url := ""

    if api_request {
        base_url = formservice.jira_api_base_url
    } else {
        base_url = formservice.jira_base_url
    }

    req, _ := http.NewRequest("GET", base_url + uri, nil)
    req.SetBasicAuth(formservice.username, formservice.jira_api_token)

    res, _ := http.DefaultClient.Do(req)
    resstring, _ := io.ReadAll(res.Body)

    return resstring

}
