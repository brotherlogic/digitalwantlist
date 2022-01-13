package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func main() {
	ctx, cancel := utils.ManualContext("dwcli", time.Hour)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "digitalwantlist")
	if err != nil {
		log.Fatalf("Pah: %v", err)
	}
	defer conn.Close()

	switch os.Args[1] {
	case "fullping":
		client := pbrc.NewClientUpdateServiceClient(conn)
		ctx2, cancel2 := utils.ManualContext("recordcollectioncli-"+os.Args[1], time.Hour)
		defer cancel2()

		conn2, err := utils.LFDialServer(ctx2, "recordcollection")
		if err != nil {
			log.Fatalf("Cannot reach rc: %v", err)
		}
		defer conn2.Close()

		registry := pbrc.NewRecordCollectionServiceClient(conn2)
		ids, err := registry.QueryRecords(ctx2, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_FolderId{673768}})
		if err != nil {
			log.Fatalf("Bad query: %v", err)
		}
		for i, id := range ids.GetInstanceIds() {
			_, err := client.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: id})

			fmt.Printf("%v and %v\n", i, err)
		}
	case "sforce":
		val, _ := strconv.Atoi(os.Args[2])
		client := pbrc.NewClientUpdateServiceClient(conn)
		resp, err := client.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: int32(val)})

		fmt.Printf("%v and %v\n", resp, err)
	}
}
