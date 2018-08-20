package main

import (
	"os"
	"fmt"
	"time"
	"github.com/ejunjsh/dl/termutil"
	"strings"
)

func main()  {
	if len(os.Args)==1{
		fmt.Println("usage: dl [url...]")
		return
	}else {
		ts:=make([]*task,len(os.Args)-1)
		for i,url:=range os.Args[1:]{
			t:=newTask(url)
			if t!=nil{
				go t.start()
				ts[i]=t
			}
		}


		width,_:=termutil.TerminalWidth()

		ticker:=time.NewTicker(time.Second)

		go func() {
			for{
				select{
					case <-ticker.C:
						for _,t:=range ts{
							var buf string
							if t.err!=nil{
								buf=fmt.Sprintf("error:%s",t.err.Error())
							}else if  t.getReadNum()>0{
								var eta string
								if t.fileSize<=0 ||  t.getBps()==0{
									eta="--"
								}else {
									eta=formatTime((t.fileSize-t.getReadNum())/int64(t.getBps()))
								}
								buf=fmt.Sprintf("%s/%s(%.2f%%) ETA %s (%s/s)",formatBytes(t.getReadNum()),formatBytes(t.fileSize),100*float64(t.getReadNum())/float64(t.fileSize),eta,formatBytes(int64(t.getBps())))
							}else {
								buf="waiting..."
							}
							count:=cellCount(buf)
							r:=width-count
							if r>0{
								buf+=strings.Repeat(" ",r)
							} else if r<0{
								buf=buf[:width]
							}
							fmt.Println(buf)

						}
						termutil.ClearLines(int16(len(ts)))
				}
			}
		}()

		for _,t:=range ts{
			if t!=nil{
				<-t.done
			}
		}

		fmt.Println("finished")
	}
}