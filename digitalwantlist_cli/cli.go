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
	case "sforce":
		val, _ := strconv.Atoi(os.Args[2])
		client := pbrc.NewClientUpdateServiceClient(conn)
		resp, err := client.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: int32(val)})

		fmt.Printf("%v and %v\n", resp, err)
	}
}
