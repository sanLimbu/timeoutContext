package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	userID := 1
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}
	for _, va := range val {
		fmt.Println(va)
	}
	fmt.Println("took :", time.Since(start))
}

type Response struct {
	value string
	err   error
}

func fetchUserData(ctx context.Context, userID int) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()
	respch := make(chan Response, 100)
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go fetchThirdPartyApi(respch, wg)
	go fetchThirdPartyApi2(respch, wg)
	go fetchThirdPartyApi3(respch, wg)

	go func() {
		wg.Wait()
		close(respch)
	}()
	var responses []string
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			// Do not send response if context is cancelled
		case resp := <-respch:
			responses = append(responses, resp.value)
			// return resp.value, resp.err
		}
	}
	return responses, nil
}

func fetchThirdPartyApi(respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 400)
	val := "response value 1"
	respch <- Response{
		value: val,
		err:   nil,
	}
	wg.Done()
}
func fetchThirdPartyApi2(respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 5000)
	val := "response value 2"
	respch <- Response{
		value: val,
		err:   nil,
	}
	wg.Done()
}

func fetchThirdPartyApi3(respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 500)
	val := "response value 3"
	respch <- Response{
		value: val,
		err:   nil,
	}
	wg.Done()
}
