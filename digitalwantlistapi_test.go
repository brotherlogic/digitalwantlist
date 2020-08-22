package main

import (
	"context"
	"testing"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

func InitTestServer() *Server {
	return Init()
}

func TestClientUpdate(t *testing.T) {
	s := InitTestServer()
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})
	if err != nil {
		t.Errorf("Could not perform update: %v", err)
	}
}
