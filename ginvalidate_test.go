package ginvalidate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	"github.com/mitchellh/mapstructure"
	R "github.com/rumis/govalidate"
	E "github.com/rumis/govalidate/executor"
	V "github.com/rumis/govalidate/validator"
)

type SlideReq struct {
	Name     string   `json:"name" url:"name"`
	Ids      string   `json:"ids" url:"ids"`
	Grade    int      `json:"grade" url:"grade"`
	Subjects []int    `json:"subjects" url:"subjects"`
	Ctime    string   `json:"ctime" url:"ctime"`
	Email    string   `json:"email" url:"email"`
	Phone    string   `json:"phone" url:"phone"`
	Stat     int      `json:"stat" url:"stat"`
	School   int      `json:"school" url:"school"`
	Cname    []string `json:"cname" url:"cname"`
	Page     int      `json:"page" url:"page"`
}

type SlideResp struct {
	Name     string    `json:"name"`
	Ids      []int     `json:"ids"`
	Grade    int       `json:"grade"`
	Subjects []int     `json:"subjects"`
	Ctime    time.Time `json:"ctime"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Stat     int       `json:"stat"`
	School   int       `json:"school"`
	Cname    []string  `json:"cname"`
	Page     int       `json:"page"`
}

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var router *gin.Engine

var rules = map[string]V.FilterItem{
	"name":     R.Filter([]V.Validator{V.Required()}),
	"ids":      R.Filter([]V.Validator{V.Required(), V.DotInt(), V.Dotint2Slice(), V.IntSlice([]E.IntExecutor{E.Between(1, 100)})}),
	"grade":    R.Filter([]V.Validator{V.Required(), V.Int(), V.Between(1, 100)}),
	"subjects": R.Filter([]V.Validator{V.Required(), V.IntSlice()}),
	"ctime":    R.Filter([]V.Validator{V.Required(), V.Datetime()}),
	"email":    R.Filter([]V.Validator{V.Required(), V.Email()}),
	"phone":    R.Filter([]V.Validator{V.Required(), V.Phone()}),
	"stat":     R.Filter([]V.Validator{V.Required(), V.EnumInt([]int{1, 2, 3, 4, 5})}),
	"school":   R.Filter([]V.Validator{V.Required(), V.Int()}),
	"cname":    R.Filter([]V.Validator{V.Required(), V.StringSlice()}),
	"page":     R.Filter([]V.Validator{V.Optional(101), V.Int()}),
}

func init() {
	router = gin.Default()

	router.POST("/json", jsonHandler)
	router.POST("/form", formHandler)
	router.POST("/query", queryHandler)
}

func TestBindJSON(t *testing.T) {

	s1 := SlideReq{
		Name:     "课件",
		Ids:      "1,2,3",
		Grade:    2,
		Subjects: []int{3, 4, 12},
		Ctime:    time.Now().Format("2006-01-02 15:04:05"),
		Email:    "liumurong1@tal.com",
		Phone:    "15810562936",
		Stat:     3,
		School:   1,
		Cname:    []string{"a", "b", "c"},
	}
	s1Byte, _ := json.Marshal(s1)

	req := httptest.NewRequest("POST", "/json", bytes.NewReader(s1Byte))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("code error: %d", w.Code)
	}

	res := w.Result()
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp Resp
	json.Unmarshal(body, &resp)
	var out SlideResp
	mapstructure.Decode(resp.Data, &out)
	if out.Name != s1.Name {
		t.Error("error")
	}
	if out.Ids[2] != 3 {
		t.Error("dotint to int slice error")
	}

	// if out.Page != 101 {
	// 	t.Error("optional default value error")
	// }

}

func TestFormJSON(t *testing.T) {

	s1 := SlideReq{
		Name:     "课件",
		Ids:      "1,2,3",
		Grade:    2,
		Subjects: []int{3, 4, 12},
		Ctime:    time.Now().Format("2006-01-02 15:04:05"),
		Email:    "liumurong1@tal.com",
		Phone:    "15810562936",
		Stat:     3,
		School:   1,
		Cname:    []string{"a", "b", "c"},
	}

	vals, _ := query.Values(s1)
	str := vals.Encode()

	req := httptest.NewRequest("POST", "/form", strings.NewReader(str))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("code error: %d", w.Code)
	}

	res := w.Result()
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp Resp
	json.Unmarshal(body, &resp)
	var out SlideResp
	mapstructure.Decode(resp.Data, &out)
	if out.Name != s1.Name {
		t.Error("error")
	}
	if out.Ids[2] != 3 {
		t.Error("dotint to int slice error")
	}

}

func TestQueryJSON(t *testing.T) {

	s1 := SlideReq{
		Name:     "课件",
		Ids:      "1,2,3",
		Grade:    2,
		Subjects: []int{3, 4, 12},
		Ctime:    time.Now().Format("2006-01-02 15:04:05"),
		Email:    "liumurong1@tal.com",
		Phone:    "15810562936",
		Stat:     3,
		School:   1,
		Cname:    []string{"a", "b", "c"},
	}

	vals, _ := query.Values(s1)
	str := vals.Encode()

	req := httptest.NewRequest("POST", "/query?"+str, nil)

	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("code error: %d", w.Code)
	}

	res := w.Result()
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp Resp
	json.Unmarshal(body, &resp)
	var out SlideResp
	mapstructure.Decode(resp.Data, &out)
	if out.Name != s1.Name {
		t.Error("error")
	}
	if out.Ids[2] != 3 {
		t.Error("dotint to int slice error")
	}

}

// JSON传参
func jsonHandler(c *gin.Context) {
	var s SlideResp
	code, err := BindJsonStruct(c, rules, &s)
	if err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: int(code),
			Msg:  err.Error(),
			Data: gin.H{},
		})
		return
	}
	c.JSON(http.StatusOK, Resp{
		Code: int(code),
		Msg:  "",
		Data: s,
	})
}

// Query传参
func queryHandler(c *gin.Context) {
	var s SlideResp
	code, err := BindQueryStruct(c, rules, &s)
	if err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: int(code),
			Msg:  err.Error(),
			Data: gin.H{},
		})
		return
	}
	c.JSON(http.StatusOK, Resp{
		Code: int(code),
		Msg:  "",
		Data: s,
	})
}

// 简单表单
func formHandler(c *gin.Context) {
	var s SlideResp
	code, err := BindFormStruct(c, rules, &s)
	if err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: int(code),
			Msg:  err.Error(),
			Data: gin.H{},
		})
		return
	}
	c.JSON(http.StatusOK, Resp{
		Code: int(code),
		Msg:  "",
		Data: s,
	})
}
