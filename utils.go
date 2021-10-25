package ginvalidate

import "github.com/mitchellh/mapstructure"

// mapDecode map转对象
func mapDecode(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
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
