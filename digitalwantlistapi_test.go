package main

import (
	"context"
	"testing"

	keystoreclient "github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/digitalwantlist/proto"
	rapb "github.com/brotherlogic/recordadder/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), CONFIG, &pb.Config{})
	return s
}

func TestClientUpdate(t *testing.T) {
	s := InitTestServer()
	_, err := s.ClientUpdate(context.Background(), &rcpb.ClientUpdateRequest{})
	if err != nil {
		t.Errorf("Could not perform update: %v", err)
	}
}

func TestClientAddUpdate(t *testing.T) {
	s := InitTestServer()
	_, err := s.ClientAddUpdate(context.Background(), &rapb.ClientAddUpdateRequest{})
	if err != nil {
		t.Errorf("Could not perform update: %v", err)
	}
}
