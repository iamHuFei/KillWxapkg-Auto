package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 定义一个结构体来映射JSON数据
type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AppID       string   `json:"appid"`
		Name        string   `json:"name"`
		Logo        string   `json:"logo"`
		Description string   `json:"description"`
		Category    []string `json:"category"`
		Username    string   `json:"username"`
		Verify      struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"verify"`
		AppUpdateTime string `json:"app_update_time"`
	} `json:"data"`
}

// 请求函数
func SendRequest(appID string) ([]byte, error) {
	// 目标URL, 对应appid查询小程序：网络游客工具箱
	url := "https://www.webstr.top/tool/jx_appid/app_getData?v=1.2&from=wxapp&sessionkey="

	// 创建POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("appid="+appID)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Host", "www.webstr.top")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36 MicroMessenger/7.0.20.1781(0x6700143B) NetType/WIFI MiniProgramEnv/Windows WindowsWechat/WMPF WindowsWechat(0x63090c11)XWEB/11275")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://servicewechat.com/wx9d4af941caa4defe/7/page-frame.html")

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
//	body, err := SendRequest("wx2e6588fd778fa5c6")
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
//	fmt.Println("名称:", responseData.Data.Name)
//}
