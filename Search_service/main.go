package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// код писать тут

func SearchServer(body string, w *http.ResponseWriter) (string, int) {
	file, err := os.Open("./dataset.xml")
	if err != nil {
		panic(err)
	}

	input := bufio.NewReader(file)
	decoder := xml.NewDecoder(input)
	var name string
	var id int
	for {
		tok, tokenError := decoder.Token()
		if tokenError != nil && tokenError != io.EOF {
			fmt.Println("error happend", tokenError)
			break
		} else if tokenError == io.EOF {
			break
		}
		if tok == nil {
			fmt.Println("t is nil break")
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == "id" {
				if err := decoder.DecodeElement(&id, &tok); err != nil {
					fmt.Println("error happend", err)
				}
			}
			if tok.Name.Local == "first_name" {
				if err := decoder.DecodeElement(&name, &tok); err != nil {
					fmt.Println("error happend", err)
				} else if name == body {
					fmt.Println("Name: ", name)
					fmt.Fprintf(*w, "Id: %v, Name %s\n", id, name)
					return name, id
				}
			}
		}

	}
	return "", 0
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "getHandler: incoming request %#v\n", r)
		//fmt.Fprintf(w, "getHandler: r.Url %#v\n", r.URL)
		myParam := r.URL.Query().Get("param") //http://127.0.0.1:8081/?param=Glenn
		if myParam != "" {
			fmt.Fprintln(w, "‘myParam‘ is", myParam)
		}
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close() // важный пункт!
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		name, id := SearchServer(myParam, &w) //Запускаем поисковик
		fmt.Fprintf(w, "postHandler: raw body %s\n", string(body))
		fmt.Fprintf(w, "Id: %v, Name %v\n", id, name)

	})

	fmt.Println("starting server at :8081")
	http.ListenAndServe(":8081", nil)
}
