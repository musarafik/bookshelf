package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var books []BookSchema

type BookSchema struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	var bookSchema BookSchema

	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		log.Println("Invalid method")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error reading body", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &bookSchema)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error parsing body", err)
	}

	books = append(books, bookSchema)

	fmt.Fprintf(w, "%s %s", bookSchema.Title, bookSchema.Author)
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.bookList(books)
	component.Render(context.Background(), w)
}

func main() {
	http.HandleFunc("/books/add", addBookHandler)
	http.HandleFunc("/", getBookHandler)

	fmt.Println("Starting server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting server", err)
	}
}
