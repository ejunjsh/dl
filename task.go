package main

import (
	"io"
	"time"
	"sync"
	"sync/atomic"
	"os"
	"net/http"
	"errors"
	"fmt"
)

type task struct {
	done chan struct{}
	src io.ReadCloser
	dst io.WriteCloser
	bytePerSecond float64
	err error
	startTime time.Time
	endTime time.Time
	mutex sync.Mutex
	readNum int64
	fileSize int64
	filename string
	buffer []byte
	client *client
}

func (t *task) getReadNum() int64{
	if t == nil {
		return 0
	}
	return atomic.LoadInt64(&t.readNum)
}

func newTask(url string) *task{
	return &task{client:newClient(url),done:make(chan struct{},1),buffer:make([]byte,32*1024)}
}

func (t *task) start(){
	var dst *os.File
	var rn,wn int
	var filename string
	req,_:= http.NewRequest("GET",t.client.url,nil)
	rep,err:= t.client._client.Do(req)
	if err!=nil{
		goto done
	}else if rep.StatusCode!=200{
		err=errors.New(fmt.Sprintf("wrong response %d",rep.StatusCode))
		goto done
	}

	filename, err = guessFilename(rep)

	dst,err=os.Create(filename)
	if err!=nil{
		goto done
	}
	t.dst=dst
	t.src=rep.Body
	t.fileSize=rep.ContentLength

	go t.bps()

	t.startTime=time.Now()

loop:
	rn,err=t.src.Read(t.buffer)

	if err!=nil||rn==0{
		goto done
	}

	wn,err=t.dst.Write(t.buffer[:rn])

	if err!=nil{
		goto done
	} else if rn!=wn {
		err = io.ErrShortWrite
		goto done
	}else{
		atomic.AddInt64(&t.readNum,int64(rn))
		goto loop
	}

done:
	t.err=err
	close(t.done)
	t.endTime=time.Now()
	return
}

func (t *task) bps(){
	var prev int64
	then := t.startTime

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.done:
			return

		case now := <-ticker.C:
			d := now.Sub(then)
			then = now

			cur := t.getReadNum()
			bs := cur - prev
			prev = cur

			t.mutex.Lock()
			t.bytePerSecond = float64(bs) / d.Seconds()
			t.mutex.Unlock()
		}
	}
}

func (t *task) getSpeed() string{
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return formatBytes(int64(t.bytePerSecond))
}


func (t *task) getETA() string{
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.fileSize==0||t.bytePerSecond==0{
		return "--"
	}else {
		return formatTime((t.fileSize-t.getReadNum())/int64(t.bytePerSecond))
	}
}