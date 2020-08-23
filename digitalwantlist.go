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

	pb "github.com/brotherlogic/digitalwantlist/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
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
		&pbg.State{Key: "magic", Value: int64(13)},
	}
}

func (s *Server) initConfig() error {
	s.Log(fmt.Sprintf("Initializing Digital Wantlist Config"))
	cancel, err := s.Elect()
	defer cancel()

	if err != nil {
		return err
	}

	config := &pb.Config{}

	ctx, cancel := utils.ManualContext("dwl", "dwl", time.Hour*2, false)
	defer cancel()

	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)

	ids, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_All{true}})
	if err != nil {
		return err
	}

	for _, id := range ids.GetInstanceIds() {
		r, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: id, Validate: false})
		if err != nil {
			return err
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
	}

	return nil
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
	if err != nil {
		log.Fatalf("Error loading: %v", err)
	}

	server.Log(fmt.Sprintf("Loaded config and ready to server: %v", len(config.GetPurchased())))
	time.Sleep(time.Second * 5)

	fmt.Printf("%v", server.Serve())
}
