package main
import "fmt"
func main() {
   var a int = 10
   var b int = 20
 
   //if...
   if a < 20 {
       fmt.Printf("a 小于 20\n" )
   }

    //if...else...
   if a < 20 {
       fmt.Printf("a 小于 20\n" );
   } else {
       fmt.Printf("a 不小于 20\n" );
   }

   //if 嵌套
   if a == 10 {
       if b == 20 {
          fmt.Printf("a 的值为 10 ， b 的值为 20\n" );
       }
   }

   switch a {
      case 10: fmt.Printf("11111\n" )
      case 20: fmt.Printf("22222\n" )
      case 30,40,50 : fmt.Printf("33333\n" )
      default: fmt.Printf("00000\n" )
   }

   	//select是Go中的一个控制结构，类似于用于通信的switch语句。
   	//每个case必须是一个通信操作，要么是发送要么是接收。
	//select随机执行一个可运行的case。如果没有case可运行，
	//它将阻塞，直到有case可运行。一个默认的子句应该总是可运行的。

    var c1, c2, c3 chan int
    var i1, i2 int
    select {
      case i1 = <-c1:
         fmt.Printf("received ", i1, " from c1\n")
      case c2 <- i2:
         fmt.Printf("sent ", i2, " to c2\n")
      case i3, ok := (<-c3):  // same as: i3, ok := <-c3
         if ok {
            fmt.Printf("received ", i3, " from c3\n")
         } else {
            fmt.Printf("c3 is closed\n")
         }
      default:
         fmt.Printf("no communication\n")
    }    
}