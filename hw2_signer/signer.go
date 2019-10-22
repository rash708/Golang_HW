package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

func worker(j job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done() // уменьшаем счетчик на 1
	defer close(out)
	j(in, out)
}

//ExecutePipeline SingleHash -> Multihash -> CombineResults
func ExecutePipeline(f ...job) error {
	wg := &sync.WaitGroup{} // инициализируем группу
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)
	for _, val := range f {
		wg.Add(1) // добавляем воркер
		go worker(val, in, out, wg)

		//Реализация конвеера ->in:out->in:out->in:out
		in = out
		outNew := make(chan interface{}, 100)
		out = outNew
	}
	wg.Wait() // ожидаем, пока waiter.Done() не приведёт счетчик к 0

	return nil
}

//SingleHash crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	for dataRaw := range in {
		dataTmp, _ := dataRaw.(int)
		data := strconv.Itoa(dataTmp)
		res := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
		fmt.Println(res)
		out <- res
	}
}

//MultiHash crc32(th+data)), где data - резульат SingleHash
func MultiHash(in, out chan interface{}) {
	var temp string
	var cur string

	for dataRaw := range in {
		data, _ := dataRaw.(string)

		temp = ""
		cur = ""

		for th := 0; th < 6; th++ {
			cur = DataSignerCrc32(strconv.Itoa(th) + data)
			fmt.Printf("%v "+cur+"\n", th)
			temp += cur
		}
		fmt.Println(temp)
		out <- temp
	}
}

//CombineResults res1_res2_.._resN
func CombineResults(in, out chan interface{}) {
	var arr []string
	var str string

	for dataRaw := range in {
		data, _ := dataRaw.(string)
		arr = append(arr, data)
	}

	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	for pos, temp := range arr {
		if pos == 0 {
			str += temp
		} else {
			str += "_" + temp
		}
	}

	//fmt.Printf("%v \n", arr)
	fmt.Printf("%v", str)
	out <- str
}

//Тестирование
func main() {
	inputData := []int{0, 1}
	Hash := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}

	ExecutePipeline(Hash...)
}
