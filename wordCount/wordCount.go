package wordcount

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	workerpool "../workerPool"
)

//FileChan struct
type FileChan struct {
	FileName string
	Channel  chan map[string]int
}

//RunJob implement for worker job
func (f FileChan) RunJob() {
	countedMap := wordCount(f.FileName)
	f.Channel <- countedMap
}

// ListAllTxt return all *.txt in folder and sub-folder
func ListAllTxt(pathToFolder string) []string {
	var files []string

	err := filepath.Walk(pathToFolder, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".txt") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// readFile return the context of file
func readFile(pathToFile string) string {
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		fmt.Println(pathToFile, "File reading error", err)
		return ""
	}
	return string(data)
}

// wordCount count exit words in the file
func wordCount(pathToFile string) map[string]int {
	context := readFile(pathToFile)
	listWords := strings.Split(context, " ")
	countTable := make(map[string]int)
	for _, word := range listWords {
		countTable[word]++
	}
	return countTable
}

//Combine all result return from channel then stop worker pool
func Combine(channel chan map[string]int, cnt int, p *workerpool.Pool, resultchannel chan map[string]int) {
	result := make(map[string]int)
	for i := 0; i < cnt; i++ {
		tmp := <-channel
		for key, value := range tmp {
			result[key] += value
		}
	}
	p.Destroy()
	resultchannel <- result
}

//PrintCountTable print easy looking format
func PrintCountTable(table map[string]int) {
	for key, value := range table {
		fmt.Println(key, " = ", value)
	}
}
