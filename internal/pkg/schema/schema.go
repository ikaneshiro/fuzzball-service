// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package schema

// String returns the GraphQL schema.
func String() string {
	// TODO: splint into .graphql files and go-gen into this file.
	return `
schema {
  query: Query
}

type Query {
  job(id: ID!): Job
}

type Job {
  id: ID!
  name: String!
}
	`
}