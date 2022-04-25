package fetch

import "fmt"

func ExampleGetFileLength() {
	requestURL := "http://httpbin.org/bytes/%d"

	var length1K int64 = 1024
	lengthResult, _ := GetFileLength(fmt.Sprintf(requestURL, length1K))
	fmt.Println(lengthResult)
	// Output:
	// 1024
}
