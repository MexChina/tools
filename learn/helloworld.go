
//包声明
package main
//包引入
import "fmt"

//函数
func main(){

	//常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型。
	const LENGTH int = 10	//显式类型定义
    const WIDTH = 5   		//隐式类型定义
    var area int
    const a, b, c = 1, false, "str" //多重赋值

    area = LENGTH * WIDTH
    fmt.Printf("面积为 : %d", area)
    println()
    println(a, b, c)  
}