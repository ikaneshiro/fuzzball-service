// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"os"
	"runtime"

	"github.com/pbnjay/memory"
)

// IOFetcher is the interface where IO data is retrieved.
type IOFetcher interface {
	JobOutputFetcher
}

// Persister is the interface by which all data is persisted.
type Persister interface {
	WorkflowPersister
	JobPersister
	VolumePersister
}

// Resolver is the root type for resolving GraphQL queries.
type Resolver struct {
	p  Persister
	f  IOFetcher
	si SystemInfo
}

// New creates a new GraphQL resolver.
func New(p Persister, f IOFetcher) (*Resolver, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Resolver{
		p: p,
		f: f,
		si: SystemInfo{
			HostName:        hostName,
			CPUArchitecture: runtime.GOARCH,
			OSPlatform:      runtime.GOOS,
			Memory:          memory.TotalMemory() / 1024 / 1024, // megabytes
		},
	}, nil
}
