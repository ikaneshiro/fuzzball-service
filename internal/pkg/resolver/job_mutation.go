// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// CreateJob creates a new job.
func (r Resolver) CreateJob(ctx context.Context, args struct {
	Name string
}) (*JobResolver, error) {
	j := model.Job{
		Name: args.Name,
	}
	j, err := r.p.CreateJob(ctx, j)
	if err != nil {
		return nil, err
	}
	return &JobResolver{&j}, nil
}
