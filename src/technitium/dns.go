package technitium

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"sync"
)

const TOKEN = "0663aea32e1520aeefd175f8f9b9656394ac8012568259fd1dce0b0ebbe4bf18"

type queue struct {
	fifo *goconcurrentqueue.FIFO
	done chan struct{}
	wg   sync.WaitGroup
}

func newQueue() *queue {
	return new(queue)
}

func (q *queue) init() {
	q.fifo = goconcurrentqueue.NewFIFO()
	q.done = make(chan struct{})

	go func() {
		q.wg.Wait()
		q.done <- struct{}{}
	}()

	go func() {
		for {
			value, err := q.fifo.DequeueOrWaitForNextElement()
			// switch on method to determine which func to call
			updateDNS(&value)

			log.Fatal().Err(err).Msg("Queue is full")
			q.wg.Done()
		}
	}()

}

func (q *queue) push(r *Record) {
	if q.done == nil {
		q.init()
	}

	q.wg.Add(1)
	err := q.fifo.Enqueue(r)

	if err != nil {
		fmt.Println(err)
	}

}

var q = newQueue()

type Record struct {
	dnsRecordName string
	dnsRecordType string
	dnsRecordData string
}

func NewRecord(event events.Message) (*Record, error) {

	if event.Actor.Attributes["hostname"] == "" || event.Actor.Attributes["domain"] == "" {
		return nil, errors.New("event does not contain hostname or domain")
	}

	rec := new(Record)
	rec.dnsRecordName = event.Actor.Attributes["hostname"]
	rec.dnsRecordType = event.Actor.Attributes["type"]
	rec.dnsRecordData = event.Actor.Attributes["content"]

	fmt.Println(rec.dnsRecordName, rec.dnsRecordType, rec.dnsRecordData)

	return rec, nil
}

func (r *Record) Process() {
	q.push(r)
}

func updateDNS(record Record) {
	//baseUrl := "http://192.168.2.4:5380/api/zones/"

	// Check if the domain already exists
	//http://192.168.2.4:5380/api/zones/records/get
	// token
	// domain
	// list zone = true
	resp, err := http.Get("http://192.168.2.4:5380/api/zones/records/get")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // response body is []byte

	var result dnsRecord

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	PrettyPrint(result)
}

func PrettyPrint(i interface{}) string {

	s, _ := json.MarshalIndent(i, "", "\t")

	return string(s)

}

//resp, err := http.Get(url)
//if err != nil {
//	// we will get an error at this stage if the request fails, such as if the
//	// requested URL is not found, or if the server is not reachable.
//	log.Fatal().Err(err).Msg("Failed to send request")
//}
//defer resp.Body.Close()
//
//// if we want to check for a specific status code, we can do so here
//// for example, a successful request should return a 200 OK status
//if resp.StatusCode != http.StatusOK {
//	// if the status code is not 200, we should log the status code and the
//	// status string, then exit with a fatal error
//	//log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
//	log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)

//}

//// print the response
//data, err := io.ReadAll(resp.Body)
//if err != nil {
//	log.Fatal().Err(err).Msg("Failed to read response")
//}
//fmt.Println(string(data))
//}
