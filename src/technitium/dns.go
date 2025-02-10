package technitium

import (
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

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
			value, _ := q.fifo.DequeueOrWaitForNextElement()
			post(value.(string))
			//log.Info().Msg(value.(string))
			q.wg.Done()
		}
	}()
}

func (q *queue) push(s string) {
	if q.done == nil {
		q.init()
	}
	q.wg.Add(1)
	err := q.fifo.Enqueue(s)

	if err != nil {
		fmt.Println(err)
	}
}

func (q *queue) action(a func()) {

}

var q = newQueue()

type Event struct {
	message string
}

func NewEvent(event events.Message) *Event {

	if event.Actor.Attributes["hostname"] != "" {
		return &Event{message}
	}

	if message != "" && message != "testing" {

	} else {
		log.Info().Msgf("%s contains testing:", message)
		return nil
	}
}

func (d *Event) Process() {
	q.push(d.message).action(func() {

		log.Info().Msg(d.message)

	})
}

func post(message string) {
	r := rand.Intn(2)
	time.Sleep(time.Duration(r) * time.Second)
	fmt.Printf("Message: %s\n", message)
}
