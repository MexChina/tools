package tools

import (
	"github.com/robfig/cron"
	"fmt"
	"time"
)

func main()  {

	//  seconds   minute   hour   dayofmonth    month   dayofweek

	cs := cron.New()
	cs.AddFunc("*/10 * * * * *", func() {
		//log.Println("10 16 14 * * *")
		fmt.Println("每10秒执行")
		//db.Exec("replace into visual_scheduler_chan(chan_id,status) values(?,?)",chan_id,2)

	})

	cs.Start()
	i := 0
	for{
		if i == 2{
			time.Sleep(10*time.Second)
			cs := cron.New()
			cs.AddFunc("*/5 * * * * *", func() {
				fmt.Println("每5秒执行")
			})
			cs.Start()
		}
		i++
	}
}