package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/andreiko/alfred-sources/sources"
)

type responseItemType struct {
	Autocomplete string                 `json:"autocomplete"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type responseType struct {
	Items []responseItemType `json:"items"`
}

type SourceServer struct {
	sources    map[string]sources.Source
	urlPattern *regexp.Regexp
}

func (srv *SourceServer) AddSource(src sources.Source) {
	srv.sources[src.Id()] = src
}

func (srv *SourceServer) Start() {
	server := &http.Server{
		Handler: srv,
		Addr:    "127.0.0.1:8080",
	}

	server.ListenAndServe()
}

func (srv *SourceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	url_submatches := srv.urlPattern.FindStringSubmatch(req.URL.Path)
	if len(url_submatches) < 2 {
		resp.WriteHeader(404)
		fmt.Fprint(resp, "Not found")
		return
	}

	source_id := url_submatches[1]
	src, ok := srv.sources[source_id]
	if !ok {
		resp.WriteHeader(404)
		fmt.Fprintf(resp, "No such source: %s", source_id)
		return
	}

	qs := req.URL.Query()
	queryList, ok := qs["query"]
	if !ok {
		resp.WriteHeader(400)
		fmt.Fprint(resp, "Parameter required: query")
		return
	}

	if len(queryList) != 1 {
		resp.WriteHeader(400)
		fmt.Fprint(resp, "Parameter requires exactly one value: query")
		return
	}

	query := queryList[0]

	response := &responseType{
		Items: make([]responseItemType, 0),
	}
	for _, item := range src.Query(query) {
		response.Items = append(response.Items, responseItemType{
			Autocomplete: item.Autocomplete(),
			Attributes:   item.Attributes(),
		})
	}

	if resp_stream, err := json.Marshal(response); err == nil {
		resp.Header()["Content-Type"] = []string{"application/json"}
		resp.Write(resp_stream)
	} else {
		resp.WriteHeader(500)
		fmt.Fprintf(resp, "Error encoding response: %s", err.Error())
	}
}

func NewSourceServer() *SourceServer {
	urlPattern, err := regexp.Compile("^/sources/(\\w+)/?$")
	if err != nil {
		panic(err.Error())
	}

	return &SourceServer{
		sources:    make(map[string]sources.Source),
		urlPattern: urlPattern,
	}
}
