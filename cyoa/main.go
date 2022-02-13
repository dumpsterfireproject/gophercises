package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Story map[string]Chapter

var storyTemplateHtml = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
	  {{range .Story}}
        <p>{{.}}</p>
      {{end}}
      {{if .Options}}
        <ul>
        {{range .Options}}
          <li><a href="/{{.Arc}}">{{.Text}}</a></li>
        {{end}}
        </ul>
      {{else}}
        <h3>The End</h3>
      {{end}}
  </body>
</html>
`

var storyTemplate *template.Template

func initTemplate() *template.Template {
	return template.Must(template.New("").Parse(storyTemplateHtml))
}

func main() {
	f := flag.String("file", "gopher.json", "the json file with the CYOA story")
	flag.Parse()
	story, err := loadFromFile(*f)
	if err != nil {
		fmt.Printf("Error parsing story %v", err)
		os.Exit(1)
	}
	storyTemplate = initTemplate()
	mux := defaultMux()
	storyHandler := StoryHandler(story, mux)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", storyHandler)
}

func loadFromFile(fileName string) (Story, error) {
	var story Story
	f, err := os.Open(fileName)
	if err != nil {
		return story, err
	}
	defer f.Close()
	story, err = readStory(f)
	return story, err
}

func readStory(reader io.Reader) (Story, error) {
	var story Story
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func StoryHandler(story Story, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := strings.Trim(request.URL.Path, "/")
		// TODO: better default
		if path == "" {
			path = "intro"
		}
		if chapter, ok := story[path]; ok {
			storyTemplate.Execute(writer, chapter)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
