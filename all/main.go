package main

import (
	"models"
	"strings"
	"crypto/tls"
	"time"
	"github.com/astaxie/beego/httplib"
	log "github.com/sirupsen/logrus"
	"fmt"
	"encoding/json"
	"os"
	"bufio"
	"flag"
)

var filePath = flag.String("file", "", "爬虫输出")

func main() {
	flag.Parse()
	if *filePath == "" {
		log.Fatalf("no file path")
	}
	hotMomentChannel := make(chan *models.CommonResponse, 100)
	logs := make(chan *models.Response, 100)
	response, err := HotRequest(models.HotmomentRequestHost)
	if err != nil {
		return
	}
	hotMomentChannel <- response
	for i := 0; i < 30; i++ {
		go ResponseLog(hotMomentChannel, logs)
	}
	go PrintJson(logs, *filePath)
	for {
		fmt.Println("lll")
		if response.Ok == true {
			if id := response.Data.(*models.HotMomentData).LastId; id != "" {
				strTmp := models.HotmomentRequestHost + "?lastId=" + id
				fmt.Println(strTmp)
			try:	response, err = HotRequest(strTmp)
				if err != nil {
					fmt.Println("请求出错, err:", err)
					time.Sleep(20* time.Second)
					goto try
				}
				fmt.Println(response)
				hotMomentChannel <- response
			} else {
				fmt.Println("程序退出了。。。。。。。。。")
				break
			}
		} else {
			time.Sleep(2*time.Minute)
		}
	}

}

func HotRequest(host string) (*models.CommonResponse, error) {
	commonresponse := models.CommonResponse{}
	data := models.HotMomentData{}
	commonresponse.Data = &data
	req := httplib.Get(host)
	if strings.Contains(host, "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	req.SetTimeout(10*time.Second, 5*time.Second)

	err := req.ToJSON(&commonresponse)
	if err != nil {
		log.Errorf("放回值解析成json失败,err:%v", err)
		/*b,_ := req.Bytes()
		fmt.Println(string(b))*/
		return nil, err
	}
	if commonresponse.ErrorCode != 0 {
		log.Errorf("放回的状态码不是0，更新fake评论失败")
		return nil, err
	}
	return &commonresponse, nil
}

func CommentRequest(host string) (*models.HotMomentComment, error) {
	/*	defer func() {
			if err := recover(); err != nil {
				fmt.Println("请求comment的时候出现panic")
			}
		}()*/
	commonresponse := models.HotMomentComment{}
	data := []models.HotMomentComentData{}
	commonresponse.Data = data
	req := httplib.Get(host)
	if strings.Contains(host, "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	req.SetTimeout(5*time.Second, 5*time.Second)
	err := req.ToJSON(&commonresponse)
	if err != nil {
		log.Errorf("放回值解析成json失败,err:%v", err)
		/*b,_ := req.Bytes()
		fmt.Println(string(b))*/
		return nil, err
	}
	if commonresponse.ErrorCode != 0 {
		log.Errorf("放回的状态码不是0，更新fake评论失败")
		return nil, err
	}
	return &commonresponse, nil
}

func GetUserInfo(host string) (int, string, error) {
	timeFormat := "2006-01-02T15:04:05.999Z"
	commonresponse := models.CommonResponse{}
	data := models.UserData{}
	commonresponse.Data = &data
	req := httplib.Get(host)
	if strings.Contains(host, "https") {
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	req.SetTimeout(5*time.Second, 5*time.Second)

	err := req.ToJSON(&commonresponse)
	if err != nil {
		log.Errorf("放回值解析成json失败,err:%v", err)
		/*b,_ := req.Bytes()
		fmt.Println(string(b))*/
		return 0, "", err
	}
	if commonresponse.ErrorCode != 0 {
		log.Errorf("放回的状态码不是0")
		return 0, "", err
	}

	year, err := time.Parse(timeFormat, commonresponse.Data.(*models.UserData).Birthday)
	if err != nil {
		log.Errorf("转换失败")
		return 0, "", err
	}
	return time.Now().Year() - year.Year(), commonresponse.Data.(*models.UserData).City, nil
}

func ResponseLog(ch chan *models.CommonResponse, logs chan *models.Response) {
	timeFormat := "2006-01-02T15:04:05.999Z"
	for {
		hot := <-ch
		for _, entry := range hot.Data.(*models.HotMomentData).Entries {
			fmt.Println("id值：", entry.Id)
			if data, err := CommentRequest(fmt.Sprintf(models.HotMomentCommentHost, entry.Id)); err != nil {
				fmt.Println("errrrrrrrrr", err)
				continue
			} else {
				age, _, err := GetUserInfo(fmt.Sprintf(models.UserInfoHost, entry.Author.UserId))
				if err != nil {
					age = 0
				}
				response := models.Response{}
				response.Platform = "keep"
				response.Category = "热门动态"
				response.ImageUrl = entry.Photo
				response.Title = ""
				response.Author = entry.Author.UserName
				if entry.Author.Gender == "M" {
					response.Gender = 1
				}
				response.Age = age
				response.Location = entry.City
				response.Article = entry.Content
				response.LikeCount = entry.Likes
				response.CommentCount = entry.Comments
				response.CollectCount = entry.FavoriteCount
				tt, err := time.Parse(timeFormat, entry.Created)
				if err == nil {
					response.Date = tt.Format("2006-01-02")
				}
				for _, value := range data.Data {
					gender := 0
					if value.Author.Gender == "M" {
						gender = 1
					}
					age, city, err := GetUserInfo(fmt.Sprintf(models.UserInfoHost, value.Author.UserId))
					if err != nil {
						age = 0
					}

					response.Comments = append(response.Comments, models.CommentsResponse{
						Text:     value.Content,
						Gender:   gender,
						Age:      age,
						Location: city,
					})
				}
				var length int
				if length = len(data.Data); length < 20 {
					logs <- &response
					continue
				}
				fmt.Println("id值：", data.Data[length-1].Id)
				for data, err = CommentRequest(fmt.Sprintf(models.HotMomentCommentHost, data.Data[length-1].Id)); err == nil; {
					for _, value := range data.Data {
						gender := 0
						if value.Author.Gender == "M" {
							gender = 1
						}
						age, city, err := GetUserInfo(fmt.Sprintf(models.UserInfoHost, value.Author.UserId))
						if err != nil {
							age = 0
						}

						response.Comments = append(response.Comments, models.CommentsResponse{
							Text:     value.Content,
							Gender:   gender,
							Age:      age,
							Location: city,
						})
					}
				}
				logs <- &response
			}
		}
	}
}

func PrintJson(ch chan *models.Response, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal("打开文件失败")
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	for {
		b, err := json.Marshal(<-ch)
		if err != nil {
			fmt.Println(err)
			continue
		}
		b = append(b, '\n')
		buf.Write(b)
		buf.Flush()
	}
}
