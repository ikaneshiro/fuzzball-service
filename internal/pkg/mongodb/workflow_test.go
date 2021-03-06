// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package mongodb

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/sylabs/fuzzball-service/internal/pkg/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// getTestWorkflow generates a workflow for use in testing. The attributes of the workflow are
// varied based on the value of i.
func getTestWorkflow(i int32) core.Workflow {
	return core.Workflow{
		Name: fmt.Sprintf("workflow-%02d", i),
	}
}

// insertTestWorkflow inserts a workflow into the DB.
func insertTestWorkflow(t *testing.T, db *mongo.Database) core.Workflow {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt32)))
	if err != nil {
		t.Fatalf("failed to generate random int: %v", err)
	}
	j := getTestWorkflow(int32(n.Int64()))
	sr, err := db.Collection(workflowCollectionName).InsertOne(context.Background(), j)
	if err != nil {
		t.Fatalf("failed to insert: %s", err)
	}
	j.ID = sr.InsertedID.(primitive.ObjectID).Hex()
	return j
}

// deleteTestWorkflow deletes a workflow.
func deleteTestWorkflow(t *testing.T, db *mongo.Database, id string) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Fatalf("failed to parse object ID: %v", err)
	}
	m := bson.M{"_id": oid}
	if err := db.Collection(workflowCollectionName).FindOneAndDelete(context.Background(), m).Err(); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}
}

func TestCreateWorkflow(t *testing.T) {
	orig := core.Workflow{
		ID:   "blah",
		Name: "test",
	}

	// Create should succeed.
	w, err := testConnection.CreateWorkflow(context.Background(), orig)
	if err != nil {
		t.Fatalf("failed to create: %s", err)
	}
	defer deleteTestWorkflow(t, testConnection.db, w.ID)

	// Verify returned workflow. Force ID and CreatedAt since they are set by CreateWorkflow.
	orig.ID = w.ID
	orig.CreatedAt = w.CreatedAt
	if _, err := primitive.ObjectIDFromHex(w.ID); err != nil {
		t.Fatalf("workflow has invalid ID")
	}
	if got, want := w, orig; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Get should succeed.
	w, err = testConnection.GetWorkflow(context.Background(), w.ID)
	if err != nil {
		t.Fatalf("failed to get: %s", err)
	}

	// Verify returned workflow.
	if got, want := w, orig; !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDeleteWorkflow(t *testing.T) {
	w := insertTestWorkflow(t, testConnection.db)

	// Delete should succeed.
	if _, err := testConnection.DeleteWorkflow(context.Background(), w.ID); err != nil {
		t.Fatalf("failed to delete: %s", err)
	}

	// Get should fail.
	if _, err := testConnection.GetWorkflow(context.Background(), w.ID); err == nil {
		t.Error("unexpected success")
	}

	// Delete should fail the second time.
	if _, err := testConnection.DeleteWorkflow(context.Background(), w.ID); err == nil {
		t.Error("unexpected success")
	}

	// Delete should fail with bad BSON ID.
	if _, err := testConnection.DeleteWorkflow(context.Background(), "oops"); err == nil {
		t.Error("unexpected success")
	}
}

func TestGetWorkflow(t *testing.T) {
	w := insertTestWorkflow(t, testConnection.db)
	defer deleteTestWorkflow(t, testConnection.db, w.ID)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"Found", w.ID, false},
		{"NotFound", primitive.NewObjectID().Hex(), true},
		{"BadID", "1234", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testConnection.GetWorkflow(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetWorkflowStatus(t *testing.T) {
	w := insertTestWorkflow(t, testConnection.db)
	defer deleteTestWorkflow(t, testConnection.db, w.ID)

	// Get should return an unfinished workflow.
	w, err := testConnection.GetWorkflow(context.Background(), w.ID)
	if err != nil {
		t.Error("unexpected failure")
	}

	if w.Status != "" {
		t.Errorf("unexpected status: got %q, want %q", w.Status, "")
	}

	if err := testConnection.SetWorkflowStatus(context.Background(), w.ID, "newStatus"); err != nil {
		t.Error("unexpected failure")
	}

	// Get should return a finished workflow.
	w, err = testConnection.GetWorkflow(context.Background(), w.ID)
	if err != nil {
		t.Error("unexpected failure")
	}

	if w.Status != "newStatus" {
		t.Errorf("unexpected status: got %q, want %q", w.Status, "newStatus")
	}
}
