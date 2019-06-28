package vodka

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

func validateQueryStringParamsForRequest(req *http.Request, rules Rule) (Values, VariantMap) {
	params := make(map[string]interface{})
	amUtil := new(ArrayMapUtil)
	opts := govalidator.Options{
		Request:  req,
		Rules:    rules,
		Messages: govalidator.MapData{},
	}
	v := govalidator.New(opts)
	errsBag := v.Validate()
	for k, v := range req.URL.Query() {
		params[k] = v[0]
	}
	amUtil.FilterOutRangeFields(params, amUtil.KeysA(rules))
	for rName, rValue := range rules {
		if !amUtil.IsArrayInclude(rValue, "allow_empty") {
			if amUtil.IsMapIncludeKey(params, rName) && len(params[rName].(string)) == 0 {
				delete(params, rName)
			}
		}
	}
	return errsBag, params
}

func validateJSONParamsForRequest(req *http.Request, rules Rule, allowEmptyBody bool) (Values, VariantMap) {
	params := make(map[string]interface{})
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

func removeQueryStringEmptyField(req *http.Request) {
	resultQuery := url.Values{}
	for k, v := range req.URL.Query() {
		if strings.TrimSpace(v[0]) != "" {
			resultQuery.Set(k, v[0])
		}
	}
	req.URL.RawQuery = resultQuery.Encode()
}
