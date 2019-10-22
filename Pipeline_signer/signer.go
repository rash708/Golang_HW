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

// func single(j func(string) string, in, out chan string) {
// 	data := <-in
// 	out <- j(data)
// }

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

	//type hash func(func(data string) string, j job)

	// Conv := []hash{
	// 	job(SingleHash),
	// 	job(MultiHash),
	// 	job(CombineResults),
	// }

	//mu := &sync.Mutex{}
	for dataRaw := range in {
		var resTemp = make(map[string]string)
		dataTmp, _ := dataRaw.(int)
		data := strconv.Itoa(dataTmp)

		// in := make(chan string, 10)
		// out := make(chan string, 10)
		// in <- data

		// go single(DataSignerMd5, in, out)
		// go single(DataSignerCrc32, out, in)

		mu := &sync.Mutex{}
		wg := &sync.WaitGroup{} // инициализируем группу

		wg.Add(2) // добавляем воркер

		go func(res map[string]string, mu *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			Md5 := DataSignerCrc32(DataSignerMd5(data))
			mu.Lock()
			res["Md5"] = Md5
			mu.Unlock()
		}(resTemp, mu, wg)

		go func(res map[string]string, mu *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			crc32 := DataSignerCrc32(data)
			//crc32Md5 := DataSignerCrc32(res["Md5"])
			mu.Lock()
			res["crc32"] = crc32
			//res["Md5"] = crc32Md5
			mu.Unlock()
		}(resTemp, mu, wg)

		wg.Wait()

		res := resTemp["crc32"] + "~" + resTemp["Md5"]
		fmt.Println(res)
		out <- res

	}
}

//MultiHash crc32(th+data)), где data - резульат SingleHash
func MultiHash(in, out chan interface{}) {
	for dataRaw := range in {
		data, _ := dataRaw.(string)

		var temp = make(map[int]string)
		cur := ""

		mu := &sync.Mutex{}
		wg := &sync.WaitGroup{} // инициализируем группу

		//wg.Add(1) // добавляем воркер

		for th := 0; th < 6; th++ {
			wg.Add(1) // добавляем воркер
			go func(th int, data, cur string, temp map[int]string, wg *sync.WaitGroup, mu *sync.Mutex) {
				cur = DataSignerCrc32(strconv.Itoa(th) + data)
				fmt.Printf("%v "+cur+"\n", th)
				mu.Lock()
				temp[th] = cur
				mu.Unlock()
				wg.Done() // уменьшаем счетчик на 1
			}(th, data, cur, temp, wg, mu)

		}

		wg.Wait()

		//fmt.Println(temp[0] + temp[1] + temp[2] + temp[3] + temp[4] + temp[5])
		out <- (temp[0] + temp[1] + temp[2] + temp[3] + temp[4] + temp[5])
		//out <- temp
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
