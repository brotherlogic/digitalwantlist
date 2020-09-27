package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/digitalwantlist/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	rapb "github.com/brotherlogic/recordadder/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	rcpb.RegisterClientUpdateServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "magic", Value: int64(134)},
	}
}

func (s *Server) initConfig() error {
	s.Log(fmt.Sprintf("Initializing Digital Wantlist Config"))

	config := &pb.Config{}

	ctx, cancel := utils.ManualContext("dwl", "dwl", time.Hour*2, false)
	defer cancel()

	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)

	ids, err := client.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_All{true}})
	if err != nil {
		return err
	}

	for _, id := range ids.GetInstanceIds() {
		_, err := s.processRecord(ctx, client, id, config)
		if err != nil {
			return err
		}
	}

	return s.KSclient.Save(ctx, CONFIG, config)
}

func (s *Server) getBoughtRecords(ctx context.Context) ([]int32, error) {
	conn, err := s.FDialServer(ctx, "recordadder")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rapb.NewAddRecordServiceClient(conn)
	resp, err := client.ListQueue(ctx, &rapb.ListQueueRequest{})
	if err != nil {
		return nil, err
	}

	res := []int32{}
	for _, r := range resp.GetRequests() {
		res = append(res, r.GetId())
	}

	return res, err
}

func (s *Server) getRecord(ctx context.Context, client rcpb.RecordCollectionServiceClient, id int32) (*rcpb.Record, error) {
	r, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: id, Validate: false})
	if err != nil {
		return nil, err
	}

	return r.GetRecord(), nil

}

func (s *Server) getRecords(ctx context.Context, client rcpb.RecordCollectionServiceClient, id int32) ([]int32, error) {
	r, err := client.QueryRecords(ctx, &rcpb.QueryRecordsRequest{Query: &rcpb.QueryRecordsRequest_ReleaseId{id}})
	if err != nil {
		return nil, err
	}

	return r.GetInstanceIds(), nil

}

func (s *Server) processRecord(ctx context.Context, client rcpb.RecordCollectionServiceClient, id int32, config *pb.Config) (*rcpb.Record, error) {
	r, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: id, Validate: false})
	if err != nil {
		return nil, err
	}

	found := false
	for _, id := range config.Purchased {
		if id == r.GetRecord().GetRelease().GetId() {
			found = true
		}
	}

	if !found {
		config.Purchased = append(config.Purchased, r.GetRecord().GetRelease().GetId())
	}

	return r.GetRecord(), nil
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("digitalwantlist", false, true)
	if err != nil {
		return
	}

	ctx, cancel := utils.ManualContext("dwlm", "dwlm", time.Minute, false)
	config, err := server.loadConfig(ctx)
	cancel()
	code := status.Convert(err).Code()
	if err != nil {
		ecancel, eerr := server.Elect()
		if eerr != nil {
			log.Fatalf("Quitting: %v", eerr)
		}

		if code == codes.InvalidArgument {
			server.initConfig()
		}
		ecancel()

		// Silent exit if there's not keystore
		if code == codes.NotFound || code == codes.DeadlineExceeded {
			return
		}

		server.Log(fmt.Sprintf("Error Loading Config: %v", err))
		time.Sleep(time.Second * 5)
		log.Fatalf("Error loading: %v", err)
	}

	/*	if len(config.GetPurchased()) > 1 {
		config.Purchased = []int32{1}
		ctx, cancel := utils.ManualContext("dwlw", "dwlw", time.Minute, false)
		defer cancel()
		e2 := server.KSclient.Save(ctx, CONFIG, config)
		time.Sleep(time.Second * 5)
		server.Log(fmt.Sprintf("boune: %v", e2))
		time.Sleep(time.Second * 5)
		log.Fatalf("back")
	}*/

	time.Sleep(time.Second * 5)
	server.Log(fmt.Sprintf("Loaded config and ready to server: %v", len(config.GetPurchased())))
	time.Sleep(time.Second * 5)

	fmt.Printf("%v", server.Serve())
}
