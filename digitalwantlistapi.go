package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"

	gdpb "github.com/brotherlogic/godiscogs"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	rwpb "github.com/brotherlogic/recordwants/proto"
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
	record, err := s.processRecord(ctx, client, req.GetInstanceId(), config)

	cdPurchased := false
	for _, purchased := range config.GetPurchased() {
		for _, dv := range record.GetRelease().GetDigitalVersions() {
			if dv == purchased {
				cdPurchased = true
			}
		}
	}

	s.Log(fmt.Sprintf("Here: %v and %v -> %v -> %v", cdPurchased, record.GetMetadata().GetGoalFolder(), record.GetRelease().GetDigitalVersions(), record.GetMetadata().GetOverallScore()))
	if !cdPurchased && record.GetMetadata().GetGoalFolder() == 242017 {
		conn, err := s.FDialServer(ctx, "recordwants")
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		rwclient := rwpb.NewWantServiceClient(conn)
		for _, dv := range record.GetRelease().GetDigitalVersions() {
			_, err = rwclient.AddWant(ctx, &rwpb.AddWantRequest{ReleaseId: dv})
			if err == nil {
				_, err = rwclient.Update(ctx, &rwpb.UpdateRequest{Want: &gdpb.Release{Id: dv}, Level: rwpb.MasterWant_ANYTIME})
			}

			if err != nil {
				return nil, err
			}
		}
	}

	purchased.Set(float64(len(config.GetPurchased())))
	return &rcpb.ClientUpdateResponse{}, nil
}
