package tools

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
)

func main()  {



	for{
		send()
	}

}

func send()  {
	jsonStr := `{
	"header": {
		"post_url": "http://bi.rpc/bi_tobbi",
		"local_ip": "127.0.0.1",
		"log_id": "123456",
		"session_id": "",
		"product_name": "wiki"
	},
	"request": {
		 "c": "tobbi", "m": "bi_company_jd_weekly_created_flash_count","p": {"ids":"1847729,2189688,2204717","src":"11,2"}
	}
}`

	req, err := http.NewRequest("POST", "http://127.0.0.1:51081/bi_tobbi", bytes.NewBuffer([]byte(jsonStr)))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}