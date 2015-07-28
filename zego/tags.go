package zego

import (
	"encoding/json"
	"fmt"
)

type TagArray struct {
	Tags []string `json:"tags"`
}

func (a Auth) ShowTicketTags(ticket_id string) (*Resource, error) {

	path := "/tickets/" + ticket_id + "/tags.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return resource, nil

}

func (a Auth) ShowTopicTags(topic_id string) (*Resource, error) {

	path := "/topics/" + topic_id + "/tags.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return resource, nil

}

func (a Auth) ShowOrganizationTags(organization_id string) (*Resource, error) {

	path := "/organizations/" + organization_id + "/tags.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return resource, nil

}

func (a Auth) setTags(kind string, id uint32, tags []string) (*[]string, error) {
	path := fmt.Sprintf("/%s/%d/tags.json", kind, id)
	bytes, err := json.Marshal(TagArray{tags})
	if err != nil {
		return nil, err
	}
	resource, err := api(a, "POST", path, string(bytes))
	if err != nil {
		return nil, err
	}

	tagsResponse := &TagArray{}
	json.Unmarshal([]byte(resource.Raw), tagsResponse)

	return &tagsResponse.Tags, nil
}

func (a Auth) SetUserTags(user_id uint32, tags []string) (*[]string, error) {
	return a.setTags("users", user_id, tags)
}

func (a Auth) SetOrganizationTags(org_id uint32, tags []string) (*[]string, error) {
	return a.setTags("organizations", org_id, tags)
}
