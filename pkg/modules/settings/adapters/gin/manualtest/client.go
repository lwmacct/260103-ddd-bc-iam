package manualtest

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/lwmacct/260101-go-pkg-gin/pkg/response"
)

// Client HTTP 测试客户端。
type Client struct {
	resty     *resty.Client
	devSecret string
	token     string
}

// newClient 创建测试客户端（内部使用）。
func newClient(baseURL string) *Client {
	r := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(10*time.Second).
		SetHeader("Content-Type", "application/json")

	return &Client{
		resty: r,
	}
}

// SetToken 设置访问令牌（使用 IAM 登录获取）。
func (c *Client) SetToken(token string) {
	c.token = token
	c.resty.SetAuthToken(token)
}

// Token 返回当前访问令牌。
func (c *Client) Token() string {
	return c.token
}

// R 返回 resty Request，用于自定义请求。
func (c *Client) R() *resty.Request {
	return c.resty.R()
}

// Get 发送 GET 请求并解析响应。
func Get[T any](c *Client, path string, queryParams map[string]string) (*T, error) {
	var result response.DataResponse[T]

	req := c.resty.R().SetResult(&result)
	if queryParams != nil {
		req.SetQueryParams(queryParams)
	}

	resp, err := req.Get(path)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}

// GetList 发送 GET 请求并解析列表响应。
func GetList[T any](c *Client, path string, queryParams map[string]string) ([]T, *response.PaginationMeta, error) {
	var result response.PagedResponse[T]

	req := c.resty.R().SetResult(&result)
	if queryParams != nil {
		req.SetQueryParams(queryParams)
	}

	resp, err := req.Get(path)
	if err != nil {
		return nil, nil, fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		return nil, nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return result.Data, result.Meta, nil
}

// Post 发送 POST 请求并解析响应。
func Post[T any](c *Client, path string, body any) (*T, error) {
	var result response.DataResponse[T]

	resp, err := c.resty.R().
		SetBody(body).
		SetResult(&result).
		Post(path)

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		// 如果响应体解析失败，使用原始响应体
		if result.Message == "" {
			return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), string(resp.Body()))
		}
		return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}

// Put 发送 PUT 请求并解析响应。
func Put[T any](c *Client, path string, body any) (*T, error) {
	var result response.DataResponse[T]

	resp, err := c.resty.R().
		SetBody(body).
		SetResult(&result).
		Put(path)

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}

// Patch 发送 PATCH 请求并解析响应。
func Patch[T any](c *Client, path string, body any) (*T, error) {
	var result response.DataResponse[T]

	resp, err := c.resty.R().
		SetBody(body).
		SetResult(&result).
		Patch(path)

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return &result.Data, nil
}

// Delete 发送 DELETE 请求。
func (c *Client) Delete(path string) error {
	var result response.DataResponse[any]

	resp, err := c.resty.R().
		SetResult(&result).
		Delete(path)

	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("状态码 %d: %s", resp.StatusCode(), result.Message)
	}

	return nil
}
