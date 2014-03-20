package weixin

import (
    "time"
    "strconv"
)

func Int64Assert(val interface{}) int64 {
	switch val.(type) {
	case uint:
		return int64(val.(uint))
	case int:
		return int64(val.(int))
	case int64:
		return val.(int64)
	case int32:
		return int64(val.(int32))
	case int16:
		return int64(val.(int16))
	case int8:
		return int64(val.(int8))
	case uint8:
		return int64(val.(uint8))
	case uint16:
		return int64(val.(uint16))
	case uint32:
		return int64(val.(uint32))
	case uint64:
		return int64(val.(uint64))
	case float64:
		return int64(val.(float64))
	case float32:
		return int64(val.(float32))
	case nil:
		return int64(0)
	case string:
		n, e := strconv.ParseInt(val.(string), 10, 64)
		if e != nil {
			return int64(0)
		}
		return int64(n)
	case bool:
		switch val.(bool) {
		case true:
			return int64(1)
		case false:
			return int64(0)
		}
	}
	return val.(int64)
}

func ParseCreateTime(createTime interface{}) time.Time {
    c := Int64Assert(createTime)
    return time.Unix(c, 0)
}
