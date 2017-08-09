package models

//"/social/v3/timeline/hot?lastId=59795807ca9aed15ca3a0c12"
var HotmomentRequestHost string = "https://api.gotokeep.com/social/v3/timeline/hot"

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
	City    string `json:"city"`
	Content string `json:"content"`
	Comments int   `json:"comments"`
	Id      string `json:"id"` //用于匹配hot下面的topic
	Likes   int `json:"likes"`
	Photo  string `json:"photo"`
	Author  WithAuthor `json:"author"`
	Created string `json:"created"`
	FavoriteCount int `json:"favoriteCount"`
}

type WithAuthor struct {
	Gender   string `json:"gender"`
	UserId   string `json:"_id"`
	Avatar   string `json:"avatar"`
	UserName string `json:"username"`
}

/**************************************************************/
//https://api.gotokeep.com/v1.1/entries/5979c08d879cca091d91461d/comments?limit=20&reverse=true
var HotMomentCommentHost string = "https://api.gotokeep.com/v1.1/entries/%s/comments?limit=20&reverse=true"
var HotMomentCommentHost2 string = "https://api.gotokeep.com/v1.1/entries/%s/comments?lastId=%s&limit=20&reverse=true"

type HotMomentComment struct {
	Ok        bool `json:"ok"`
	ErrorCode int `json:"errorCode"`
	Text      string `json:"text"`
	Data      []HotMomentComentData `json:"data"`
	Count     int `json:"count"`
}

type HotMomentComentData struct {
	Author  WithAuthor `json:"author"`
	Created string `json:"created"`
	Content string `json:"content"`
	Likes   int `json:"likes"`
	Id      string `json:"id"`
}

/*********************************************/

type Response struct {
	Platform string `json:"platform"`
	Category string `json:"category"`
	ImageUrl string `json:"image_url"`
	Title string `json:"title"`
	Author string `json:"author"`
	Gender int `json:"gender"`
	Age     int `json:"age"`
	Location string `json:"location"`
	Article string `json:"article"`
	LikeCount int `json:"likeCount"`
	CommentCount int `json:"commentCount"`
	CollectCount int `json:"collectCount"`
	Date string `json:"date"`
	Comments []CommentsResponse `json:"comments"`

}

type CommentsResponse struct {
	Text string `json:"text"`
	Gender int  `json:"gender"`
	Age int `json:"age"`
	Location string `json:"location"`
}

//https://api.gotokeep.com/account/v3/userinfo/5652eb19a00ca510d518c36b
var UserInfoHost string = "https://api.gotokeep.com/account/v3/userinfo/%s"
/************************************/
type UserData struct {
	Birthday string `json:"birthday"`
	City     string `json:"city"`
}

/*
{
"platform": "小红书",
"category": "美妆",
"image_url": "xxx",
"title": "xxx",
"author": "xxx",
"gender": 0, // 性别：0：女，1：男
"age": 25,  // 年龄，如果 APP 提供的是出生日期，则请换算成年龄
"location": "上海",
"article": "xxx",
"likeCount": 123,
"commentCount": 3,
"collectCount": 234,
"date": "2017-06-12",
"comments": [
{
"text": "xxx",
"gender": 0,
"age": 18,
"location": "北京"
},
{
"text": "xxx",
"gender": 0,
"age": 18,
"location": "北京"
}
]
}*/
