/*
package linear

import (
	"fmt"
	"mailru/hw2_signer"
	"sort"
)

var steps []struct{}

func init() {
	steps = make([]struct{}, 6)
	fmt.Printf("%+v\n", steps)
}

func ExecutePipeline(data []string) {
	for _, d := range data {
		s1 := SingleHash(d)
		s2 := MultiHash(s1)
		CombineResults([]string{s2})
	}
}

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(data string) (res string) {
	fmt.Println(data, "SingleHash data", data)
	p1 := DataSignerCrc32(data)
	fmt.Println(data, "SingleHash crc32(data)", p1)
	p21 := DataSignerMd5(data)
	fmt.Println(data, "SingleHash md5(data)", p21)
	p22 := DataSignerCrc32(p21)
	fmt.Println(data, "SingleHash crc32(md5(data))", p22)
	res = fmt.Sprintf("%s~%s", p1, p22)
	fmt.Println(data, "SingleHash result", res)
	return
}

// MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5
// ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход (и ушло на выход из SingleHash)
func MultiHash(data string) (res string) {
	for _, i := range []int{0, 1, 2, 3, 4, 5} {
		hi := DataSignerCrc32(string(i) + data)
		fmt.Printf("%s MultiHash: crc32(th+step1)) %d %s\n", data, i, hi)
		res += hi
	}
	fmt.Printf("%s MultiHash: result %s\n", data, res)
	return
}

//CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in []string) {
	sort.Strings(in)
	fmt.Printf("CombineResults ")
	for i, s := range in {
		fmt.Printf("%s", s)

		if i != len(in)-1 {
			fmt.Printf("_")
		}
	}
	fmt.Println()
}
*/