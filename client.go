package main

import (
	"net/http"
	"os"
	"log"
)

type client struct {
	_client *http.Client
}

func newClient() *client{
	return &client{
		&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
}



func (c *client) do(req *http.Request) *task{
	rep,err:= c._client.Do(req)
	if err!=nil{
		log.Println(err)
		return nil
	}else if rep.StatusCode!=200{
		log.Printf("wrong response %d \n",rep.StatusCode)
		return nil
	}
	dst,err:=os.Create("xxx")
	if err!=nil{
		log.Println(err)
		return nil
	}
	return &task{src:rep.Body,err:err,dst:dst,done:make(chan struct{},1),buffer:make([]byte,32*1024)}
}
