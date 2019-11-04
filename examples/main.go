package main

import (
	"fmt"
	"time"

	"github.com/andreiavrammsd/progress"
)

type service struct {
}

func (service) Error(err error) {
	fmt.Println(err.Error())
}

func main() {
	service := &service{}
	c := progress.Config{
		Drawer: &progress.Characters{},
		Error:  service.Error,
	}
	p := progress.New(&c)

	time.AfterFunc(time.Second*2, func() {
		p.Stop()
	})

	// Process making progress
	for j := 0; j <= 1; j++ {
		go func(p *progress.Progress, x int) {
			i := 0
			for {
				p.Progress()
				fmt.Println(x, i)

				i++
				if i > 15 {
					time.Sleep(time.Millisecond * 30)
					continue
				}
				if i > 10 {
					time.Sleep(time.Millisecond * 300)
					continue
				}
				if i > 5 {
					time.Sleep(time.Millisecond * 500)
					continue
				}
				time.Sleep(time.Millisecond * 150)
			}
		}(p, j)
	}

	fmt.Println("here")

	time.Sleep(time.Hour)
}
