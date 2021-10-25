package ginvalidate

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/rumis/govalidate"
)

// BindJsonMap 解析请求参数
// Content-type:application/json
func BindJsonMap(c *gin.Context, rules map[string]govalidate.FilterItem) (map[string]interface{}, int32, error) {
	if c.ContentType() != "application/json" {
		return nil, 0, errors.New("不支持的Content-Type")
	}
	defer c.Request.Body.Close()
	decoder := json.NewDecoder(c.Request.Body)
	decoder.UseNumber()
	tmpRes := make(map[string]interface{})
	err := decoder.Decode(&tmpRes)
	if err != nil {
		return nil, 0, err
	}
	res, errCode, err := govalidate.Validate(tmpRes, rules)
	if err != nil {
		return nil, 0, err
	}
	return res, errCode, nil
}

// BindJsonStruct 返回值为对象
// Content-type:application/json
func BindJsonStruct(c *gin.Context, rules map[string]govalidate.FilterItem, obj interface{}) (int32, error) {
	res, errCode, err := BindJsonMap(c, rules)
	if err != nil {
		return 0, err
	}
	err = mapstructure.Decode(res, obj)
	if err != nil {
		return 0, err
	}
	return errCode, nil
}

// BindQueryMap 解析Query部分参数
func BindQueryMap(c *gin.Context, rules map[string]govalidate.FilterItem) (map[string]interface{}, int32, error) {
	tmpRes := make(map[string]interface{})
	values := c.Request.URL.Query()
	for k, v := range values {
		if len(v) == 1 {
			tmpRes[k] = v[0]
		}
		tmpRes[k] = v
	}
	res, errCode, err := govalidate.Validate(tmpRes, rules)
	if err != nil {
		return nil, 0, err
	}
	return res, errCode, nil
}

// BindQueryStruct 解析Query参数
func BindQueryStruct(c *gin.Context, rules map[string]govalidate.FilterItem, obj interface{}) (int32, error) {
	res, errCode, err := BindQueryMap(c, rules)
	if err != nil {
		return 0, err
	}
	err = mapstructure.Decode(res, obj)
	if err != nil {
		return 0, err
	}
	return errCode, nil
}