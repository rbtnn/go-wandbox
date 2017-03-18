package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

const (
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

func executeCompile(data, compiler string) error {
	in := &WandboxInput{
		data,
		compiler,
	}
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", compileURL, buf)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var out WandboxOutputCompile
	err = json.NewDecoder(resp.Body).Decode(&out)
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
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	var out []WandboxOutputList
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return err
	}
	m := map[string][]WandboxOutputList{}
	for _, x := range out {
		m[x.Language] = append(m[x.Language], x)
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		compilers := m[key]
		for i, x := range compilers {
			if i == 0 {
				fmt.Println("[" + key + "]")
			}
			fmt.Println("  " + x.Name)
		}
	}
	return nil
}

func run() int {
	var source, code, compiler string
	var list bool

	flag.StringVar(&source, "source", "", "source file")
	flag.StringVar(&code, "code", "", "code")
	flag.StringVar(&compiler, "compiler", "", "compiler")
	flag.BoolVar(&list, "list", false, "compiler list")
	flag.Parse()

	if list {
		err := executeList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: list error: %v\n", os.Args[0], err)
			return 1
		}
	} else {
		data := ""
		if code != "" {
			data = code
		} else if source != "" {
			xs, err := ioutil.ReadFile(source)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: read file: %v\n", os.Args[0], err)
				return 1
			}
			data = string(xs)
		} else {
			flag.Usage()
			return 0
		}
		err := executeCompile(data, compiler)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: compile error: %v\n", os.Args[0], err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
