package main
//eg1 hello
//var c = make(chan int)
//var a string
//func f()  {
//	a = "hello"
//	c <- 0
//}
//func main() {
//	go f()
//	<- c
//	print(a)
//}
//
//eg2 1 10
func main()  {
	println(v1())
	println(v2())
}

func v1()(value int)  {
	defer func() {
		value++
	}()
	return value
}

func v2() int {
	value := 10
	defer func() {
		value ++
	}()
	return value
}