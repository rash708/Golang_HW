package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func parse(out io.Writer, tabPrev, path string, printFiles bool) {
	var temp []os.FileInfo
	//Просматриваем указанную дирректорию
	files, err := ioutil.ReadDir(path)
	//Обозначаем отступы
	tab := `│	`
	tabLast := `	`

	if err != nil {
		log.Fatal(err)
	}

	//Формируем слайс включая/не включая файлы

	for _, t := range files {

		if t.IsDir() || printFiles {
			temp = append(temp, t)
			//fmt.Printf("%v\n", t.Name())
		}
	}

	//Обход всех папок и файлов
	for _, f := range temp {
		last := temp[len(temp)-1]

		//Если это дирректория то вывести имя и продолжить обход,
		//иначе вывести только имя и размер
		if f.IsDir() {

			if f != last {
				fmt.Fprintf(out, tabPrev+"├───%v\n", f.Name())
				//Продолжаем обход папок
				parse(out, tabPrev+tab, path+f.Name()+"/", printFiles)
			} else {
				fmt.Fprintf(out, tabPrev+"└───%v\n", f.Name())
				//Продолжаем обход папок
				parse(out, tabPrev+tabLast, path+f.Name()+"/", printFiles)
			}

		} else {
			if f != last {

				if f.Size() == 0 {
					fmt.Fprintf(out, tabPrev+"├───%v (empty)\n", f.Name())
				} else {
					fmt.Fprintf(out, tabPrev+"├───%v (%vb)\n", f.Name(), f.Size())
				}

			} else {

				if f.Size() == 0 {
					fmt.Fprintf(out, tabPrev+"└───%v (empty)\n", f.Name())
				} else {
					fmt.Fprintf(out, tabPrev+"└───%v (%vb)\n", f.Name(), f.Size())
				}

			}

		}

	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	path = path + "/"
	parse(out, "", path, printFiles)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
