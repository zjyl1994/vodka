package main

import (
	"fmt"
	"time"

	"github.com/zjyl1994/vodka"
)

func main() {
	vodka.Init("Vodka Demo")
	vodka.RegisterError("custom", 418)
	vodka.Handle("GET", "/ping", vodka.Rules{
		QueryString: vodka.Rule{
			"id": []string{"required", "string", "alpha_num", "max:50"},
		},
	}, true, pingHandler)
	vodka.Handle("GET", "/optional", vodka.Rules{
		QueryString: vodka.Rule{
			"param": []string{"string", "max:50"},
		},
	}, true, optionalHandler)
	vodka.Handle("GET", "/error", vodka.Rules{}, true, eHandler)
	vodka.Handle("GET", "/nowtime", vodka.Rules{}, true, nowHandler)
	vodka.Handle("GET", "/db", vodka.Rules{}, true, dbHandler)
	if err := vodka.Run(); err != nil {
		fmt.Println("ERROR", err)
	}
}

func optionalHandler(c *vodka.Context, params vodka.Datas) (interface{}, vodka.Error) {
	param, ok := params.QueryString["param"]
	if ok {
		return vodka.H{
			"hasParam": true,
			"param":    param,
		}, vodka.NewError(nil)
	} else {
		return vodka.H{
			"hasParam": false,
		}, vodka.NewError(nil)
	}
}

func pingHandler(c *vodka.Context, params vodka.Datas) (interface{}, vodka.Error) {
	id := params.QueryString["id"].(string)
	return vodka.H{
		"id": id,
	}, vodka.NewError(nil)
}

func eHandler(c *vodka.Context, params vodka.Datas) (interface{}, vodka.Error) {
	return "just error", vodka.NewMessageError("custom", "this is error msg %d", 777)
}

func dbHandler(c *vodka.Context, params vodka.Datas) (interface{}, vodka.Error) {
	return vodka.DB, vodka.NoError
}

func nowHandler(c *vodka.Context, params vodka.Datas) (interface{}, vodka.Error) {
	nowtime := vodka.Time(time.Now())
	return vodka.H{
		"now": nowtime,
	}, vodka.NewError(nil)
}
