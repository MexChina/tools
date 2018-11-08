package tools

import (
	_ "net/http/pprof"
	"log"
	"net/http"
)

func main(){

	go func(){
		log.Fatal(http.ListenAndServe(":6060",nil))
	}()

}
