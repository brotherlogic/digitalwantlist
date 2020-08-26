package main

import (
	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	purchased = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "digitalwantlist_purchased",
		Help: "The size of the print queue",
	})
)

//ClientUpdate on an updated record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	config, err := s.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	s.processRecord(ctx, client, req.GetInstanceId(), config)

	purchased.Set(float64(len(config.GetPurchased())))
	return &rcpb.ClientUpdateResponse{}, nil
}
