package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/digitalwantlist/proto"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true

	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), CONFIG, &pb.Config{})

	return s
}

func TestLoad(t *testing.T) {
	s := InitTest()
	_, err := s.loadConfig(context.Background())
	if err != nil {
		t.Errorf("Bad config load: %v", err)
	}
}

func TestBadLoad(t *testing.T) {
	s := InitTest()
	s.GoServer.KSclient.Fail = true
	conf, err := s.loadConfig(context.Background())
	if err == nil {
		t.Errorf("Bad config load: %v", conf)
	}
}

func TestBasic(t *testing.T) {
	s := InitTest()
	err := s.runComputation(context.Background())
	if err != nil {
		t.Errorf("Bad run: %v", err)
	}
}
