package ginvalidate

import (
	"context"
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
func BindJsonMap(c *gin.Context, rules []validator.Filter) (map[string]interface{}, int32, error) {
	defer c.Request.Body.Close()

	// 解析body
	decoder := json.NewDecoder(c.Request.Body)
	decoder.UseNumber()
	tmpRes := make(map[string]interface{})
	err := decoder.Decode(&tmpRes)
	if err != nil && !errors.Is(err, io.EOF) {
		return tmpRes, 0, err
	}
	// 解析header参数
	for k, v := range c.Request.Header {
		tmpRes[strings.ToLower(k)] = strings.Join(v, ",")
	}
	// 校验
	res, errCode, err := govalidate.Validate(tmpRes, rules)
	if err != nil {
		return tmpRes, 0, err
	}
	return res, errCode, nil
}

// BindJsonMapContent 解析请求参数
// Content-type:application/json
func BindJsonMapContext(c *gin.Context, rules []validator.Filter) (map[string]interface{}, int32, error) {
	defer c.Request.Body.Close()
	// 解析body参数
	decoder := json.NewDecoder(c.Request.Body)
	decoder.UseNumber()
	tmpRes := make(map[string]interface{})
	err := decoder.Decode(&tmpRes)
	if err != nil && !errors.Is(err, io.EOF) {
		return tmpRes, 0, err
	}
	// 解析header参数
	for k, v := range c.Request.Header {
		tmpRes[strings.ToLower(k)] = strings.Join(v, ",")
	}
	// 校验
	res, errCode, err := govalidate.Validate1(toContext(c), tmpRes, rules)
	if err != nil {
		return tmpRes, 0, err
	}
	return res, errCode, nil
}

// BindJsonStruct 返回值为对象
// Content-type:application/json
func BindJsonStruct(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, error) {
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

// BindJsonStructContent 返回值为对象
// Content-type:application/json
func BindJsonStructContext(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, error) {
	res, errCode, err := BindJsonMapContext(c, rules)
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
func BindJsonStructRaw(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, interface{}, error) {
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

// BindJsonStructRawContent 返回值为对象
// 如果校验失败，返回原始数据内容
// Content-type:application/json
func BindJsonStructRawContext(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, interface{}, error) {
	res, errCode, err := BindJsonMapContext(c, rules)
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
func BindQueryMap(c *gin.Context, rules []validator.Filter) (map[string]interface{}, int32, error) {
	pCol := NewParamsCollection()
	// 解析查询参数
	values := c.Request.URL.Query()
	for k, v := range values {
		k = FormatKey(k)
		pCol.Set(k, v)
	}
	// 解析header参数
	for k, v := range c.Request.Header {
		pCol.Set(strings.ToLower(k), v)
	}
	// 校验
	res, errCode, err := govalidate.Validate(pCol.To(), rules)
	if err != nil {
		return res, 0, err
	}
	return res, errCode, nil
}

// BindQueryMapContent 解析Query部分参数
func BindQueryMapContext(c *gin.Context, rules []validator.Filter) (map[string]interface{}, int32, error) {
	pCol := NewParamsCollection()
	// 解析查询参数
	values := c.Request.URL.Query()
	for k, v := range values {
		k = FormatKey(k)
		pCol.Set(k, v)
	}
	// 解析Header参数
	for k, v := range c.Request.Header {
		pCol.Set(strings.ToLower(k), v)
	}
	// 校验
	res, errCode, err := govalidate.Validate1(toContext(c), pCol.To(), rules)
	if err != nil {
		return res, 0, err
	}
	return res, errCode, nil
}

// BindQueryStruct 解析Query参数
func BindQueryStruct(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, error) {
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

// BindQueryMapContent 解析Query参数
func BindQueryStructContext(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, error) {
	res, errCode, err := BindQueryMapContext(c, rules)
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
func BindQueryStructRaw(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, interface{}, error) {
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

// BindQueryStructRawContent 解析Query参数
// 若解析失败，返回原始数据内容
func BindQueryStructRawContext(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, interface{}, error) {
	res, errCode, err := BindQueryMapContext(c, rules)
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
func BindFormMap(c *gin.Context, rules []validator.Filter) (map[string]interface{}, int32, error) {
	if err := c.Request.ParseForm(); err != nil {
		return nil, 0, err
	}
	pCol := NewParamsCollection()
	values := c.Request.PostForm
	for k, v := range values {
		k = FormatKey(k)
		pCol.Set(k, v)
	}

	// 解析MultipartForm
	err := c.Request.ParseMultipartForm(10240)
	if err == nil {
		for k, v := range c.Request.MultipartForm.Value {
			k = FormatKey(k)
			pCol.Set(k, v)
		}
	}

	// 解析Header参数
	for k, v := range c.Request.Header {
		pCol.Set(strings.ToLower(k), v)
	}

	res, errCode, err := govalidate.Validate(pCol.To(), rules)
	if err != nil {
		return pCol.To(), 0, err
	}
	return res, errCode, nil
}

// BindFormMapContent 解析form数据
func BindFormMapContext(c *gin.Context, rules []validator.Filter) (map[string]interface{}, int32, error) {
	if err := c.Request.ParseForm(); err != nil {
		return nil, 0, err
	}
	// 解析form
	pCol := NewParamsCollection()
	values := c.Request.PostForm
	for k, v := range values {
		k = FormatKey(k)
		pCol.Set(k, v)
	}
	// 解析MultipartForm
	err := c.Request.ParseMultipartForm(10240)
	if err == nil {
		for k, v := range c.Request.MultipartForm.Value {
			k = FormatKey(k)
			pCol.Set(k, v)
		}
	}
	// 解析Header参数
	for k, v := range c.Request.Header {
		pCol.Set(strings.ToLower(k), v)
	}
	// 校验
	res, errCode, err := govalidate.Validate1(toContext(c), pCol.To(), rules)
	if err != nil {
		return pCol.To(), 0, err
	}
	return res, errCode, nil
}

// BindFormStruct 解析Form参数
func BindFormStruct(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, error) {
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

// BindFormStructContent 解析Form参数
func BindFormStructContext(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, error) {
	res, errCode, err := BindFormMapContext(c, rules)
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
func BindFormStructRaw(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, interface{}, error) {
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

// BindFormStructRawContent 解析Form参数
// 若校验失败，返回map格式的原始数据
func BindFormStructRawContext(c *gin.Context, rules []validator.Filter, obj interface{}) (int32, interface{}, error) {
	res, errCode, err := BindFormMapContext(c, rules)
	if err != nil {
		return 0, res, err
	}
	err = mapDecode(res, obj)
	if err != nil {
		return 0, res, err
	}
	return errCode, nil, nil
}

func toContext(c *gin.Context) context.Context {
	stdCtx := context.Background()
	for k, v := range c.Keys {
		stdCtx = context.WithValue(stdCtx, k, v)
	}
	return stdCtx
}
