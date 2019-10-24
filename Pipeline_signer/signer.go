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

func hashing(data, hash string, counter int, res map[int]map[string]string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	rs := DataSignerCrc32(data)
	mu.Lock()
	res[counter][hash] = rs
	mu.Unlock()

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
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{} // инициализируем группу
	counter := 0
	var resTemp = make(map[int]map[string]string) //
	for dataRaw := range in {
		counter++

		mu.Lock()
		resTemp[counter] = make(map[string]string)
		mu.Unlock()

		dataTmp := dataRaw.(int)
		data := strconv.Itoa(dataTmp)
		Md5 := DataSignerMd5(data)
		wg.Add(2) // добавляем воркер

		go hashing(Md5, "Md5", counter, resTemp, mu, wg)
		go hashing(data, "crc32", counter, resTemp, mu, wg)

	}

	wg.Wait()
	for _, tmp := range resTemp {
		res := tmp["crc32"] + "~" + tmp["Md5"]
		fmt.Println(res)
		out <- res
	}
}

type myMap map[int]string

//MultiHash crc32(th+data)), где data - резульат SingleHash
func MultiHash(in, out chan interface{}) {
	cur := ""
	counter := 0
	var temp = make(map[int]map[int]string)

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{} // инициализируем группу

	for dataRaw := range in {
		data := dataRaw.(string)
		counter++

		mu.Lock()
		temp[counter] = make(map[int]string)
		mu.Unlock()

		for th := 0; th < 6; th++ {
			wg.Add(1) // добавляем воркер
			go func(th, counter int, data, cur string, temp map[int]map[int]string, mu *sync.Mutex, wg *sync.WaitGroup) {
				cur = DataSignerCrc32(strconv.Itoa(th) + data)
				fmt.Printf("%v "+cur+"\n", th)
				mu.Lock()
				temp[counter][th] = cur
				mu.Unlock()
				wg.Done() // уменьшаем счетчик на 1
			}(th, counter, data, cur, temp, mu, wg)

		}

	}
	wg.Wait()
	for _, tmp := range temp {
		out <- (tmp[0] + tmp[1] + tmp[2] + tmp[3] + tmp[4] + tmp[5])
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

	fmt.Printf("%v", str)
	out <- str
}

// Тестирование
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
