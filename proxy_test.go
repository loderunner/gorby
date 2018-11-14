package main

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestSubscribeAndDispatch(t *testing.T) {
	expect := RequestResponse{newTestRequest(time.Now()), newTestResponse(time.Now())}
	p := NewProxyHandler()
	c := p.Subscribe()
	go p.dispatch(expect.Request, expect.Response)

	got := <-c

	if !reflect.DeepEqual(expect, got) {
		t.Fatalf("got %#v, expected %#v", got, expect)
	}
}

func TestDispatchTimeout(t *testing.T) {
	expect := RequestResponse{newTestRequest(time.Now()), newTestResponse(time.Now())}
	p := NewProxyHandler()
	c := p.Subscribe()
	p.dispatch(expect.Request, expect.Response)
	p.dispatch(expect.Request, expect.Response)

	got := <-c

	if !reflect.DeepEqual(expect, got) {
		t.Fatalf("got %#v, expected %#v", got, expect)
	}

	got, ok := <-c
	if ok {
		t.Fatalf("unexpected success reading from channel")
	}
}

func TestConcurrency(t *testing.T) {
	p := NewProxyHandler()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		p.dispatch(nil, nil)
		p.dispatch(nil, nil)
		wg.Done()
	}()
	go func() {
		p.Subscribe()
		p.dispatch(nil, nil)
		wg.Done()
	}()
	wg.Wait()
}
