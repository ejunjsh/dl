package main

import (
	"os"
	"fmt"
	"net/http"
	"time"
	"io"
)

func main()  {
	if len(os.Args)==1{
		fmt.Println("usage: dl [url...]")
		return
	}else {
		ts:=make([]*task,len(os.Args)-1)
		for i,url:=range os.Args[1:]{
			req,_:= http.NewRequest("GET",url,nil)
			t:=newClient().do(req)
			if t!=nil{
				go t.start()
				ts[i]=t
			}
		}


		ticker:=time.NewTicker(200*time.Millisecond)

		go func() {
			for{
				select{
					case <-ticker.C:
						for _,t:=range ts{
							fmt.Printf("Downloaded %s(%s/s)\n",formatBytes(t.getReadNum()),formatBytes(int64(t.getBps())))
						}
				}
			}
		}()

		for _,t:=range ts{
			if t!=nil{
				<-t.done
				if t.err!=nil&&t.err!=io.EOF{
					fmt.Println(t.err)
				} else {
					fmt.Println("finished")
				}
			}
		}
	}
}