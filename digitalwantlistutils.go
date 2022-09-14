package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/digitalwantlist/proto"
)

const (
	// CONFIG - Where we store incoming requests
	CONFIG = "/github.com/brotherlogic/digitalwantslist/queue"
)

func (s *Server) runComputation(ctx context.Context) error {
	t := time.Now()
	sum := 0
	for i := 0; i < 10000; i++ {
		sum += i
	}
	s.CtxLog(ctx, fmt.Sprintf("Sum is %v -> %v", sum, time.Now().Sub(t).Nanoseconds()/1000000))
	return nil
}

func (s *Server) loadConfig(ctx context.Context) (*pb.Config, error) {
	data, _, err := s.KSclient.Read(ctx, CONFIG, &pb.Config{})
	if err != nil {
		return nil, err
	}
	config := data.(*pb.Config)
	return config, nil

}
