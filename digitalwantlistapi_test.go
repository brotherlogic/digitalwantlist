package main

import (
	"context"
	"testing"

	keystoreclient "github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/digitalwantlist/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), CONFIG, &pb.Config{})
	return s
}

func TestClientUpdateFail(t *testing.T) {
	// Pass
}
