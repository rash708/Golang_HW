package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// type TestCase struct {
// 	Request *SearchRequest
// 	Result  *SearchResponse
// 	IsError bool
// }

// type Databse struct {
// }

// type UserBase struct {
// 	ID        int    `xml: "id"`
// 	Firstname string `xml: "first_name"`
// 	Lastname  string `xml: "last_name"`
// 	About     string `xml: "about"`
// 	Age       int    `xml: "age"`
// }

// код писать тут

func SearchServer(body []byte) {

	// file, err := os.Open("./dataset.xml")
	// if err != nil {
	// 	panic(err)
	// }

	// input := bufio.NewReader(file)
	// decoder := xml.NewDecoder(input)
	// var name string
	// for {
	// 	tok, tokenError := decoder.Token()
	// 	if tokenError != nil && tokenError != io.EOF {
	// 		fmt.Println("error happend", tokenError)
	// 		break
	// 	} else if tokenError == io.EOF {
	// 		break
	// 	}
	// 	if tok == nil {
	// 		fmt.Println("t is nil break")
	// 	}
	// 	switch tok := tok.(type) {
	// 	case xml.StartElement:
	// 		if tok.Name.Local == "first_name" {
	// 			if err := decoder.DecodeElement(&name, &tok); err != nil {
	// 				fmt.Println("error happend", err)
	// 			} else {
	// 				fmt.Println("Name: ", name)
	// 			}
	// 		}
	// 	}

	// }

}

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "getHandler: incoming request %#v\n", r)
	// 	fmt.Fprintf(w, "getHandler: r.Url %#v\n", r.URL)
	// 	body, err := ioutil.ReadAll(r.Body)
	// 	defer r.Body.Close() // важный пункт!
	// 	if err != nil {
	// 		http.Error(w, err.Error(), 500)
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "postHandler: raw body %s\n", string(body))
	// })

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "getHandler: incoming request %#v\n", r)
		fmt.Fprintf(w, "getHandler: r.Url %#v\n", r.URL)
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close() // важный пункт!
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		SearchServer(body) //Запускаем поисковик
		fmt.Fprintf(w, "postHandler: raw body %s\n", string(body))
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// func FindUsersDummy(w http.ResponseWriter, r *http.Request) {

// }

// func SearchServer(t *testing.T) {
// 	cases := []TestCase{
// 		TestCase{
// 			Request: &SearchRequest,
// 			Result: &SearchResponse{
// 				Users:    []User{},
// 				NextPage: true,
// 			},
// 			IsError: false,
// 		},
// 		TestCase{
// 			Request: "100500",
// 			Result: &SearchResponse{
// 				Users:    []User{},
// 				NextPage: true,
// 			},
// 			IsError: false,
// 		},
// 		TestCase{
// 			Request:      "__broken_json",
// 			Result:  nil,
// 			IsError: true,
// 		},
// 		TestCase{
// 			Request:      "__internal_error",
// 			Result:  nil,
// 			IsError: true,
// 		},
// 	}

// 	ts := httptest.NewServer(http.HandlerFunc(FindUsersDummy))

// 	for caseNum, item := cases {
// 		sc := &SearchClient{
// 			URL: ts.URL,
// 			AccessToken: ts.URL,
// 		}
// 	}

// 	ts.Close()
// }
