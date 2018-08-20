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
		isfirst:=true
		go func() {
			for{
				select{
					case <-ticker.C:

						if !isfirst{
							termutil.ClearLines(int16(len(ts)))
						}
						for _,t:=range ts{
							var buf string
							var etaBuf string
							if t.err!=nil{
								buf=fmt.Sprintf("error:%s",t.err.Error())
								etaBuf=""
							}else if  t.getReadNum()>0{
								var eta string
								if t.fileSize<=0 ||  t.getBps()==0{
									eta="--"
								}else {
									eta=formatTime((t.fileSize-t.getReadNum())/int64(t.getBps()))
								}
								etaBuf=fmt.Sprintf("%s (%s/s)",eta,formatBytes(int64(t.getBps())))
								buf=fmt.Sprintf("%s/%s(%.2f%%)",formatBytes(t.getReadNum()),formatBytes(t.fileSize),100*float64(t.getReadNum())/float64(t.fileSize))
							}else {
								buf="waiting..."
								etaBuf=""
							}
							count:=cellCount(buf)
							etacount:=cellCount(etaBuf)
							r:=width-etacount-count
							if r>0 && t.fileSize>0 {
								buf+="["
								etaBuf="]"+etaBuf

								ratio:=float64(t.getReadNum())/float64(t.fileSize)
								r-=2
								bar:=strings.Repeat(" ",r)
								c:= int(float64(r)*ratio)
								progress:=""
								if c!=0{
									progress=strings.Repeat("=",c)
								}
								bar=strings.Join([]string{progress,">",bar[c+1:]},"")
								buf=strings.Join([]string{buf,bar,etaBuf},"")

							} else if r<0{
								buf=buf[:width]
							}
							fmt.Println(buf)

						}
						isfirst=false

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

