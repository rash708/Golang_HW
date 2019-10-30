package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//easyjson:json
type JSONData struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
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
	user := &JSONData{}
	i := -1             //Итератор
	var end error = nil //Флаг конца строки
	var browsers []string

	fmt.Fprintln(out, "found users:")

	rLine, _, end := reader.ReadLine() //Начало чтения файла
	for end == nil {

		err := user.UnmarshalJSON(rLine)
		if err != nil {
			panic(err)
		}

		i++
		isAndroid := false
		isMSIE := false
		browsers = user.Browsers

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

		rLine, _, end = reader.ReadLine() //Переход на следующую строку с проверкой конца строки

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", 1)
		fmt.Fprintln(out, "["+strconv.Itoa(i)+"] "+user.Name+" <"+email+">")
	}

	fmt.Fprintln(out, "\n"+"Total unique browsers", len(seenBrowsers))
}
