package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var result string = ""

func worker(j job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done() // уменьшаем счетчик на 1
	defer close(out)
	j(in, out)

}

//ExecutePipeline SingleHash -> Multihash -> CombineResults
func ExecutePipeline(f ...job) {
	wg := &sync.WaitGroup{} // инициализируем группу
	wg.Add(1)               // добавляем воркер
	in := make(chan interface{}, 50)
	out := make(chan interface{}, 50)
	for _, val := range f {
		go worker(val, in, out, wg)
	}
	time.Sleep(time.Millisecond)
	wg.Wait() // ожидаем, пока waiter.Done() не приведёт счетчик к 0
}

//SingleHash crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	in <- out
	data := <-in
	fmt.Println(DataSignerCrc32(data.(string)) + "~" + DataSignerCrc32(DataSignerMd5(data.(string))))
	out <- DataSignerCrc32(data.(string)) + "~" + DataSignerCrc32(DataSignerMd5(data.(string)))
}

//MultiHash crc32(th+data)), где data - резульат SingleHash
func MultiHash(in, out chan interface{}) {
	var temp string
	var cur string

	in <- out
	data := <-in

	for th := 0; th < 6; th++ {
		cur = DataSignerCrc32(strconv.Itoa(th) + data.(string))
		fmt.Printf("%v "+cur+"\n", th)
		temp += cur
	}
	fmt.Println(temp)
	out <- temp
}

//CombineResults res1_res2_.._resN
func CombineResults(in, out chan interface{}) {
	in <- out
	data := <-in

	result += data.(string)
}

func main() {
	Hash := []job{
		job(func(in, out chan interface{}) {
			in <- uint32(0)
			in <- uint32(1)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}

	ExecutePipeline(Hash...)

	fmt.Println(result)

	//fmt.Println(SingleHash(os.Args[1]))
	//fmt.Println(MultiHash(SingleHash(os.Args[1])))
}

// сюда писать код
