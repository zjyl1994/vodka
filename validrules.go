package vodka

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"regexp"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	govalidator.AddCustomRule("not_all", func(field string, rule string, message string, value interface{}) error {
		replacement := strings.TrimPrefix(rule, "not_all:")

		if value != nil {
			if reflect.TypeOf(value).String() == "string" && len(strings.Replace(value.(string), replacement, "", -1)) == 0 {
				if message != "" {
					return errors.New(message)
				}
				return fmt.Errorf("The %s field must a string and not be empty", field)
			}
		}

		return nil
	})

	govalidator.AddCustomRule("string", func(field string, rule string, message string, value interface{}) error {
		if value != nil && reflect.TypeOf(value).String() != "string" {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be a string", field)
		}

		return nil
	})

	govalidator.AddCustomRule("min_numeric", func(field string, rule string, message string, value interface{}) error {
		minNumber, _ := strconv.ParseInt(strings.TrimPrefix(rule, "min_numeric:"), 10, 64)
		val, err := strconv.ParseInt(value.(string), 10, 64)

		if err != nil || val < minNumber {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be an int numeric and greate than %d or equal to", field, minNumber)
		}

		return nil
	})

	govalidator.AddCustomRule("max_numeric", func(field string, rule string, message string, value interface{}) error {
		maxNumber, _ := strconv.ParseInt(strings.TrimPrefix(rule, "max_numeric:"), 10, 64)
		val, err := strconv.ParseInt(value.(string), 10, 64)

		if err != nil || val > maxNumber {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be an int numeric and less than %d or equal to", field, maxNumber)
		}

		return nil
	})

	govalidator.AddCustomRule("utc_timestamp", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		msInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be numeric", field)
		}

		// The value must be greater than utc-timestamp for 1970-01-01 08:00:12.133
		if msInt < 12133 {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be right utc-timestamp(mill-second) and greater than 12133 (utc-timestamp for 1970-01-01 08:00:12.133)", field)
		}

		return nil
	})

	govalidator.AddCustomRule("unix_timestamp", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		msInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be numeric", field)
		}

		// The value must be greater than unix-timestamp for 1970-01-01 08:00:12.133
		if msInt < 12 {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be right unix-timestamp(second) and greater than 12 (unix-timestamp for 1970-01-01 08:00:12.133)", field)
		}

		return nil
	})

	govalidator.AddCustomRule("allow_empty", func(field string, rule string, message string, value interface{}) error {
		return nil
	})

	govalidator.AddCustomRule("array", func(field string, rule string, message string, value interface{}) error {
		amUtil := new(ArrayMapUtil)
		if amUtil.KindOfData(value) != reflect.Array{
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must be a array", field)
		}
		return nil
	})

	govalidator.AddCustomRule("datetime", func(field string, rule string, message string, value interface{}) error {
		const datetimeRegex=`^\d\d\d\d-([0][1-9]|1[0-2])-([0][1-9]|[12][0-9]|3[01]) (00|[0][0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9]):([0-9]|[0-5][0-9])$`
		if datetimeStr,ok := value.(string);ok{
			if match, _ := regexp.MatchString(datetimeRegex, datetimeStr); match{
				return nil
			}
		}
		if message != "" {
			return errors.New(message)
		}
		return fmt.Errorf("The %s field must be yyyy-MM-dd HH:mm:ss format", field)
	})
}
