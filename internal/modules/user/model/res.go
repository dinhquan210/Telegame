package usermodel

type GetResult struct {
	NickName   string  `json:"nickName"`
	TimeRes    int64   `json:"timeRes"`
	Result     string  `json:"result"`
	TotalPoint float64 `json:"totalPoint"`
}
