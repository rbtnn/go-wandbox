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
	Name                  string `json:"name"`
	Language              string `json:"language"`
	DisplayCompileCommand string `json:"display-compile-command"`
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
	var m map[string][]WandboxOutputList = map[string][]WandboxOutputList{}
	for _, x := range out {
		m[x.Language] = append(m[x.Language], x)
	}
	for key, compilers := range m {
		for i, x := range compilers {
			if i == 0 {
				fmt.Println("[" + key + "]")
			}
			fmt.Println("  " + x.Name)
		}
	}
}

func main() {
	var source = flag.String("source", "", "source file")
	var code = flag.String("code", "", "code")
	var compiler = flag.String("compiler", "", "compiler")
	var list = flag.Bool("list", false, "compiler list")
	flag.Parse()

	if *list {
		execute_list()
	} else {
		data := ""
		if 0 < len(*code) {
			data = *code
		} else {
			xs, err := ioutil.ReadFile(*source)
			if err != nil {
				return
			}
			data = string(xs)
		}
		execute_compile(WandboxInput{
			data,
			*compiler,
		})
	}

}
