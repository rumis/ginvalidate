package ginvalidate

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rumis/govalidate"
	"github.com/rumis/govalidate/validator"
)

// BindJsonMap 解析请求参数
// Content-type:application/json
func BindJsonMap(c *gin.Context, rules []validator.FilterItem) (map[string]interface{}, int32, error) {
	defer c.Request.Body.Close()

	decoder := json.NewDecoder(c.Request.Body)
	decoder.UseNumber()
	tmpRes := make(map[string]interface{})
	err := decoder.Decode(&tmpRes)
	if err != nil && !errors.Is(err, io.EOF) {
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
func BindJsonStruct(c *gin.Context, rules []validator.FilterItem, obj interface{}) (int32, error) {
	res, errCode, err := BindJsonMap(c, rules)
	if err != nil {
		return 0, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, err
	}
	return errCode, nil
}

// BindQueryMap 解析Query部分参数
func BindQueryMap(c *gin.Context, rules []validator.FilterItem) (map[string]interface{}, int32, error) {
	tmpRes := make(map[string]interface{})
	values := c.Request.URL.Query()
	for k, v := range values {
		k = strings.TrimRight(k, "[]")
		if len(v) == 1 {
			tmpRes[k] = v[0]
			continue
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
func BindQueryStruct(c *gin.Context, rules []validator.FilterItem, obj interface{}) (int32, error) {
	res, errCode, err := BindQueryMap(c, rules)
	if err != nil {
		return 0, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, err
	}
	return errCode, nil
}

// BindFormMap 解析form数据
func BindFormMap(c *gin.Context, rules []validator.FilterItem) (map[string]interface{}, int32, error) {
	if err := c.Request.ParseForm(); err != nil {
		return nil, 0, err
	}
	tmpRes := make(map[string]interface{})
	values := c.Request.PostForm
	for k, v := range values {
		k = strings.TrimRight(k, "[]")
		if len(v) == 1 {
			tmpRes[k] = v[0]
			continue
		}
		tmpRes[k] = v
	}
	res, errCode, err := govalidate.Validate(tmpRes, rules)
	if err != nil {
		return nil, 0, err
	}
	return res, errCode, nil
}

// BindFormStruct 解析Form参数
func BindFormStruct(c *gin.Context, rules []validator.FilterItem, obj interface{}) (int32, error) {
	res, errCode, err := BindFormMap(c, rules)
	if err != nil {
		return 0, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, err
	}
	return errCode, nil
}
