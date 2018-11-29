package main

import "github.com/zjyl1994/vodka"

func main() {
	vodka.Init("Vodka Demo")
	vodka.RegisterError("custom",418)
	vodka.Handle("GET","/ping",vodka.Rules{
		QueryString:vodka.Rule{
			"id":              []string{"required","string", "alpha_num", "max:50"},
		},
	},true,pingHandler)
	vodka.Handle("GET","/error",vodka.Rules{},true,eHandler)
	vodka.Handle("GET","/db",vodka.Rules{},true,dbHandler)
	vodka.Run()
}

func pingHandler(c *vodka.Context,params vodka.Datas)(interface{},vodka.Error){
	id := params.QueryString["id"].(string) 
	return vodka.H{
		"id":id,
	},vodka.NewError(nil)
}

func eHandler(c *vodka.Context,params vodka.Datas)(interface{},vodka.Error){
	return "just error",vodka.NewMessageError("custom","this is error msg %d",777)
}

func dbHandler(c *vodka.Context,params vodka.Datas)(interface{},vodka.Error){
	return vodka.DB,vodka.NoError
}