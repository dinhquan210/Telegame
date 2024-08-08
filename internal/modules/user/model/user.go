package usermodel

type User struct {
	Id         int64   `json:"id"`
	NickName   string  `json:"nickName"`
	TotalPoint float64 `json:"totalPoint"`
}
