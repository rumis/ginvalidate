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
	"github.com/hetiansu5/urlquery"
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

var rules = []V.FilterItem{
	R.NewFilter("name", []V.Validator{V.Required()}),
	R.NewFilter("ids", []V.Validator{V.Required(), V.DotInt(), V.Dotint2Slice(), V.IntSlice([]E.IntExecutor{E.Between(1, 100)})}),
	R.NewFilter("grade", []V.Validator{V.Required(), V.Int(), V.Between(1, 100)}),
	R.NewFilter("subjects", []V.Validator{V.Required(), V.IntSlice()}),
	R.NewFilter("ctime", []V.Validator{V.Required(), V.Datetime()}),
	R.NewFilter("email", []V.Validator{V.Required(), V.Email()}),
	R.NewFilter("phone", []V.Validator{V.Required(), V.Phone()}),
	R.NewFilter("stat", []V.Validator{V.Required(), V.EnumInt([]int{1, 2, 3, 4, 5})}),
	R.NewFilter("school", []V.Validator{V.Required(), V.Int()}),
	R.NewFilter("cname", []V.Validator{V.Required(), V.StringSlice()}),
	R.NewFilter("page", []V.Validator{V.Optional(101), V.Int()}),
}

func init() {
	router = gin.Default()

	router.POST("/json", jsonHandler)
	router.POST("/jsonraw", jsonRawErrorHandler)
	router.POST("/form", formHandler)
	router.POST("/query", queryHandler)
}

func TestBindJSON(t *testing.T) {

	s1 := map[string]interface{}{
		"name":     "??????",
		"ids":      "1,2,3",
		"grade":    2,
		"subjects": []int{3, 4, 12},
		"ctime":    time.Now().Format("2006-01-02 15:04:05"),
		"email":    "liumurong1@tal.com",
		"phone":    "15810562936",
		"stat":     3,
		"school":   1,
		"cname":    []string{"a", "b", "c"},
	}
	s1Byte, _ := json.Marshal(s1)

	req := httptest.NewRequest("POST", "/json", bytes.NewReader(s1Byte))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// ???????????????handler??????
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
	if out.Name != "??????" {
		t.Error("error")
	}
	if out.Ids[2] != 3 {
		t.Error("dotint to int slice error")
	}

	if out.Page != 101 {
		t.Error("optional default value error")
	}

}

// ???????????????????????????????????????
func TestBindJSONRaw(t *testing.T) {

	s1 := map[string]interface{}{
		"ids":      "1,2,3",
		"grade":    2,
		"subjects": []int{3, 4, 12},
		"ctime":    time.Now().Format("2006-01-02 15:04:05"),
		"email":    "liumurong1@tal.com",
		"phone":    "15810562936",
		"stat":     3,
		"school":   1,
		"cname":    []string{"a", "b", "c"},
	}
	s1Byte, _ := json.Marshal(s1)

	req := httptest.NewRequest("POST", "/jsonraw", bytes.NewReader(s1Byte))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	// ???????????????handler??????
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("code error: %d", w.Code)
	}

	res := w.Result()
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var resp Resp
	json.Unmarshal(body, &resp)
	result, _ := resp.Data.(map[string]interface{})

	if email, ok := result["email"]; !ok {
		t.Error("validate error response column [email] not exist")
		if es, ok := email.(string); !ok || es != "liumurong1@tal.com" {
			t.Error("validate error response column [email] type error")
		}
	}
	if ids, ok := result["ids"]; !ok {
		t.Error("validate error response column [ids] not exist")
		if eids, ok := ids.(string); !ok || eids != "liumurong1@tal.com" {
			t.Error("validate error response column [ids] type error")
		}
	}
}

func TestFormJSON(t *testing.T) {

	s1 := map[string]interface{}{
		"name":     "??????",
		"ids":      "1,2,3",
		"grade":    2,
		"subjects": []int{3, 4, 12},
		"ctime":    time.Now().Format("2006-01-02 15:04:05"),
		"email":    "liumurong1@tal.com",
		"phone":    "15810562936",
		"stat":     3,
		"school":   1,
		"cname":    []string{"a", "b", "c"},
	}

	// vals, _ := query.Values(s1)
	bs, err := urlquery.Marshal(s1)
	if err != nil {
		t.Error("map convent to query string error")
	}

	req := httptest.NewRequest("POST", "/form", strings.NewReader(string(bs)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	// ???????????????handler??????
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
	if out.Name != "??????" {
		t.Error("error")
	}
	if out.Ids[2] != 3 {
		t.Error("dotint to int slice error")
	}

}

func TestQueryJSON(t *testing.T) {

	s1 := map[string]interface{}{
		"name":     "??????",
		"ids":      "1,2,3",
		"grade":    2,
		"subjects": []int{3, 4, 12},
		"ctime":    time.Now().Format("2006-01-02 15:04:05"),
		"email":    "liumurong1@tal.com",
		"phone":    "15810562936",
		"stat":     3,
		"school":   1,
		"cname":    []string{"a", "b", "c"},
	}

	str, _ := urlquery.Marshal(s1)

	req := httptest.NewRequest("POST", "/query?"+string(str), nil)

	w := httptest.NewRecorder()
	// ???????????????handler??????
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
	if out.Name != "??????" {
		t.Error("error")
	}
	if out.Ids[2] != 3 {
		t.Error("dotint to int slice error")
	}

}

// JSON??????
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

// JSON??????
// ???????????????????????????????????????
func jsonRawErrorHandler(c *gin.Context) {
	var s SlideResp
	code, raw, err := BindJsonStructRaw(c, rules, &s)
	if err != nil {
		c.JSON(http.StatusOK, Resp{
			Code: int(code),
			Msg:  err.Error(),
			Data: raw,
		})
		return
	}
	c.JSON(http.StatusOK, Resp{
		Code: int(code),
		Msg:  "",
		Data: s,
	})
}

// Query??????
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

// ????????????
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
