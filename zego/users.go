package zego

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UserArray struct {
	Users []User `json:"users"`
}

type SingleUser struct {
	User *User `json:"user"`
}

type Tags struct {
	Tags []string `json:"tags"`
}

type User struct {
	Id                    uint32         `json:"id,omitempty"`
	Url                   string         `json:"url,omitempty"`
	Name                  string         `json:"name,omitempty"`
	External_id           string         `json:"external_id,omitempty"`
	Alias                 string         `json:"alias,omitempty"`
	Created_at            string         `json:"created_at,omitempty"`
	Updated_at            string         `json:"updated_at,omitempty"`
	Active                bool           `json:"active,omitempty"`
	Verified              bool           `json:"verified,omitempty"`
	Shared                bool           `json:"shared,omitempty"`
	Shared_agent          bool           `json:"shared_agent,omitempty"`
	Locale                string         `json:"locale,omitempty"`
	Locale_id             uint32         `json:"locale_id,omitempty"`
	Time_zone             string         `json:"time_zone,omitempty"`
	Last_login_at         string         `json:"last_login_at,omitempty"`
	Email                 string         `json:"email,omitempty"`
	Phone                 string         `json:"phone,omitempty"`
	Signature             string         `json:"signature,omitempty"`
	Details               string         `json:"details,omitempty"`
	Notes                 string         `json:"notes,omitempty"`
	Organization_id       uint32         `json:"organization_id,omitempty"`
	Role                  string         `json:"role,omitempty"`
	Customer_role_id      uint32         `json:"custom_role_id,omitempty"`
	Moderator             bool           `json:"moderator,omitempty"`
	Ticket_restriction    string         `json:"ticket_restriction,omitempty"`
	Only_private_comments bool           `json:"only_private_comments,omitempty"`
	Tags                  []string       `json:"tags,omitempty"`
	Restricted_agent      bool           `json:"restricted_agent,omitempty"`
	Suspended             bool           `json:"suspended,omitempty"`
	Photo                 []*Users_Photo `json:"photo,omitempty"`
	User_fields           []*User_Field  `json:"user_fields,omitempty"`
}

type Users_Photo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ContentUrl  string `json:"content_url"`
	ContentType string `json:"content_type"`
	Size        int    `json:"size"`
}

type User_Field struct {
	UserDecimal  float32 `json:"user_decimal"`
	UserDropdown string  `json:"user_dropdown"`
	UserDate     string  `json:"user_date"`
}

func parseUsers(resource *Resource) (*[]User, error) {
	users := &UserArray{}
	json.Unmarshal([]byte(resource.Raw), users)
	return &users.Users, nil
}

func (a Auth) ListUsers() (*[]User, error) {
	path := "/users.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return parseUsers(resource)
}

func (a Auth) ShowUser(user_id string) (*User, error) {
	UserStruct := &SingleUser{}

	path := "/users/" + user_id + ".json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(resource.Raw), UserStruct)

	return UserStruct.User, nil
}

func (a Auth) CreateUser(user *User) (*User, error) {
	path := "/users.json"
	bytes, err := json.Marshal(SingleUser{user})
	if err != nil {
		return nil, err
	}
	resource, err := api(a, "POST", path, string(bytes))
	if err != nil {
		return nil, err
	}
	singleUser := &SingleUser{}
	json.Unmarshal([]byte(resource.Raw), singleUser)

	return singleUser.User, nil
}

func (a Auth) UpdateUser(user *User) (*User, error) {
	path := fmt.Sprintf("/users/%d.json", user.Id)
	bytes, err := json.Marshal(SingleUser{user})
	if err != nil {
		return nil, err
	}
	resource, err := api(a, "PUT", path, string(bytes))
	if err != nil {
		return nil, err
	}

	singleUser := &SingleUser{}
	json.Unmarshal([]byte(resource.Raw), singleUser)

	return singleUser.User, nil
}

func (a Auth) SetUserTags(user_id uint32, tags []string) (*[]string, error) {
	path := fmt.Sprintf("/users/%d/tags.json", user_id)
	bytes, err := json.Marshal(Tags{tags})
	if err != nil {
		return nil, err
	}
	resource, err := api(a, "POST", path, string(bytes))
	if err != nil {
		return nil, err
	}

	tagsResponse := &Tags{}
	json.Unmarshal([]byte(resource.Raw), tagsResponse)

	return &tagsResponse.Tags, nil
}

func (a Auth) DeleteUser(user_id uint32) error {
	path := fmt.Sprintf("/users/%d.json", user_id)
	_, err := api(a, "DELETE", path, "")
	if err != nil {
		return err
	}
	return nil
}

func (a Auth) SearchUserByExternalId(external_id string) (*[]User, error) {
	data := url.Values{}
	data.Set("external_id", external_id)
	path := fmt.Sprintf("/users/search.json?%s", data.Encode())
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return parseUsers(resource)
}

func (a Auth) ShowUserRelated(user_id string) (*Resource, error) {

	path := "/users/" + user_id + "/related.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return resource, nil

}

func (a Auth) ListCollaborators(ticket_id string) (*Resource, error) {

	path := "/tickets/" + ticket_id + "/collaborators.json"
	resource, err := api(a, "GET", path, "")
	if err != nil {
		return nil, err
	}

	return resource, nil

}
