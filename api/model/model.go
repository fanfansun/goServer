package model

type UserModelInterface struct {
	// 开头大写，导出 json 标记
	Username string `json:"user_name"`
	Password string `json:"password"`
}

/*
{
 user_name:xxx
 password:xxx
}

*/
// 请求json返回格式

//response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

// Data model
type VideoInfo struct {
	Id           string
	AuthorId     int
	Name         string
	DisplayCtime string
}

type Comment struct {
	Id      string
	VideoId string
	Author  string
	Content string
}

type SimpleSession struct {
	Username string //login name
	TTL      int64
}
