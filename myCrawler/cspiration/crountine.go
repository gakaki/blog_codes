package cspiration

import "fmt"

type Job struct {
	name string
}

func process(job Job) {
	fmt.Println(job.name)
}

func worker(jobChan <-chan Job) {
	for job := range jobChan {
		process(job)
	}
}

//
//func main(){
//	fmt.Println("hello crountine")
//
//	// make a channel with a capacity of 100.
//	jobChan := make(chan Job, 100)
//
//	// start the worker
//	for {
//		go worker(jobChan)
//	}
//	for i := 1; i <= 1000; i++ {
//		jobChan <- Job{ fmt.Sprintf("%d",i) }
//	}
//
//}
