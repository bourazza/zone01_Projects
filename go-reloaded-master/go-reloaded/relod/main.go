package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"relod"
	"strings"
)

func main() {
	if len(os.Args) != 3 || !(strings.HasSuffix(os.Args[1], ".txt")) || !(strings.HasSuffix(os.Args[2], ".txt") || os.Args[1] == os.Args[2]) {
		log.Println("errr")
		return

	} else {
		var holder string
		text := os.Args[1]
		text2 := os.Args[2]

		data, err := os.ReadFile(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		var g string

		words := strings.Split(string(data), "\n")

		for i := 0; i < len(words); i++ {
			words[i] = strings.ReplaceAll(string(words[i]), "\r", "")
			words[i] = strings.ReplaceAll(string(words[i]), "\n", "")

			holder = relod.TransformText(string(words[i]))
			holder = relod.FormatText(holder)

			holder = relod.Fixit(holder)
			g += holder + "\n"
		}

		file, err := os.Create(text2)
		if err != nil {
			panic(err)
		}
		fdata, err := io.WriteString(file, g)
		if err != nil {
			log.Fatal(err)
		} else if fdata == 0 {
			fmt.Printf("the file is empty")
		}
	}

}
