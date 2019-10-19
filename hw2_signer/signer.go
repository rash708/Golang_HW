package main

import (
	"fmt"
	"os"
	"strconv"
)

//ExecutePipeline SingleHash -> Multihash -> CombineResults
func ExecutePipeline(data string) {
	CombineResults(MultiHash(SingleHash(data)))
}

//SingleHash crc32(data)+"~"+crc32(md5(data))
func SingleHash(data string) string {
	return DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

//MultiHash crc32(th+data)), где data - htpkmnfn
func MultiHash(data string) string {
	var out string
	var temp string
	for th := 0; th < 6; th++ {
		temp = DataSignerCrc32(strconv.Itoa(th) + data)
		fmt.Printf("%v "+temp+"\n", th)
		out += temp
	}
	return out
}

//CombineResults res1_res2_.._resN
func CombineResults(data string) string {
	out := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
	return out
}

func main() {
	fmt.Println(SingleHash(os.Args[1]))
	fmt.Println(MultiHash(SingleHash(os.Args[1])))
}

// сюда писать код
