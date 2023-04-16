package direct_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/murtaza-u/ddos/direct"
	"github.com/murtaza-u/ddos/mock/direct"
)

func TestSrvDirect(t *testing.T) {
	mock.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	director := mock.NewSrvDirector(ctx)
	out := director.Outlet()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		err := director.Start()
		if err != nil {
			log.Printf("[Director] Start: %v\n", err)
		}
		cancel()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-ctx.Done():
				log.Println("[send]: for-select: context canceled")
				break Loop
			case <-t.C:
				err := director.Direct(&mock.Request{Id: "XXX"})
				if err != nil {
					log.Printf("[send]: for-select: %v\n", err)
					cancel()
					break Loop
				}
			}
		}

		wg.Done()
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("[receive]: for-select: context canceled")
			break Loop
		case req := <-out:
			log.Printf(
				"received request: id: %s | method: %s | issuer: %s",
				req.GetId(), req.GetMethod(), req.GetIssuer(),
			)
		}
	}

	wg.Wait()
}

func TestSrvWatchDirect(t *testing.T) {
	mock.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	director := mock.NewSrvDirector(ctx)
	out := director.Outlet()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		err := director.Start()
		if err != nil {
			log.Printf("[Director] Start: %v\n", err)
		}
		cancel()
		wg.Done()
	}()

	ch := make(chan direct.IdGetter, 100)

	wg.Add(1)
	go func() {
		err := director.WatcherDirect(ch)
		if err != nil {
			log.Printf("[watcher director]: %v\n", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-ctx.Done():
				log.Println("[watcher] context canceled")
				break Loop
			case <-t.C:
				resp := &mock.Request{Id: "XXX"}
				ch <- resp
			}
		}

		wg.Done()
	}()

Loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("[receive]: for-select: context canceled")
			break Loop
		case req := <-out:
			log.Printf(
				"received request: id: %s | method: %s | issuer: %s",
				req.GetId(), req.GetMethod(), req.GetIssuer(),
			)
		}
	}

	wg.Wait()
}
