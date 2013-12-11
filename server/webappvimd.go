package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keep94/weblogs"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

type response struct {
	Header []string `json:"header"`
	Body   string   `json:"body"`
	Status int      `json:"status"`
}

type request struct {
	Header []string `json:"header"`
	Path   string   `json:"path"`
	Query  string   `json:"query"`
	Body   string   `json:"body"`
	Method string   `json:"method"`
}

var re_header = regexp.MustCompile(`^([a-zA-Z][-_a-zA-Z0-9]*):\s*(.*)`)

func main() {
	b, err := exec.Command("vim", "--serverlist").Output()
	if err != nil {
		log.Fatal(err)
	}
	vim := ""
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		b, err = exec.Command("vim", "--servername", line, "--remote-expr", `string(function('webapp#serve'))`).Output()
		if err == nil && strings.TrimSpace(string(b)) == `function('webapp#serve')` {
			vim = line
			break
		}

	}
	if vim == "" {
		log.Fatal("vim doesn't support remote protocol")
	}
	log.Print("Registered vim server: ", vim)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var buf bytes.Buffer
		err = r.Header.Write(&buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req := &request{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Header: strings.Split(strings.TrimSpace(buf.String()), "\r\n"),
			Body:   string(body),
		}
		b, err = json.Marshal(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		input := string(b)
		input = strings.Replace(input, `'`, `\x27`, -1)
		payload := fmt.Sprintf(`webapi#json#encode(webapp#serve(webapi#json#decode('%s')))`, input)
		b, err = exec.Command("vim", "--servername", vim, "--remote-expr", payload).CombinedOutput()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(b) > 0 && b[0] == 'E' {
			http.Error(w, string(b), http.StatusInternalServerError)
			return
		}
		var res response
		err = json.Unmarshal(b, &res)
		if err != nil {
			http.Error(w, err.Error()+": "+string(b), http.StatusInternalServerError)
			return
		}
		for _, header := range res.Header {
			kv := re_header.FindStringSubmatch(header)
			if len(kv) == 3 {
				w.Header().Set(kv[1], kv[2])
			}
		}
		if res.Status != 0 {
			w.WriteHeader(res.Status)
		}
		w.Write([]byte(res.Body))
	})

	err = http.ListenAndServe("127.0.0.1:9001", weblogs.Handler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
