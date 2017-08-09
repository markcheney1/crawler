package main

import (
	"github.com/astaxie/beego/httplib"
	log "github.com/sirupsen/logrus"
	"fmt"
	"time"
	"flag"
	"os"
	"bufio"
	"bytes"
)

var HotmomentRequestHost string = "https://api.gotokeep.com/social/v3/timeline/hot"
var host string = "https://api.gotokeep.com/social/v3/timeline/hot?lastId=%s"
type CommonResponse struct {
	ErrorCode int `json:"errorCode"`
	Text      string `json:"text"`
	Data      interface{} `json:"data"`
	Ok        bool `json:"ok"`
}

type HotMomentData struct {
	Entries []Entry `json:"entries"`
	LastId  string `json:"lastId"`
}

type Entry struct {
	Id      string `json:"id"` //用于匹配hot下面的topic
}


var filePath = flag.String("file", "", "file用于存储topicId")

func main()  {
	flag.Parse()
	if *filePath == "" {
		log.Fatal("file的值为空")
	}
	f, err := os.Create(*filePath)
	if err != nil {
		log.Fatalf("创建文件失败,err:%v", err)
		fmt.Println()
	}
	buf := bufio.NewWriter(f)
	ids, lastId, err := GetId(HotmomentRequestHost)
	if err != nil {
		log.Errorf("请求获取hotcomment的话题id失败，err:%v", err)
	}
	buf.Write(ids)
	buf.Flush()
	var tmp string
	for {
		ids, lastId, err = GetId(fmt.Sprintf(host, lastId))
		if err != nil {
			lastId = tmp
			time.Sleep(30*time.Second)
			continue
		}
		if lastId == "" {
			break
		}
		buf.Write(ids)
		buf.Flush()
		tmp = lastId
	}
}

func GetId(host string) ([]byte, string, error) {
	common := CommonResponse{}
	data := HotMomentData{}
	common.Data = &data
	req := httplib.Get(host)
	err := req.ToJSON(&common)
	if err != nil {
		return nil, "", err
	}
	buf := bytes.Buffer{}
	for _, v := range common.Data.(*HotMomentData).Entries {
		buf.WriteString(v.Id)
		buf.WriteByte('\n')
	}
	return buf.Bytes(), common.Data.(*HotMomentData).LastId, nil
}


