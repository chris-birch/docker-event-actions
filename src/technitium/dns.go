package technitium

import (
	"context"
	"fmt"
	"github.com/chris-birch/docker-dns-sync/proto/technitium/v1/message"
	"github.com/chris-birch/docker-dns-sync/proto/technitium/v1/service"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Technitium struct {
	client *grpc.ClientConn
	lock   chan bool
	msg    chan *message.DnsRecord
}

func (t *Technitium) Init() {
	// gRPC Client Setup
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient("localhost:50051", opts...) //TODO get from envar
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to server")
	}

	// Channel setup
	t.client = conn
	t.lock = make(chan bool, 1)
	t.lock <- false
}
func (t *Technitium) Close() {
	err := t.client.Close()
	if err != nil {
		fmt.Println("Failed to close connection")
	}
}

func (t *Technitium) SendMsg(rec *message.DnsRecord) {
	select {
	case <-t.lock: // Proceed only if unlocked
		t.msg = make(chan *message.DnsRecord, 5)
		t.msg <- rec
		fmt.Println("lock")

		go func(msg chan *message.DnsRecord) { // Start gRPC stream
			fmt.Println("channel routine start")

			srv := service.NewTechnitiumServiceClient(t.client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
			defer cancel()

			ticker := time.NewTicker(3 * time.Second)
			defer ticker.Stop()

			trecord, err := srv.ProcessRecord(ctx)
			if err != nil {
				log.Fatal().Msgf("could not process record: %v", err)
			}

		loop:
			for {
				select {
				case chmsg := <-msg:
					err := trecord.Send(chmsg)
					if err != nil {
						log.Fatal().Msgf("could not send message: %v", err.Error())
					}
					ticker.Reset(2 * time.Second)
				case <-ticker.C:
					fmt.Println("ticker tick")
					break loop
				}
			}

			// Finished sending msg
			err = trecord.CloseSend()
			if err != nil {
				log.Fatal().Msgf("could not close: %v", err)
			}
			t.lock <- false
			return
		}(t.msg)

	default:
		t.msg <- rec
		fmt.Println("Default")
	}
}

func NewRecord(event events.Message, c *client.Client) (*message.DnsRecord, error) {
	// Validate event data

	//if event.Actor.Attributes["hostname"] == "" || event.Actor.Attributes["domain"] == "" {
	//	return nil, errors.New("event does not contain hostname or domain")
	//}

	// Connect to Docker to get container info
	inspect, err := c.ContainerInspect(context.Background(), event.Actor.ID)
	if err != nil {
		return nil, err
	}

	//rec := new(Record)
	//rec.dnsRecordName = inspect.Config.Hostname
	//rec.dnsRecordType = event.Actor.Attributes["type"]
	//rec.dnsRecordData = string(event.Action)

	rec := message.DnsRecord{
		Name: inspect.Config.Hostname,
		Type: message.Type_TYPE_CNAME,
		Data: string(event.Action),
	}

	return &rec, nil
}

//func Process(record *Record, c *Client) {
//
//	// Process event and create pb message
//	// Create new go routine that handles all streaming (including setting up and taking down the client)
//	// Send pb messages to this routine via a channel
//	// When the client is taken down and the stream closed, end the routine
//	// Check if we already have a routine using a chanel block (unblocked when the routine ends)
//	// Use a select to check on the routine status
//	// The gRPC client should be ended when main() ends, same as the docker one
//
//	srv := service.NewTechnitiumServiceClient(c.conn)
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	trecord, err := srv.ProcessRecord(ctx)
//	if err != nil {
//		log.Fatalf("could not process record: %v", err)
//	}
//
//	rec := message.DnsRecord{
//		Name: record.dnsRecordName,
//		Type: message.Type_TYPE_CNAME,
//		Data: record.dnsRecordData,
//	}
//
//	err = trecord.Send(&rec)
//	if err != nil {
//		log.Fatalf("could not send dns record: %v", err)
//	}
//
//	err = trecord.CloseSend()
//	if err != nil {
//		log.Fatalf("could not close: %v", err)
//	}
//
//}
