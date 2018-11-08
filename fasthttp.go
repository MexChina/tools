package tools

import (
    "fmt"

    // "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	// "time"
	"runtime"
	"log"
)

// fasthttprouter.Params 是路由匹配得到的参数，如规则 /hello/:name 中的 :name
// func httpHandle(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
//     fmt.Fprintf(ctx, "hello fasthttp")
// }

func httpHandle(ctx *fasthttp.RequestCtx) {
    // resCh := make(chan string, 1)
    //go func() {

		//defer func(){
		//	err := recover()
		//	if err != nil {
		//		fmt.Println("error to chan put.")
		//	}
		//}()

		// 这里使用 ctx 参与到耗时的逻辑中
		//fmt.Println(ctx.Host)
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	log.Println(memStat.Sys,memStat.Alloc)
		ctx.WriteString("get: abc = ")

        // time.Sleep(1 * time.Second)
        // resCh <- string("aaaaa")
    //}()

    // RequestHandler 阻塞，等着 ctx 用完或者超时
    // select {
    // case <-time.After(3 * time.Second):
    //     ctx.TimeoutError("timeout")
    // case r := <-resCh:
    //     ctx.WriteString("get: abc = " + r)
    // }
}

func main() {
    // 使用 fasthttprouter 创建路由
    // router := fasthttprouter.New()
    // router.GET("/", httpHandle)
    // if err := fasthttp.ListenAndServe("0.0.0.0:12345", router.Handler); err != nil {
    //     fmt.Println("start fasthttp fail:", err.Error())
	// }
	

	if err := fasthttp.ListenAndServe("0.0.0.0:12345", httpHandle); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}

}


