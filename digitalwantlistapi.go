package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gdpb "github.com/brotherlogic/godiscogs"
	rapb "github.com/brotherlogic/recordadder/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	rspb "github.com/brotherlogic/recordsales/proto"
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

	// Don't process keepers
	if record.GetMetadata().GetKeep() == rcpb.ReleaseMetadata_KEEPER {
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

	c2, err2 := s.FDialServer(ctx, "recordsales")
	if err2 != nil {
		return err
	}
	defer c2.Close()
	rsclient := rspb.NewSaleServiceClient(c2)

	for _, dv := range record.GetRelease().GetDigitalVersions() {
		sprice, err := rsclient.GetPrice(ctx, &rspb.GetPriceRequest{Id: dv})
		if err != nil {
			return err
		}

		if sprice.GetPrices().GetLatest().GetPrice() < float32(record.GetMetadata().GetSalePrice())/100 {
			_, err = rwclient.AddWant(ctx, &rwpb.AddWantRequest{ReleaseId: dv})
			if status.Convert(err).Code() == codes.OK || status.Convert(err).Code() == codes.FailedPrecondition {
				_, err = rwclient.Update(ctx, &rwpb.UpdateRequest{Want: &gdpb.Release{Id: dv}, Level: rwpb.MasterWant_WANT_DIGITAL})
			} else {
				return err
			}
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

//ClientAddUpdate deal with a new record addition from record adder
func (s *Server) ClientAddUpdate(ctx context.Context, req *rapb.ClientAddUpdateRequest) (*rapb.ClientAddUpdateResponse, error) {
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rclient := rcpb.NewRecordCollectionServiceClient(conn)
	rel, err := rclient.GetRecord(ctx, &rcpb.GetRecordRequest{ReleaseId: req.GetId()})
	if err != nil {
		return nil, err
	}

	iids, err := rclient.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_MasterId{MasterId: rel.GetRecord().GetRelease().GetMasterId()}})
	if err != nil {
		return nil, err
	}

	for _, iid := range iids.GetInstanceIds() {
		_, err := rclient.UpdateRecord(ctx, &rcpb.UpdateRecordRequest{Reason: "dwl-blank", Update: &rcpb.Record{Release: &gdpb.Release{InstanceId: iid}}})
		if err != nil {
			return nil, err
		}
	}

	return &rapb.ClientAddUpdateResponse{}, nil
}
