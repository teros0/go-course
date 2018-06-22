package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

var steps []struct{}

func init() {
	steps = make([]struct{}, 6)
	fmt.Printf("%+v\n", steps)
}

func ExecutePipeline(jobs ...job) {
	in, out := make(chan interface{}, 100), make(chan interface{}, 100)
	wg := &sync.WaitGroup{}
	for _, j := range jobs {
		wg.Add(1)
		go worker(wg, j, in, out)
		in = out
		out = make(chan interface{}, 100)
	}
	wg.Wait()
}

func worker(wg *sync.WaitGroup, j job, in, out chan interface{}) {
	j(in, out)
	wg.Done()
}

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}) {
	for v := range in {
		go func(interface{}) {
			i := v.(int)
			s := strconv.Itoa(i)
			fmt.Println(s, "SingleHash data", s)
			p1 := DataSignerCrc32(s)
			fmt.Println(s, "SingleHash crc32(data)", p1)
			p21 := DataSignerMd5(s)
			fmt.Println(s, "SingleHash md5(data)", p21)
			p22 := DataSignerCrc32(p21)
			fmt.Println(s, "SingleHash crc32(md5(data))", p22)
			res := fmt.Sprintf("%s~%s", p1, p22)
			fmt.Println(s, "SingleHash result", res)
			out <- res
		}(v)
	}
	return
}

// MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5
// ( т.е. 6 хешей на каждое входящее значение ), потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход (и ушло на выход из SingleHash)
func MultiHash(in, out chan interface{}) {
	var res string
	for v := range in {
		go func() {
			s := v.(string)
			for _, i := range []int{0, 1, 2, 3, 4, 5} {
				go func(i int) {
					hi := DataSignerCrc32(string(i) + s)
					fmt.Printf("%s MultiHash: crc32(th+step1)) %d %s\n", s, i, hi)
					res += hi
					out <- hi
				}(i)
			}
		}()
	}
	fmt.Printf("MultiHash: result %s\n", res)
	return
}

//CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
func CombineResults(in, out chan interface{}) {
	var res []string
	for v := range in {
		res = append(res, v.(string))
	}

	sort.Strings(res)
	fmt.Printf("CombineResults ")
	for i, s := range res {
		fmt.Printf("%s", s)

		if i != len(in)-1 {
			fmt.Printf("_")
		}
	}
	fmt.Println()
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	// inputData := []int{0,1}
	var testResult string

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			fmt.Println("doing first job")
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			testResult = data
		}),
	}
	start := time.Now()
	ExecutePipeline(hashSignJobs...)
	end := time.Now().Sub(start)

	fmt.Println(testResult)
	fmt.Println(end)
}
