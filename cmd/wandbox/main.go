package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	method     = "POST"
	listURL    = "http://melpon.org/wandbox/api/list.json"
	compileURL = "http://melpon.org/wandbox/api/compile.json"
)

// WandboxInput mean request structure to the compile.
type WandboxInput struct {
	Code     string `json:"code"`
	Compiler string `json:"compiler"`
}

// WandboxOutputCompile mean response structure from the compile.
type WandboxOutputCompile struct {
	ProgramError   string `json:"program_error"`
	ProgramMessage string `json:"program_message"`
	Status         string `json:"status"`
}

// WandboxOutputList mean response structure from the list.
type WandboxOutputList struct {
	Name                  string `json:"name"`
	Language              string `json:"language"`
	DisplayCompileCommand string `json:"display-compile-command"`
}

func executeCompile(in *WandboxInput) error {
	bytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(bytes))
	req, err := http.NewRequest(method, compileURL, reader)
	if err != nil {
		return err
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var out WandboxOutputCompile
	err = json.Unmarshal(bs, &out)
	if err != nil {
		return err
	}
	fmt.Println(out.ProgramMessage)
	return nil
}

func executeList() error {
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return err
	}
	resp, _ := (&http.Client{}).Do(req)
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var out []WandboxOutputList
	err = json.Unmarshal(bs, &out)
	if err != nil {
		return err
	}
	m := map[string][]WandboxOutputList{}
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
	return nil
}

func main() {
	var source = flag.String("source", "", "source file")
	var code = flag.String("code", "", "code")
	var compiler = flag.String("compiler", "", "compiler")
	var list = flag.Bool("list", false, "compiler list")
	flag.Parse()

	if *list {
		executeList()
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
		executeCompile(&WandboxInput{
			data,
			*compiler,
		})
	}

}
