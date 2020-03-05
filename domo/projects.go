package domo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Project struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	MemberIDs []int `json:"members,omitempty"`
	CreatedByUserID int `json:"createdBy,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty"`
	Public bool `json:"public,omitempty"`
	Description string `json:"description, omitempty"`
	DueDate time.Time `json:"dueDate, omitempty"`
}
// ProjectsService handles communication with the projects
// related methods of the Domo API.
//
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference
type ProjectsService service


// List the projects the user of the Domo Client Credentials has access to.
//
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference#Retrieve%20all%20projects
func (s *ProjectsService) List(ctx context.Context) ([]*Project, *http.Response, error) {
	u := "v1/projects"
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var projects []*Project
	resp, err := s.client.Do(ctx, req, &projects)
	if err != nil {
		return nil, resp, err
	}
	return projects, resp, nil
}

// Info about a specific Project. Use the special project ID "me" to return your personal project.
//
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference#Retrieve%20individual%20project
func (s *ProjectsService) Info(ctx context.Context, projectID string) (*Project, *http.Response, error) {
	u := fmt.Sprintf("v1/projects/%s", projectID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var p *Project
	resp, err := s.client.Do(ctx, req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, err
}

// Create a project. Name, Members, and Public field are required. description and dueDate are optional.
//
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference#Create%20a%20project
func (s *ProjectsService) Create(ctx context.Context, project Project) (*Project, *http.Response, error) {
	if project.Name == "" {
		return nil, nil, fmt.Errorf("Expected a project Name. Name is a required field to create a Project")
	}
	if len(project.MemberIDs) < 1  {
		return nil, nil, fmt.Errorf("Expected at least one Member ID. MemberIDs is a required field to create a Project")
	}
	u := "v1/projects"
	req, err := s.client.NewRequest("POST", u, project)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var p *Project
	resp, err := s.client.Do(ctx, req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

// Update a project in Domo.
// The Following fields are read-only and cannot be updated with this API:
// `id`
// `members`
// `createdBy`
// `createdDate`
//
// To update Members use the Update Members API Method.
//
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference#Update%20a%20project
func (s *ProjectsService) Update(ctx context.Context, project Project) (*Project, *http.Response, error) {
	if project.ID == "" {
		return nil, nil, fmt.Errorf("Expected a project ID to identify the Project to Update.")
	}
	u := fmt.Sprintf("v1/projects/%s", project.ID)
	updateBody := struct {
		Name string `json:"name,omitempty"`
		Public bool `json:"public,omitempty"`
		Description string `json:"description, omitempty"`
		DueDate time.Time `json:"dueDate, omitempty"`
	}{project.Name, project.Public, project.Description, project.DueDate}
	req, err := s.client.NewRequest("PUT", u, updateBody)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var p *Project
	resp, err := s.client.Do(ctx, req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

// Delete a project in Domo for the given projectID.
// WARNING: This is destructive and cannot be reversed.
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference#Delete%20a%20project
func (s *ProjectsService) Delete(ctx context.Context, projectID string) (*http.Response, error) {
	u := fmt.Sprintf("v1/projects/%s", projectID)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ProjectMembers IDs for a given Domo Project.
//
// Domo API Docs: https://developer.domo.com/docs/projectsandtasks/projects-tasks-api-reference#Retrieve%20project%20members
func (s *ProjectsService) ProjectMembers(ctx context.Context, projectID string) ([]*int, *http.Response, error) {
	u := fmt.Sprintf("v1/projects/%s/members", projectID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")

	var memberIDs []*int
	resp, err := s.client.Do(ctx, req, &memberIDs)
	if err != nil {
		return nil, resp, err
	}
	return memberIDs, resp, nil
}