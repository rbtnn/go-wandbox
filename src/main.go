package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const method = "POST"
const list_url = "http://melpon.org/wandbox/api/list.json"
const compile_url = "http://melpon.org/wandbox/api/compile.json"

type WandboxInput struct {
	Code     string `json:"code"`
	Compiler string `json:"compiler"`
}

type WandboxOutputCompile struct {
	ProgramError   string `json:"program_error"`
	ProgramMessage string `json:"program_message"`
	Status         string `json:"status"`
}

type WandboxOutputList struct {
	Name string `json:"name"`
}

func execute_compile(in WandboxInput) {
	bytes, err := json.Marshal(in)
	if err != nil {
		return
	}
	reader := strings.NewReader(string(bytes))
	req, _ := http.NewRequest(method, compile_url, reader)
	resp, _ := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	var out WandboxOutputCompile
	json.Unmarshal(bs, &out)
	fmt.Println(out.ProgramMessage)
}

func execute_list() {
	req, _ := http.NewRequest("GET", list_url, nil)
	resp, _ := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	var out []WandboxOutputList
	json.Unmarshal(bs, &out)
	for _, value := range out {
		fmt.Println(value.Name)
	}
}

func main() {
	var f = flag.String("f", "", "source file")
	var c = flag.String("c", "", "compiler")
	var l = flag.Bool("list", false, "compiler list")
	flag.Parse()

	if *l {
		execute_list()
	} else {
		xs, err := ioutil.ReadFile(*f)
		if err != nil {
			return
		}
		in := WandboxInput{
			string(xs),
			*c,
		}
		execute_compile(in)
	}

}
