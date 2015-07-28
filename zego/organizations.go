package zego

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type OrganizationArray struct {
	Organizations []Organization `json:"organizations"`
}

type Organization struct {
	Id                 uint32                 `json:"id,omitempty"`
	ExternalId         string                 `json:"external_id,omitempty"`
	Url                string                 `json:"url,omitempty"`
	Name               string                 `json:"name,omitempty"`
	CreatedAt          string                 `json:"created_at,omitempty"`
	UpdatedAt          string                 `json:"updated_at,omitempty"`
	DomainNames        []string               `json:"domain_names,omitempty"`
	Details            string                 `json:"details,omitempty"`
	Notes              string                 `json:"notes,omitempty"`
	GroupId            int                    `json:"group_id,omitempty"`
	SharedTickets      bool                   `json:"shared_tickets,omitempty"`
	SharedComments     bool                   `json:"shared_comments,omitempty"`
	Tags               []string               `json:"tags,omitempty"`
	OrganizationFields []*OrganizationalField `json:"organization_fields,omitempty"`
}

type OrganizationalField struct {
	OrgDropdown string  `json:"org_dropdown"`
	OrgDecimal  float32 `json:"org_decimal"`
}

type SingleOrganization struct {
	Organization *Organization `json:"organization"`
}

func parseOrganizations(resource *Resource) (*[]Organization, error) {
	organizations := &OrganizationArray{}
	json.Unmarshal([]byte(resource.Raw), organizations)
	return &organizations.Organizations, nil
}

func (a Auth) ListOrganizations() (*[]Organization, error) {

	path := "/organizations.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return parseOrganizations(resource)

}

func (a Auth) SearchOrganizationByExternalId(external_id string) (*[]Organization, error) {
	data := url.Values{}
	data.Set("external_id", external_id)
	path := fmt.Sprintf("/organizations/search.json?%s", data.Encode())
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return parseOrganizations(resource)
}

func (a Auth) CreateOrganization(org *Organization) (*Organization, error) {
	path := "/organizations.json"
	bytes, err := json.Marshal(SingleOrganization{org})
	if err != nil {
		return nil, err
	}
	resource, err := api(a, "POST", path, string(bytes))
	if err != nil {
		return nil, err
	}
	singleOrg := &SingleOrganization{}
	json.Unmarshal([]byte(resource.Raw), singleOrg)

	return singleOrg.Organization, nil
}

func (a Auth) DeleteOrganization(org_id uint32) error {
	path := fmt.Sprintf("/organizations/%d.json", org_id)

	_, err := api(a, "DELETE", path, "")
	if err != nil {
		return err
	}

	return nil
}

func (a Auth) ListUserOrganizations(user_id string) (*Resource, error) {

	path := "/users/" + user_id + "/organizations.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return resource, nil

}
