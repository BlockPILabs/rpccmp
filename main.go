package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var fileName string

func init() {
	flag.StringVar(&fileName, "f", "cmp.json", "")
	flag.Parse()
}

type Cmp struct {
	Tag       string `json:"tag"`
	Endpoints []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"endpoints"`
	Data []string `json:"data"`
}

func main() {
	if len(fileName) == 0 {
		return
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		println(err)
		return
	}

	var cmp *Cmp
	err = json.Unmarshal(content, &cmp)
	if err != nil {
		println(err)
		return
	}

	for _, data := range cmp.Data {
		isOk := true
		var md5Hex string
		println(cmp.Tag + ":\n" + data)
		for _, endpoint := range cmp.Endpoints {
			var dataSorted []byte
			resp, err := http.Post(endpoint.Url, "application/json", bytes.NewBuffer([]byte(data)))
			if err != nil {
				println(endpoint.Name + "\t" + err.Error())
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			var tempData map[string]any
			json.Unmarshal(body, &tempData)
			dataSorted, _ = json.Marshal(tempData)
			md5Val := md5.Sum(dataSorted)
			hexVal := hex.EncodeToString(md5Val[0:])

			if len(md5Hex) == 0 {
				md5Hex = hexVal
			}

			if hexVal != md5Hex {
				isOk = false
			}

			println(endpoint.Name + "\t" + md5Hex + "  " + string(dataSorted))
		}

		if isOk {
			fmt.Printf("\033[1;32m%s\033[0m\n", "Cmp ok.\n\n")
		} else {
			fmt.Printf("\033[1;32m%s\033[0m\n", "Cmp failed.\n\n")
		}
	}
}
