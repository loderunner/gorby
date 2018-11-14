package main

type Subscriber interface {
	Subscribe() <-chan RequestResponse
}
