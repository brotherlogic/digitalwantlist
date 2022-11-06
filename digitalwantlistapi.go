package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"

	gdpb "github.com/brotherlogic/godiscogs"
	rapb "github.com/brotherlogic/recordadder/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	rspb "github.com/brotherlogic/recordsales/proto"
	wlpb "github.com/brotherlogic/wantslist/proto"
)

var (
	purchased = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "digitalwantlist_purchased",
		Help: "The size of the print queue",
	})
)

func (s *Server) adjust(ctx context.Context, client rcpb.RecordCollectionServiceClient, record *rcpb.Record) error {
	// We only consider 12 inches
	if record.GetMetadata().GetFiledUnder() != rcpb.ReleaseMetadata_FILE_12_INCH {
		s.CtxLog(ctx, fmt.Sprintf("Skipping %v because it's %v", record.GetRelease().GetInstanceId(), record.GetMetadata().GetFiledUnder()))
		return nil
	}

	//Unwant anything that we have partial or full matches on
	if record.GetMetadata().GetMatch() == rcpb.ReleaseMetadata_FULL_MATCH || record.GetMetadata().GetMatch() == rcpb.ReleaseMetadata_PARTIAL_MATCH {
		s.CtxLog(ctx, fmt.Sprintf("UNWATING %v because of match: %v", record.GetRelease().GetInstanceId(), record.GetMetadata().GetMatch()))
		return s.unwant(ctx, record)
	}

	// Only process 12 inches that are in the collection
	if record.GetRelease().GetFolderId() != int32(242017) ||
		record.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_IN_COLLECTION ||
		record.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_IN_THE_BOX {
		s.CtxLog(ctx, fmt.Sprintf("UNWANTING the %v because of situation %v, %v, %v",
			record.GetRelease().GetInstanceId(),
			record.GetRelease().GetFolderId() != int32(242017),
			record.GetMetadata().GetCategory() != rcpb.ReleaseMetadata_IN_COLLECTION,
			record.GetMetadata().GetBoxState() == rcpb.ReleaseMetadata_IN_THE_BOX))
		return s.unwant(ctx, record)
	}

	// If it's a digital keeper , then want it
	if record.GetMetadata().GetKeep() == rcpb.ReleaseMetadata_DIGITAL_KEEPER {
		return s.want(ctx, record, "digital_quick")
	}

	// Don't process keepers
	if record.GetMetadata().GetKeep() == rcpb.ReleaseMetadata_KEEPER {
		s.CtxLog(ctx, fmt.Sprintf("UNWANTING %v because of keeper", record.GetRelease().GetInstanceId()))
		return s.unwant(ctx, record)
	}

	//Unwant anything that scores under 4
	if record.GetRelease().GetRating() < 4 {
		s.CtxLog(ctx, fmt.Sprintf("UNWANTING %v because of overall score", record.GetRelease().GetInstanceId()))
		return s.unwant(ctx, record)
	}

	s.CtxLog(ctx, fmt.Sprintf("WANTING %v", record.GetRelease().GetInstanceId()))
	return s.want(ctx, record, "digital")
}

// ClientUpdate on an updated record
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
	record, _ := s.processRecord(ctx, client, req.GetInstanceId(), config)
	//s.adjust(ctx, client, record)

	return &rcpb.ClientUpdateResponse{}, s.adjust(ctx, client, record)
}

func (s *Server) want(ctx context.Context, record *rcpb.Record, list string) error {
	conn, err := s.FDialServer(ctx, "wantslist")
	if err != nil {
		return err
	}
	defer conn.Close()
	rwclient := wlpb.NewWantServiceClient(conn)

	c2, err2 := s.FDialServer(ctx, "recordsales")
	if err2 != nil {
		return err
	}
	defer c2.Close()
	rsclient := rspb.NewSaleServiceClient(c2)

	presp, err := rsclient.GetPrice(ctx, &rspb.GetPriceRequest{Ids: record.GetRelease().GetDigitalVersions()})
	if err != nil {
		return err
	}

	for _, dv := range record.GetRelease().GetDigitalVersions() {

		if presp.GetPrices()[dv].GetLatest().GetPrice() < float32(record.GetMetadata().GetCurrentSalePrice())/100 {
			_, err = rwclient.AddWantListItem(ctx, &wlpb.AddWantListItemRequest{ListName: list, Entry: &wlpb.WantListEntry{Want: dv}})
			if err != nil {
				return err
			}
		} else {
			s.CtxLog(ctx, fmt.Sprintf("Price mismatch %v vs %v", presp.GetPrices()[dv].GetLatest().GetPrice(), float32(record.GetMetadata().GetCurrentSalePrice())/100))
		}
	}
	return nil
}

func (s *Server) unwant(ctx context.Context, record *rcpb.Record) error {
	conn, err := s.FDialServer(ctx, "wantslist")
	if err != nil {
		return err
	}
	defer conn.Close()
	rwclient := wlpb.NewWantServiceClient(conn)

	for _, dv := range record.GetRelease().GetDigitalVersions() {
		_, err = rwclient.DeleteWantListItem(ctx, &wlpb.DeleteWantListItemRequest{ListName: "digital", Entry: &wlpb.WantListEntry{Want: dv}})
		if err != nil {
			return err
		}
	}
	return nil
}

// ClientAddUpdate deal with a new record addition from record adder
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
