package log

import (
	"fmt"
	"io"
	"net/http"
)

type Options struct {
	log       bool
	bodyBytes *[]byte
}

type OpFunc func(*Options)

func WithLog(log bool) OpFunc {
	return func(opts *Options) {
		opts.log = log
	}
}

func WithBodyBytes(bodyBytes *[]byte) OpFunc {
	return func(opts *Options) {
		opts.bodyBytes = bodyBytes
	}
}

func CheckHttpStatus(prefix string, resp *http.Response, opFuncs ...OpFunc) error {
	prefix = prefix + "[CheckHttpStatus] "

	opts := &Options{
		log:       false,
		bodyBytes: nil,
	}

	for _, opFunc := range opFuncs {
		opFunc(opts)
	}

	fmt.Println(prefix, "Response Status:", resp.Status)

	// 读取响应主体内容到字节数组
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(prefix, "Error reading response body:", err)
		return err
	}

	if opts.bodyBytes != nil {
		*opts.bodyBytes = bodyBytes
		//fmt.Println(prefix, string(*opts.bodyBytes))
	}

	if resp.StatusCode != 200 || opts.log {
		// 将字节数组转换为字符串并打印
		fmt.Println(prefix, "Response Body:", string(bodyBytes))
	}
	return nil
}
