package operate_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	direct "github.com/murtaza-u/ddos/mock/direct"
	store "github.com/murtaza-u/ddos/mock/store"
	"github.com/murtaza-u/ddos/operate"
	pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"
)

func TestGet(t *testing.T) {
	defer direct.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := store.NewStore()
	director := direct.NewSrvDirector(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		err := director.Start()
		if err != nil {
			log.Printf("director.Start: %v\n", err)
		}
		cancel()
		wg.Done()
	}()

	opt := operate.NewDefaultOperator(ctx, s, "resource", "id")
	err := opt.Get(director)
	if err != nil {
		t.Fatal(err)
	}

	cancel()
	wg.Wait()
}

func TestPut(t *testing.T) {
	defer direct.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := store.NewStore()
	director := direct.NewSrvDirector(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		err := director.Start()
		if err != nil {
			log.Printf("director.Start: %v\n", err)
		}
		cancel()
		wg.Done()
	}()

	opt := operate.NewDefaultOperator(ctx, s, "resource", "id")
	err := opt.Put(director, &pb.Request{
		Id: "XXX",
		Resource: &pb.Resource{
			Version:  0,
			Manifest: []byte{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	cancel()
	wg.Wait()
}

func TestWatch(t *testing.T) {
	defer direct.Reset()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := store.NewStore()
	director := direct.NewSrvDirector(ctx)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		err := director.Start()
		if err != nil {
			log.Printf("director.Start: %v\n", err)
		}
		cancel()
		wg.Done()
	}()

	opt := operate.NewDefaultOperator(ctx, s, "resource", "id")

	wg.Add(1)
	go func() {
		opt.Watch(director)
		wg.Done()
	}()

	time.Sleep(time.Second)

	cancel()
	wg.Wait()
}
