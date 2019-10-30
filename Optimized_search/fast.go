package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

//easyjson:json
type JSONData struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

var dataPool = sync.Pool{
	New: func() interface{} {
		return &JSONData{}
	},
}

//FastSearch Вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	seenBrowsers := make(map[string]bool, 300)
	uniqueBrowsers := 0

	//Variant 1
	users := make([]JSONData, 0, 1000)
	rLine, _, err := reader.ReadLine()
	for err == nil {
		user := &JSONData{}
		// fmt.Printf("%v %v\n", err, line)
		err := user.UnmarshalJSON(rLine)
		if err != nil {
			panic(err)
		}
		users = append(users, *user)
		rLine, _, err = reader.ReadLine()
		if err != nil {
			break
		}
	}

	//Variant 2
	// users := make(chan JSONData, 100)
	// go func(out chan<- JSONData) {
	// 	for err == nil {
	// 		user := &JSONData{}
	// 		// fmt.Printf("%v %v\n", err, line)
	// 		err := user.UnmarshalJSON(rLine)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		users <- (*user)
	// 		rLine, _, err = reader.ReadLine()
	// 		if err != nil {
	// 			break
	// 		}
	// 	}
	// 	close(out)
	// }(users)

	fmt.Fprintln(out, "found users:")
	i := -1
	for _, user := range users {
		i++
		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {

			notSeenBefore := true

			Android := strings.Contains(browser, "Android")
			MSIE := strings.Contains(browser, "MSIE")

			if (Android || MSIE) == true {

				if Android {
					isAndroid = true
				}

				if MSIE {
					isMSIE = true
				}

				if exist := seenBrowsers[browser]; exist == true {
					notSeenBefore = false
				}

				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}

			}

		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", 1)

		fmt.Fprintln(out, "["+strconv.Itoa(i)+"] "+user.Name+" <"+email+">")
		//temp = append(temp, "["+strconv.Itoa(i)+"] "+user.Name+" <"+email+">")
	}
	//fmt.Fprintln(out, "found users:\n"+buf.String())
	//fmt.Fprintln(out, "found users:\n"+strings.Join(temp, "\n")+"\n")
	fmt.Fprintln(out, "\n"+"Total unique browsers", len(seenBrowsers))
}
