package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("You can try (1 +2)/3 + 4*(1-1)\nInput \"q\" to exit\n> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "q" {
			break
		}
		if text == "" {
			continue
		}
		result, err := calculateString(text)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%v = %v", text, result)
		}
		fmt.Print("\n> ")
	}

}
