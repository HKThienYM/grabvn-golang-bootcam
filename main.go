package main

import (
	"log"

	wordcount "./wordCount"
	workerpool "./workerPool"
)

const (
	dirPath      = "./assest"
	numOfWWorker = 4
)

func main() {
	listfiles := wordcount.ListAllTxt(dirPath)

	if listfiles == nil {
		log.Fatalln("Can not find any *.txt file in the input folder")
	}

	centerchannel := make(chan map[string]int)
	resultchannel := make(chan map[string]int)
	p := workerpool.NewWorkerPool(numOfWWorker)

	go wordcount.Combine(centerchannel, len(listfiles), p, resultchannel)

	for _, file := range listfiles {
		filechan := wordcount.FileChan{FileName: file,
			Channel: centerchannel}
		p.Dispatch(workerpool.Job(filechan))
	}
	result := <-resultchannel
	wordcount.PrintCountTable(result)
}
