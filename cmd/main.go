package main

import (
	"context"
	"fmt"
	"net/http"
	usermodel "telegame/internal/modules/user/model"
	"telegame/utils/binance"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// save price btc redis
	go func() {
		for {
			price, err := binance.GetPrice()
			if err != nil {
				fmt.Println(err)
				return
			}
			timeSave := time.Now().Unix()
			key := "btc_price_" + fmt.Sprint(timeSave)
			rdb.Set(context.Background(), key, price, 15*time.Second)
		}
	}()

	// create api submit
	userApi := r.Group("/users")
	userApi.POST("/submit", func(c *gin.Context) {

		var request usermodel.SubmitRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "request error",
			})
		}

		// todo: get user from ctx
		user := usermodel.User{
			Id:         1,
			NickName:   "B2W Nauk",
			TotalPoint: 100,
		}

		// save request redis
		timeSubmit := time.Now().Unix()

		userRequestKey := fmt.Sprint(user.Id) + "_" + fmt.Sprint(timeSubmit)
		if err := rdb.SetNX(context.Background(), userRequestKey, request.PredictedValue, 10*time.Second).Err(); err != nil {
			fmt.Println(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "submit success",
			"timeSubmit": timeSubmit,
			"data":       request,
		})
	})

	// get result
	userApi.GET("/result", func(c *gin.Context) {

		var request usermodel.GetResultRequest
		if err := c.ShouldBind(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "request error",
			})
		}

		// todo: get user from ctx
		user := usermodel.User{
			Id:         1,
			NickName:   "B2W Nauk",
			TotalPoint: 100,
		}

		timeSubmit := request.TimeRes - 5

		fmt.Println("time submit", request)
		// get data submit by user
		userRequestKey := fmt.Sprint(user.Id) + "_" + fmt.Sprint(timeSubmit)
		value, err := rdb.Get(context.Background(), userRequestKey).Int()
		if err != nil {
			fmt.Println("get data user: "+userRequestKey, err)
			return
		}

		priceBtcSubmit, err := rdb.Get(context.Background(), "btc_price_"+fmt.Sprint(timeSubmit)).Float64()
		if err != nil {
			fmt.Println("get data price submit :"+"btc_price_"+fmt.Sprint(timeSubmit), err)
			return
		}

		priceBtcRes, err := rdb.Get(context.Background(), "btc_price_"+fmt.Sprint(request.TimeRes)).Float64()
		if err != nil {
			fmt.Println("get data price res :"+"btc_price_"+fmt.Sprint(request.TimeRes), err)
			return
		}

		response := ""
		dif := priceBtcRes - priceBtcSubmit
		if dif > 0 && value == 1 {
			response = "win"
		} else {
			response = "lose"
		}

		//todo: save history

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    response,
		})

	})
	r.Run()
}
