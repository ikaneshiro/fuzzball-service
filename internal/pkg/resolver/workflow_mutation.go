// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/fuzzball-service/internal/pkg/core"
)

// CreateWorkflow creates a new workflow.
func (r Resolver) CreateWorkflow(ctx context.Context, args struct {
	Spec core.WorkflowSpec
}) (*WorkflowResolver, error) {
	w, err := r.s.CreateWorkflow(ctx, args.Spec)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{w}, nil
}

// DeleteWorkflow deletes a workflow.
func (r Resolver) DeleteWorkflow(ctx context.Context, args struct {
	ID string
}) (*WorkflowResolver, error) {
	w, err := r.s.DeleteWorkflow(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{w}, nil
}
