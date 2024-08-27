package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type FileReader struct{}

type Post struct {
	Content string
	Author  string
	Title   string
}

// read file content {METHOD}
func (fr FileReader) Read(slug string) (string, error) {

	f, err := os.Open(slug + ".md")
	// loggging the captured slug
	//log.Println(f)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fcontent, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fcontent), nil
}

// read file content {FUNCTION}
func PostHandler(fr FileReader, tpl *template.Template) http.HandlerFunc {

	// this doesn't get rendered again and again
	// somehow , i don't understand this as of now.
	mdRenderer := goldmark.New(goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		),
	))

	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		// loggging the captured slug
		//log.Println(slug)
		postMarkdown, err := fr.Read(slug)
		if err != nil {
			http.Error(w, "Post not found !", http.StatusNotFound)
			return
		}
		var buf bytes.Buffer

		if err := mdRenderer.Convert([]byte(postMarkdown), &buf); err != nil {
			panic(err)
		}
		// print out everything from byte buffer to Writer
		// fmt.Fprint(w, buf)
		//io.Copy(w, &buf)

		err = tpl.Execute(w, &Post{
			Content: buf.String(),
			Author:  "Ashu Adhana",
			Title:   " we have only 10 days to live",
		})
		if err != nil {
			http.Error(w, "Error Executing Template", http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	mux := http.NewServeMux()

	postemplate := template.Must(template.ParseFiles("post.gohtml"))

	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}, postemplate))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
