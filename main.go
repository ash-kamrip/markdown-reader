package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type FileReader struct{}

// read file content {METHOD}
func (fr FileReader) Read(slug string) (string, error) {

	f, err := os.Open(slug + ".md")
	// loggging the captured slug
	log.Println(f)
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
func PostHandler(fr FileReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		// loggging the captured slug
		//log.Println(slug)
		postMarkdown, err := fr.Read(slug)
		if err != nil {
			http.Error(w, "Post not found !", http.StatusNotFound)
			return
		}
		fmt.Fprint(w, postMarkdown)
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
