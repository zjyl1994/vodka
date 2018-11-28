package vodka

import (
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func validateQueryStringParamsForRequest(req *http.Request, rules map[string][]string) (Values, Datas) {
	params := make(Datas)
	amUtil := new(ArrayMapUtil)
	opts := govalidator.Options{
		Request:  req,
		Rules:    rules,
		Messages: govalidator.MapData{},
	}
	v := govalidator.New(opts)
	errsBag := v.Validate()
	for k, v:= range req.URL.Query(){
		params[k] = v[0]
	}
	amUtil.FilterOutRangeFields(params, amUtil.KeysA(rules))
	for rName,rValue := range rules{
		if !amUtil.IsArrayInclude(rValue,"allow_empty"){
			if amUtil.IsMapIncludeKey(params,rName) && len(params[rName].(string))==0{
				delete(params,rName)
			}
		}
	}
	return errsBag,params
}

func validateJSONParamsForRequest(req *http.Request, rules map[string][]string, allowEmptyBody bool) (Values, Datas) {
	params := make(Datas)
	amUtil := new(ArrayMapUtil)
	opts := govalidator.Options{
		Request:  req,
		Rules:    rules,
		Messages: govalidator.MapData{},
		Data:     &params,
	}
	errsBag := Values{}
	if !allowEmptyBody {
		if req.ContentLength == 0 {
			errsBag.Add("_error", "Empty body not allowed")
			return errsBag, params
		}
	} else {
		if req.ContentLength == 0 {
			return nil, params
		}
	}
	v := govalidator.New(opts)
	errsBag = v.ValidateJSON()
	amUtil.FilterOutRangeFields(params, amUtil.KeysA(rules))
	return errsBag, params
}