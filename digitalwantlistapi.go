package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (s *Server) adjust(ctx context.Context, client rcpb.RecordCollectionServiceClient, record *rcpb.Record) error {
	// Only process 12 inches
	if record.GetRelease().GetFolderId() != int32(242017) {
		return s.unwant(ctx, record)
	}

	//Unwant anything that scores under 4
	if record.GetMetadata().GetOverallScore() < 4 {
		return s.unwant(ctx, record)
	}

	purchased := false
	for _, id := range record.GetRelease().GetDigitalVersions() {
		records, err := s.getRecords(ctx, client, id)
		if err != nil {
			return err
		}
		if len(records) > 0 {
			purchased = true
		}
	}

	records, err := s.getBoughtRecords(ctx)
	if err != nil {
		return err
	}

	for _, id := range records {
		for _, oid := range record.GetRelease().GetDigitalVersions() {
			if id == oid {
				purchased = true
			}
		}
	}

	s.Log(fmt.Sprintf("FOUND PURCHASED %v for %v", purchased, record.GetRelease().GetInstanceId()))
	if purchased {
		return s.unwant(ctx, record)
	}
	return s.want(ctx, record)
}

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
	s.adjust(ctx, client, record)

	return &rcpb.ClientUpdateResponse{}, s.adjust(ctx, client, record)
}

func (s *Server) want(ctx context.Context, record *rcpb.Record) error {
	s.Log(fmt.Sprintf("WANTING %v", record.GetRelease().GetInstanceId()))
	conn, err := s.FDialServer(ctx, "recordwants")
	if err != nil {
		return err
	}
	defer conn.Close()
	rwclient := rwpb.NewWantServiceClient(conn)

	for _, dv := range record.GetRelease().GetDigitalVersions() {
		_, err = rwclient.AddWant(ctx, &rwpb.AddWantRequest{ReleaseId: dv})
		if status.Convert(err).Code() == codes.OK || status.Convert(err).Code() == codes.FailedPrecondition {
			_, err = rwclient.Update(ctx, &rwpb.UpdateRequest{Want: &gdpb.Release{Id: dv}, Level: rwpb.MasterWant_ANYTIME})
		} else {
			return err
		}
	}
	return nil
}

func (s *Server) unwant(ctx context.Context, record *rcpb.Record) error {
	s.Log(fmt.Sprintf("UNWANTING %v", record.GetRelease().GetInstanceId()))

	conn, err := s.FDialServer(ctx, "recordwants")
	if err != nil {
		return err
	}
	defer conn.Close()
	rwclient := rwpb.NewWantServiceClient(conn)

	for _, dv := range record.GetRelease().GetDigitalVersions() {
		_, err = rwclient.AddWant(ctx, &rwpb.AddWantRequest{ReleaseId: dv})
		if status.Convert(err).Code() == codes.OK || status.Convert(err).Code() == codes.FailedPrecondition {
			_, err = rwclient.Update(ctx, &rwpb.UpdateRequest{Want: &gdpb.Release{Id: dv}, Level: rwpb.MasterWant_NEVER})
		}

		if err != nil {
			return err
		}
	}
	return nil
}
