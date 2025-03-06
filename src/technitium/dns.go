package technitium

import (
	"context"
	"fmt"
	"github.com/chris-birch/docker-dns-sync/proto/technitium/v1/message"
	"github.com/chris-birch/docker-dns-sync/proto/technitium/v1/service"
	"github.com/docker/docker/api/types/container"
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
	t.msg = make(chan *message.DnsRecord, 25)
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
		log.Debug().Msg("Create new service routine and send msg")
		t.msg <- rec

		go func(msg chan *message.DnsRecord) { // Start gRPC stream
			log.Debug().Msg("Started gRPC service routine")

			srv := service.NewTechnitiumServiceClient(t.client)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			trecord, err := srv.ProcessRecord(ctx)
			if err != nil {
				log.Error().Msgf("could not setup ClientStreamingClient: %v", err)
			}

		loop:
			for {
				select {
				case chmsg := <-msg:
					log.Debug().Msgf("Sending msg to gRPC server: %v", chmsg)
					err := trecord.Send(chmsg)
					if err != nil {
						log.Fatal().Msgf("could not send message: %v", err.Error())
					}
					ticker.Reset(2 * time.Second)
				case <-ticker.C:
					log.Debug().Msg("gRPC service routine timeout. Exiting loop.")
					break loop
				}
			}

			// Finished sending all msg
			err = trecord.CloseSend()
			if err != nil {
				log.Fatal().Msgf("could not close: %v", err)
			}
			t.lock <- false
			return
		}(t.msg)

	default:
		log.Debug().Msg("Sending record to exiting service routine")
		t.msg <- rec
	}
}

func NewRecord(event events.Message, c *client.Client, dockHost string) (*message.DnsRecord, error) {
	t := make(map[events.Action]message.Action)
	t["create"] = message.Action_ACTION_CREATE
	t["attach"] = message.Action_ACTION_ATTACH
	t["start"] = message.Action_ACTION_START
	t["die"] = message.Action_ACTION_DIE

	// Check that the event type is supported
	if !validateEvent(event, t) {
		log.Debug().Msgf("Ignoring container %v, with event type %v", event.Actor.ID, event.Action)
		return nil, nil
	}

	//if event.Actor.Attributes["hostname"] == "" || event.Actor.Attributes["domain"] == "" {
	//	return nil, errors.New("event does not contain hostname or domain")
	//}

	cntr, err := inspect(c, &event)
	if err != nil {
		return nil, err
	}

	rec := new(message.DnsRecord)
	rec.Name = cntr.Config.Hostname
	rec.Type = message.Type_TYPE_CNAME // It's always going to be a CNAME
	rec.Data = dockHost
	rec.Action = t[event.Action]
	rec.ContainerId = event.Actor.ID
	return rec, nil
}

func inspect(c *client.Client, event *events.Message) (container.InspectResponse, error) {
	// Connect to Docker to get container info
	inspect, err := c.ContainerInspect(context.Background(), event.Actor.ID)
	if err != nil {
		return inspect, err
	}
	return inspect, nil
}

func validateEvent(event events.Message, act map[events.Action]message.Action) bool {
	var valid bool
	valid = false
	for k := range act {
		if event.Action == k {
			valid = true
			break
		}
	}
	return valid
}
