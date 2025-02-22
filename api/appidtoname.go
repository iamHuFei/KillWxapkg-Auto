package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 定义一个结构体来映射JSON数据
type ResponseData struct {
	Code int `json:"code"`
	Data struct {
		Nickname      string `json:"nickname"`
		Username      string `json:"username"`
		Description   string `json:"description"`
		Avatar        string `json:"avatar"`
		AppID         string `json:"appid"`
		UsesCount     string `json:"uses_count"`
		PrincipalName string `json:"principal_name"`
	} `json:"data"`
}

// 请求函数
func SendRequest(appID string) ([]byte, error) {
	// 目标URL, 对应appid查询小程序
	url := "https://kainy.cn/api/weapp/info/"

	// 创建POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("appid="+appID)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Host", "kainy.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	return ioutil.ReadAll(resp.Body)
}

//func main() {
//	// 调用请求函数
//	body, err := SendRequest("wx310cecc58fa78fc6")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// 输出状态码
//	fmt.Println("响应内容:", string(body))
//
//	// 解析JSON响应
//	var responseData ResponseData
//	if err := json.Unmarshal(body, &responseData); err != nil {
//		fmt.Println("解析JSON失败:", err)
//		return
//	}
//
//	// 提取name字段
//	fmt.Println("名称:", responseData.Data.Nickname)
//}
