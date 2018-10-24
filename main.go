package main

import (
	"fmt"
	"github.com/ejunjsh/dl/termutil"
	"io"
	"os"
	"strings"
	"time"
	"github.com/urfave/cli"
	"log"
)

func printUsage(){
	usage := `usage: dl [--header <header> [ --header <header>]] [[rate limit:]url...]
--header: specify your http header,format is "key:value"
rate limit: limit the speed,unit is KB
url...: urls you want to download`
	fmt.Println(usage)
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "header",
		},
	}
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		printUsage()
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg()==0{
			printUsage()
			return nil
		}

		headers:=c.StringSlice("header")
		m:=parseHeaderFromStringSlice(headers)
		ts := make([]*task, c.NArg())
		for i, url := range c.Args() {
			t := newTask(url,m)
			if t != nil {
				go t.start()
				ts[i] = t
			}
		}

		var isGetWidth = true

		width, err := termutil.TerminalWidth()

		if err != nil {
			isGetWidth = false
		}

		ticker := time.NewTicker(time.Second)
		isfirst := true
		go func() {
			for {
				select {
				case <-ticker.C:
					if !isfirst {
						termutil.ClearLines(int16(len(ts)))
					}
					updateTerm(isGetWidth, ts, width)
					isfirst = false
				}
			}
		}()

		for _, t := range ts {
			if t != nil {
				<-t.done
			}
		}

		time.Sleep(time.Second)

		fmt.Println("finished")

		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func updateTerm(isGetWidth bool, ts []*task, width int) {
	for _, t := range ts {
		var buf string
		if t.err != nil && t.err != io.EOF {
			if t.filename==""{
				buf = fmt.Sprintf("error:%s",t.err.Error())
			}else {
				buf = fmt.Sprintf("%s:error:%s",t.filename,t.err.Error())
			}
		} else if t.getReadNum() > 0 {
			var etaBuf string
			var fileSizeBuf string
			var fnbuf string

			fnnum:=20
			fnbuf=showFileName(t.filename,fnnum)

			if t.fileSize <= 0 {
				fileSizeBuf = fmt.Sprintf("|%s", formatBytes(t.getReadNum()))
			} else {
				fileSizeBuf = fmt.Sprintf("|%s",formatBytes(t.fileSize))
			}

			etaBuf = fmt.Sprintf("%s|%s/s", t.getETA(), t.getSpeed())

			if isGetWidth && t.fileSize > 0 {
				r := width - cellCount(fileSizeBuf+etaBuf)-fnnum
				if r > 4 {
					fileSizeBuf += "["
					etaBuf = "]" + etaBuf

					ratio := float64(t.getReadNum()) / float64(t.fileSize)
					r -= 2
					bar := strings.Repeat(" ", r)
					c := int(float64(r) * ratio)
					progress := ""
					if c > 0 {
						progress = strings.Repeat("=", c)
					}
					if c+1<len(bar){
						bar = strings.Join([]string{progress, ">", bar[c+1:]}, "")
					}else {
						bar = strings.Join([]string{progress, ">"}, "")
					}
					buf = strings.Join([]string{fnbuf,fileSizeBuf, bar, etaBuf}, "")

				} else if r < 0 {
					buf = buf[:width]
				} else {
					buf = strings.Join([]string{fnbuf,fileSizeBuf, etaBuf}, "")
				}
			} else if t.fileSize>0 {
				buf = strings.Join([]string{fnbuf,fileSizeBuf,fmt.Sprintf("|%.2f%%",100*float64(t.getReadNum())/float64(t.fileSize)) ,etaBuf}, "")
			} else {
				buf = strings.Join([]string{fnbuf,fmt.Sprintf("|%s",formatBytes(t.getReadNum()))}, "")
			}
		} else {
			buf = "waiting..."
		}

		if isGetWidth {
			c := cellCount(buf)
			if c > width {
				buf = buf[:width]
			} else if c < width {
				buf = buf + strings.Repeat(" ", width-c)
			}
		}

		fmt.Println(buf)
	}

}

func showFileName(filename string,cap int) string{
	if len(filename)<cap{
		return strings.Join([]string{filename,strings.Repeat(" ",cap-len(filename))},"")
	}else {
		r:=[]rune(filename)
		//check if a rune is greater than 2 bytes
		if len(r)!=len(filename){

				l:=len(r)
				for
				{
					d:=string(r[:l])
					d2:=cellCount(d)
					if d2>cap{
						l--
						continue
					}else {
						return strings.Join([]string{d,strings.Repeat(" ",cap-d2)},"")
					}
				}

		}else {
			return filename[:cap]
		}

	}
}
