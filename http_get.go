package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// HttpGetContent 提供获取http内容的能力
func HttpGetContent(fullUrl string, headers map[string]string, params map[string]string) ([]byte, error) {
	// 创建查询参数
	totalUrl := url.Values{}
	for k, v := range params {
		totalUrl.Add(k, v)
	}

	u, err := url.Parse(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("parsing url failed: %s", err.Error())
	}
	// 设置查询参数
	u.RawQuery = totalUrl.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new http get request failed: %s", err.Error())
	}
	// 添加请求头
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// 发送请求
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http get failed: %s", err.Error())
	}
	defer rsp.Body.Close()

	// 保存请求结果
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %s", err.Error())
	}

	if rsp.StatusCode != http.StatusOK { // 若返回状态有问题，抛出响应内容
		return nil, fmt.Errorf("server returned %v status，body is %s", rsp.StatusCode, string(body))
	}
	return body, nil
}

// HttpGetHtml 提供获取html的能力
func HttpGetHtml(fullUrl string, headers map[string]string, params map[string]string, coder func(body io.ReadCloser) io.Reader) (string, error) {
	// 创建查询参数
	totalUrl := url.Values{}
	for k, v := range params {
		totalUrl.Add(k, v)
	}

	u, err := url.Parse(fullUrl)
	if err != nil {
		return "", fmt.Errorf("parsing url failed: %s", err.Error())
	}
	// 设置查询参数
	u.RawQuery = totalUrl.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("new http get request failed: %s", err.Error())
	}
	// 添加请求头
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// 发送请求
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do http get failed: %s", err.Error())
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned %v status", rsp.StatusCode)
	}

	var rspBody io.Reader = rsp.Body
	if coder != nil {
		rspBody = coder(rsp.Body)
	}

	// 解析html
	tokenizer := html.NewTokenizer(rspBody)
	var buf strings.Builder
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}
		if tokenType == html.TextToken {
			text := strings.TrimSpace(string(tokenizer.Text()))
			if text != "" {
				buf.WriteString(text + "\n") // 添加换行符以便后续分割
			}
		}
	}
	return buf.String(), nil
}

func GB180302UTF8(body io.ReadCloser) io.Reader {
	//decoder := encoding.Nop.NewDecoder()
	decoder := simplifiedchinese.GB18030.NewDecoder()
	return transform.NewReader(body, decoder)
}

func GBK2UTF8(body io.ReadCloser) io.Reader {
	decoder := simplifiedchinese.GBK.NewDecoder()
	return transform.NewReader(body, decoder)
}

func Windows12522UTF8(body io.ReadCloser) io.Reader {
	decoder := charmap.Windows1252.NewDecoder()
	return transform.NewReader(body, decoder)
}
