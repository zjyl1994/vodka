package vodka

import (
	"reflect"
	"strings"
)

type ArrayMapUtil struct{}

// return keys for map with interface val.
func (amu *ArrayMapUtil) Keys(m map[string]interface{}) []string {
	keys := []string{}
	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}

func (amu *ArrayMapUtil) Int64Keys(m map[int64]interface{}) []int64 {
	keys := []int64{}
	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}

// return keys for map with array val.
func (amu *ArrayMapUtil) KeysA(m map[string][]string) []string {
	keys := []string{}
	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}

func (amu *ArrayMapUtil) IsMapEmpty(m map[string]interface{}) bool {
	return len(amu.Keys(m)) == 0
}

func (amu *ArrayMapUtil) IsMapIncludeKey(m map[string]interface{}, key string) bool {
	flag := false

	for k, _ := range m {
		if k == key {
			flag = true
		}
	}

	return flag
}

func (amu *ArrayMapUtil) IsArrayInclude(a []string, val string) bool {
	flag := false

	for _, v := range a {
		if v == val {
			flag = true
		}
	}

	return flag
}

func (amu *ArrayMapUtil) IsArrayIncludeInt64(a []int64, val int64) bool {
	flag := false

	for _, v := range a {
		if v == val {
			flag = true
		}
	}

	return flag
}

func (amu *ArrayMapUtil) GetStrWithSingleQuoteForArray(a []string) string {
	temp_a := make([]string, len(a))
	for i, v := range a {
		temp_a[i] = strings.Join([]string{"'", v, "'"}, "")
	}

	return strings.Join(temp_a, ",")
}

/**
  e.g: in: ["aaa,bbb", "ccc,ddd", "eee"], out: ["aaa", "bbb", "ccc", "ddd", "eee"].
*/
func (amu *ArrayMapUtil) TransArrayWithCommaToFullArray(a []string) []string {
	temp_a := []string{}

	for _, v := range a {
		temp_b := strings.Split(v, ",")
		for _, v1 := range temp_b {
			temp_a = append(temp_a, v1)
		}
	}

	return temp_a
}

func (amu *ArrayMapUtil) IsArrayHasDuplicateItems(a []string, val string) bool {
	count := 0

	for _, v := range a {
		if v == val {
			count++
		}
	}

	return count > 1
}

func (amu *ArrayMapUtil) RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	for _, v1 := range a {
		if len(ret) == 0 {
			ret = append(ret, v1)
		} else {
			for k2, v2 := range ret {
				if v1 == v2 || v1 == "" {
					break
				} else if k2 == len(ret)-1 {
					ret = append(ret, v1)
				}
			}
		}
	}
	return
}

func (amu *ArrayMapUtil) IsMapKeyOutRange(m map[string]interface{}, rangeKeysArr []string) bool {
	flag := false

loop:
	for key, _ := range m {
		if !amu.IsArrayInclude(rangeKeysArr, key) {
			flag = true
			break loop
		}
	}

	return flag
}

func (amu *ArrayMapUtil) KindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (amu *ArrayMapUtil) FilterOutRangeFields(m map[string]interface{}, rangeKeysArr []string) {
	for key, _ := range m {
		if !amu.IsArrayInclude(rangeKeysArr, key) {
			delete(m, key)
		}
	}
}
