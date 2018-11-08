package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "abcdec"

	//包含   指定字符串中是否包含某个字符串   返回bool
	fmt.Println(strings.Contains(s, "b"), strings.Contains(s, "e"))

	//索引   在指定字符串中查找某个字符串首次出现的位置
	fmt.Println(strings.Index(s, "c"))

	//字符串切割成数组
	id := "1,2,3,4,5,6"
	id_arr := strings.Split(id, ",")
	fmt.Println(id_arr)

	//将数组合并成新的字符串
	new_id := strings.Join(id_arr, ",")
	fmt.Println(new_id)

	//判断字符串是否含有某个前缀
	fmt.Println(strings.HasPrefix(s, "a"))
	//判断字符串是否含有某个后缀
	fmt.Println(strings.HasSuffix(s, "c"))

}
