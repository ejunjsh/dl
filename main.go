package main

import (
	"fmt"
	"github.com/ejunjsh/dl/termutil"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		usage := `usage: dl [[rate limit]:url...]
rate limit: limit the speed,unit is KB
url...: urls you want to download`
		fmt.Println(usage)
		return
	} else {
		ts := make([]*task, len(os.Args)-1)
		for i, url := range os.Args[1:] {
			t := newTask(url)
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
	}
}

func updateTerm(isGetWidth bool, ts []*task, width int) {
	for _, t := range ts {
		var buf string
		if t.err != nil && t.err != io.EOF {
			buf = fmt.Sprintf("error:%s", t.err.Error())
		} else if t.getReadNum() > 0 {
			var etaBuf string
			var fileSizeBuf string
			if t.fileSize <= 0 {
				fileSizeBuf = fmt.Sprintf("%s", formatBytes(t.getReadNum()))
			} else {
				fileSizeBuf = fmt.Sprintf("%s/%s(%.2f%%)", formatBytes(t.getReadNum()), formatBytes(t.fileSize), 100*float64(t.getReadNum())/float64(t.fileSize))
			}

			etaBuf = fmt.Sprintf("%s (%s/s)", t.getETA(), t.getSpeed())

			if isGetWidth || t.fileSize > 0 {
				r := width - cellCount(fileSizeBuf+etaBuf)
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
					buf = strings.Join([]string{fileSizeBuf, bar, etaBuf}, "")

				} else if r < 0 {
					buf = buf[:width]
				} else {
					buf = strings.Join([]string{fileSizeBuf, etaBuf}, "")
				}
			} else {
				buf = strings.Join([]string{fileSizeBuf, etaBuf}, "")
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
