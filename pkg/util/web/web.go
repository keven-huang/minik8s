package web

import (
	"fmt"
	"io"
	"minik8s/pkg/util/log"
	"net/http"
)

type RequestOptions struct {
	body      io.Reader
	prefix    string
	log       bool
	bodyBytes *[]byte
}

type Option func(*RequestOptions)

func WithBody(body io.Reader) Option {
	return func(opts *RequestOptions) {
		opts.body = body
	}
}

func WithPrefix(prefix string) Option {
	return func(opts *RequestOptions) {
		opts.prefix = prefix
	}
}

func WithLog(log bool) Option {
	return func(opts *RequestOptions) {
		opts.log = log
	}
}

func WithBodyBytes(bodyBytes *[]byte) Option {
	return func(opts *RequestOptions) {
		opts.bodyBytes = bodyBytes
	}
}

func SendHttpRequest(method string, url string, options ...Option) error {
	// 默认选项
	opts := &RequestOptions{
		body:      nil,
		prefix:    "[util.web] ",
		log:       false,
		bodyBytes: nil,
	}

	// 应用选项
	for _, opt := range options {
		opt(opts)
	}

	opts.prefix = opts.prefix + "[SendHttpRequest] "

	// 创建请求
	req, err := http.NewRequest(method, url, opts.body)
	if err != nil {
		fmt.Println(opts.prefix, err)
		return err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(opts.prefix, "Error sending request:", err)
		return err
	}

	// 打印响应结果
	err = log.CheckHttpStatus(opts.prefix, resp, log.WithLog(opts.log), log.WithBodyBytes(opts.bodyBytes))
	if err != nil {
		return err
	}

	return nil
}
