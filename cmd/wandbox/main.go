package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

func detectLatestCompiler(source string) (string, error) {
	var ext2name = map[string]string{
		".c":       "C",
		".C":       "C++",
		".cc":      "C++",
		".cpp":     "C++",
		".cxx":     "C++",
		".cs":      "C#",
		".coffee":  "CoffeeScript",
		".d":       "D",
		".erl":     "Erlang",
		".ex":      "Elixir",
		".go":      "Go",
		".groovy":  "Groovy",
		".hs":      "Haskell",
		".java":    "Java",
		".js":      "JavaScript",
		".lazyk":   "Lazy K",
		".l":       "Lisp",
		".lsp":     "Lisp",
		".lisp":    "Lisp",
		".lua":     "Lua",
		".ml":      "OCaml",
		".pas":     "Pascal",
		".php":     "PHP",
		".pl":      "Perl",
		".py":      "Python",
		".rb":      "Ruby",
		".rs":      "Rust",
		".scala":   "Scala",
		".sh":      "Bash script",
		".bash":    "Bash script",
		".sqlite3": "SQL",
		".sqlite":  "SQL",
		".sql":     "SQL",
		".swift":   "Swift",
		".vim":     "Vim script",
	}
	ext := filepath.Ext(source)
	if name, ok := ext2name[ext]; ok {
		m, err := getList()
		if err != nil {
			return "", err
		}
		if compiler, ok := m[name]; ok {
			return compiler[0].Name, nil
		}
	}
	return "", errors.New(ext + " not supported")
}

func executeCompile(data, compiler string) error {
	in := &WandboxInput{
		data,
		compiler,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", compileURL, &buf)
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

func getList() (map[string][]WandboxOutputList, error) {
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out []WandboxOutputList
	// bs, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bs))
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return nil, err
	}
	m := map[string][]WandboxOutputList{}
	for _, x := range out {
		m[x.Language] = append(m[x.Language], x)
	}
	return m, nil
}

func executeList() error {
	m, err := getList()
	if err != nil {
		return err
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
			if compiler == "" {
				compiler, err = detectLatestCompiler(source)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: detect latest compiler: %v\n", os.Args[0], err)
					return 1
				}
			}
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
