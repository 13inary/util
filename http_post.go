package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HttpPost 提供post请求的能力
func HttpPost(fullUrl string, headers map[string]string, formDatas map[string]string) ([]byte, error) {
	// 创建查询参数
	data := url.Values{}
	for k, v := range formDatas {
		data.Set(k, v)
	}

	// 使用strings.NewReader(data.Encode())将表单数据转换为io.Reader
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("new http post request failed: %s", err.Error())
	}

	// 设置Content-Type头部为application/x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// 发送请求
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http get failed: %s", err.Error())
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %v status", rsp.StatusCode)
	}

	// 保存请求结果
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err.Error())
	}
	return body, nil
}
