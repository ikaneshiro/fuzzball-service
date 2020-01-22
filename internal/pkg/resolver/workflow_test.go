// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

type mockPersister struct {
	wantPA model.PageArgs
	j      model.Job
	w      model.Workflow
	jp     model.JobsPage
	wp     model.WorkflowsPage
}

func (p *mockPersister) CreateWorkflow(ctx context.Context, w model.Workflow) (model.Workflow, error) {
	if got, want := w.Name, p.w.Name; got != want {
		return model.Workflow{}, fmt.Errorf("got name %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) DeleteWorkflow(ctx context.Context, id string) (model.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return model.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) GetWorkflow(ctx context.Context, id string) (model.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return model.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) GetWorkflows(ctx context.Context, pa model.PageArgs) (model.WorkflowsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return model.WorkflowsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.wp, nil
}

func (p *mockPersister) CreateJob(ctx context.Context, j model.Job) (model.Job, error) {
	if got, want := j.Name, p.j.Name; got != want {
		return model.Job{}, fmt.Errorf("got name %v, want %v", got, want)
	}
	return p.j, nil
}

func (p *mockPersister) DeleteJobsByWorkflowID(ctx context.Context, id string) error {
	if got, want := id, p.j.WorkflowID; got != want {
		return fmt.Errorf("got ID %v, want %v", got, want)
	}
	return nil
}

func (p *mockPersister) GetJob(ctx context.Context, id string) (model.Job, error) {
	if got, want := id, p.j.ID; got != want {
		return model.Job{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.j, nil
}

func (p *mockPersister) GetJobs(ctx context.Context, pa model.PageArgs) (model.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return model.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, nil
}

func (p *mockPersister) GetJobsByWorkflowID(ctx context.Context, pa model.PageArgs, id string) (model.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return model.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, nil
}

func TestWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: model.Workflow{
				ID:   "workflowID",
				Name: "workflowName",
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		id   string
	}{
		{"OK", "workflowID"},
		{"BadID", "bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			query OpName($id: ID!)
			{
			  workflow(id: $id) {
			    id
			    name
			    createdBy {
			      id
			      login
			    }
			    createdAt
			    startedAt
			    finishedAt
			  }
			}`
			args := map[string]interface{}{
				"id": tt.id,
			}

			res := s.Exec(context.Background(), q, "", args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: model.Workflow{
				ID:   "workflowID",
				Name: "workflowName",
			},
			j: model.Job{
				ID:         "jobID",
				WorkflowID: "workflowID",
				Name:       "jobName",
				Image:      "jobImage",
				Command:    []string{"jobCommand"},
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	okMap := map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "workflowName",
			"jobs": map[string]interface{}{
				"name":    "jobName",
				"image":   "jobImage",
				"command": "jobCommand",
			},
		},
	}

	badMap := map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "bad",
			"jobs": map[string]interface{}{
				"name":    "jobName",
				"image":   "jobImage",
				"command": "jobCommand",
			},
		},
	}

	tests := []struct {
		name string
		vars map[string]interface{}
	}{
		{"OK", okMap},
		{"BadName", badMap},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			mutation OpName($spec: WorkflowSpec!) {
			  createWorkflow(spec: $spec) {
			    id
			    name
			    createdBy {
			      id
			      login
			    }
			    createdAt
			    startedAt
			    finishedAt
			  }
			}`

			res := s.Exec(context.Background(), q, "", tt.vars)
			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: model.Workflow{
				ID:   "workflowID",
				Name: "workflowName",
			},
			j: model.Job{
				ID:         "jobID",
				WorkflowID: "workflowID",
				Name:       "jobName",
				Image:      "jobImage",
				Command:    []string{"jobCommand"},
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		id   string
	}{
		{"OK", "workflowID"},
		{"BadID", "bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			mutation OpName($id: ID!) {
			  deleteWorkflow(id: $id) {
			    id
			    name
			    createdBy {
			      id
			      login
			    }
			    createdAt
			    startedAt
			    finishedAt
			  }
			}`

			args := map[string]interface{}{
				"id": tt.id,
			}

			res := s.Exec(context.Background(), q, "", args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
