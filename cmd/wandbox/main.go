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
	ext := filepath.Ext(source)
	name := ""
	switch ext {
	case ".c":
		name = "C"
	case ".C", ".cc", ".cpp", ".cxx":
		name = "C++"
	case ".cs":
		name = "C#"
	case ".coffee":
		name = "CoffeeScript"
	case ".d":
		name = "D"
	case ".erl":
		name = "Erlang"
	case ".ex":
		name = "Elixir"
	case ".go":
		name = "Go"
	case ".groovy":
		name = "Groovy"
	case ".hs":
		name = "Haskell"
	case ".java":
		name = "Java"
	case ".js":
		name = "JavaScript"
	case ".lazyk":
		name = "Lazy K"
	case ".l", ".lsp", ".lisp":
		name = "Lisp"
	case ".lua":
		name = "Lua"
	case ".ml":
		name = "OCaml"
	case ".pas":
		name = "Pascal"
	case ".php":
		name = "PHP"
	case ".pl":
		name = "Perl"
	case ".py":
		name = "Python"
	case ".rb":
		name = "Ruby"
	case ".rs":
		name = "Rust"
	case ".scala":
		name = "Scala"
	case ".sh", ".bash":
		name = "Bash script"
	case ".sqlite3", ".sqlite", ".sql":
		name = "SQL"
	case ".swift":
		name = "Swift"
	case ".vim":
		name = "Vim script"
	}
	m, err := getList()
	if err != nil {
		return "", err
	}
	if value, ok := m[name]; ok {
		return value[0].Name, nil
	} else {
		return "", errors.New(ext + " not supported")
	}
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
	m, _ := getList()
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
