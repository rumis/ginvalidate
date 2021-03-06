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
		return tmpRes, 0, err
	}
	res, errCode, err := govalidate.Validate(tmpRes, rules)
	if err != nil {
		return tmpRes, 0, err
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

// BindJsonStructRaw 返回值为对象
// 如果校验失败，返回原始数据内容
// Content-type:application/json
func BindJsonStructRaw(c *gin.Context, rules []validator.FilterItem, obj interface{}) (int32, interface{}, error) {
	res, errCode, err := BindJsonMap(c, rules)
	if err != nil {
		return 0, res, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, res, err
	}
	return errCode, nil, nil
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
		return res, 0, err
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

// BindQueryStructRaw 解析Query参数
// 若解析失败，返回原始数据内容
func BindQueryStructRaw(c *gin.Context, rules []validator.FilterItem, obj interface{}) (int32, interface{}, error) {
	res, errCode, err := BindQueryMap(c, rules)
	if err != nil {
		return 0, res, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, res, err
	}
	return errCode, nil, nil
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
		return tmpRes, 0, err
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

// BindFormStructRaw 解析Form参数
// 若校验失败，返回map格式的原始数据
func BindFormStructRaw(c *gin.Context, rules []validator.FilterItem, obj interface{}) (int32, interface{}, error) {
	res, errCode, err := BindFormMap(c, rules)
	if err != nil {
		return 0, res, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, res, err
	}
	return errCode, nil, nil
}
