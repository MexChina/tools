package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strconv"
	"bytes"
	"net/http"
)

func main() {

	//db, err := sql.Open("mysql", "icdc:rg62fnme2d68t9cmd3d@tcp(192.168.8.250:3306)/icdc_0?charset=utf8&parseTime=True&loc=Local")
	//db, err := sql.Open("mysql", "icdc:rg62fnme2d68t9cmd3d@tcp(192.168.8.251:3306)/icdc_1?charset=utf8&parseTime=True&loc=Local")


	for i := 0; i < 33; i++ {
		read(i)
	}

}

func read(i int) {
	fmt.Println(i, "start read")
	db, err := sql.Open("mysql", "devuser:devuser@tcp(192.168.1.201:3310)/icdc_1?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//db2, err := sql.Open("mysql", "icdc:rg62fnme2d68t9cmd3d@tcp(192.168.8.251:3306)/icdc_1?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//var db *sql.DB

	//if i%2 == 0 {
	//	db = db1
	//} else {
	//	db = db2
	//}
	sql := "select id from icdc_"+strconv.Itoa(i)+".resumes order by id"
	var id int
	rst, err := db.Query(sql)
	if err != nil {
		panic(err)
		return
	}

	for rst.Next() {
		if err := rst.Scan(&id); err != nil {
			panic(err)
		}
		doJson(id)
	}
	defer rst.Close()

}

func doJson(id int) {
	ids := strconv.Itoa(id)
	bodyjson := "{\"header\": {\"log_id\": "+ids+"},\"request\": {\"c\": \"Logic_refresh\",\"m\": \"cv_trade\",\"p\": {\"ids\":" + ids + "}}}"
	var jsonStr = []byte(bodyjson)
	req, err := http.NewRequest("POST", "http://dev.icdc.rpc/icdc_refresh",bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
