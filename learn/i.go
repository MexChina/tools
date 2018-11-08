package main
import "fmt"
func main() {

   var a int = 100
   var b int = 200
   var ret int
   var c int
   var d int
   ret = max(a, b)

   c,d = swap(a,b)

   fmt.Printf( "最大值是 : %d\n", ret )
   fmt.Printf("c : %d\n",c)
   fmt.Printf("d : %d\n",d)
}

/* 函数返回两个数的最大值 */
func max(num1, num2 int) int {
   var result int
   if (num1 > num2) {
      result = num1
   } else {
      result = num2
   }
   return result 
}
//函数返回多个值
func swap(x, y int) (int, int) {
   return y, x
}

/**
 * func：函数由 func 开始声明
 * function_name：函数名称，函数名和参数列表一起构成了函数签名。
 * @param parameter list：参数列表，参数就像一个占位符，
 *        当函数被调用时，你可以将值传递给参数，
 *        这个值被称为实际参数。
 *        参数列表指定的是参数类型、顺序、及参数个数。
 *        参数是可选的，也就是说函数也可以不包含参数。
 * @return return_types：返回类型，函数返回一列值。
 * return_types 是该列值的数据类型。
 * 有些功能不需要返回值，这种情况下 
 * return_types 不是必须的。函数体：函数定义的代码集合。
 */
// func function_name( [parameter list] ) [return_types] {
//    函数体
// }




