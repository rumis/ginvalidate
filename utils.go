package ginvalidate

import (
	"regexp"

	"github.com/mitchellh/mapstructure"
)

// mapDecode map转对象
func mapDecode(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
			mapstructure.StringToTimeHookFunc("2006-01-02")),
		Metadata: nil,
		Result:   out,
		TagName:  "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	err = decoder.Decode(input)
	if err != nil {
		return err
	}
	return nil
}

// FormatKey 去除数组类参数KEY中的中括号
func FormatKey(k string) string {
	reg := regexp.MustCompile(`\[\d*\]`)
	k = reg.ReplaceAllString(k, "")
	return k
}
