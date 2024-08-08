package usermodel

type SubmitRequest struct {
	PredictedValue int `json:"predictedValue"` // 1: tăng, 0: giảm
}

type GetResultRequest struct {
	TimeRes int64 `json:"timeRes" form:"timeRes"`
}
