package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Order 表示订单数据结构，与主服务保持一致
type Order struct {
	ID        string  `json:"id,omitempty"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status,omitempty"`
}

// 主函数提供简单的命令行交互界面
func main() {
	reader := bufio.NewReader(os.Stdin)
	baseURL := "http://localhost:8080"
	// baseURL := "http://localhost" // 配置Ingress、LoadBalancer 后

	for {
		printMenu()

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			checkHealth(baseURL)
		case "2":
			checkReady(baseURL)
		case "3":
			createNewOrder(baseURL, reader)
		case "4":
			listAllOrders(baseURL)
		case "5":
			getOrderByID(baseURL, reader)
		case "6":
			updateOrderByID(baseURL, reader)
		case "7":
			deleteOrderByID(baseURL, reader)
		case "8":
			fmt.Println("退出程序...")
			return
		default:
			fmt.Println("无效的选择，请重试。")
		}

		fmt.Println("\n按Enter键继续...")
		reader.ReadString('\n')
	}
}

// 打印菜单选项
func printMenu() {
	fmt.Println("======== 订单服务请求工具 ========")
	fmt.Println("1. 健康检查")
	fmt.Println("2. ready检查")
	fmt.Println("3. 创建订单")
	fmt.Println("4. 列出所有订单")
	fmt.Println("5. 根据ID获取订单")
	fmt.Println("6. 更新订单")
	fmt.Println("7. 删除订单")
	fmt.Println("8. 退出")
	fmt.Print("请选择操作 (1-8): ")
}

// 发送GET请求并显示响应
func sendGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	displayResponse(resp)
}

// 发送POST请求并显示响应
func sendPostRequest(url string, body interface{}) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	displayResponse(resp)
}

// 发送PUT请求并显示响应
func sendPutRequest(url string, body interface{}) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("JSON序列化失败: %v\n", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	displayResponse(resp)
}

// 发送DELETE请求并显示响应
func sendDeleteRequest(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	displayResponse(resp)
}

// 显示HTTP响应
func displayResponse(resp *http.Response) {
	fmt.Printf("HTTP状态码: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体失败: %v\n", err)
		return
	}

	// 格式化JSON响应以便于阅读
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		fmt.Println("响应内容:")
		fmt.Println(string(body))
	} else {
		fmt.Println("响应内容:")
		fmt.Println(prettyJSON.String())
	}
}

// 健康检查
func checkHealth(baseURL string) {
	fmt.Println("=== 健康检查 ===")
	sendGetRequest(baseURL + "/health")
}

// ready检查
func checkReady(baseURL string) {
	fmt.Println("=== ready检查 ===")
	sendGetRequest(baseURL + "/ready")
}

// 创建新订单
func createNewOrder(baseURL string, reader *bufio.Reader) {
	fmt.Println("=== 创建订单 ===")

	fmt.Print("请输入产品ID: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	fmt.Print("请输入数量: ")
	quantityStr, _ := reader.ReadString('\n')
	quantityStr = strings.TrimSpace(quantityStr)
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		fmt.Println("数量必须是数字")
		return
	}

	fmt.Print("请输入金额: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("金额必须是数字")
		return
	}

	order := Order{
		ProductID: productID,
		Quantity:  quantity,
		Amount:    amount,
	}

	sendPostRequest(baseURL+"/api/orders", order)
}

// 列出所有订单
func listAllOrders(baseURL string) {
	fmt.Println("=== 所有订单列表 ===")
	sendGetRequest(baseURL + "/api/orders")
}

// 根据ID获取订单
func getOrderByID(baseURL string, reader *bufio.Reader) {
	fmt.Println("=== 根据ID获取订单 ===")

	fmt.Print("请输入订单ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	sendGetRequest(baseURL + "/api/orders/" + id)
}

// 更新订单
func updateOrderByID(baseURL string, reader *bufio.Reader) {
	fmt.Println("=== 更新订单 ===")

	fmt.Print("请输入订单ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("请输入产品ID: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	fmt.Print("请输入数量: ")
	quantityStr, _ := reader.ReadString('\n')
	quantityStr = strings.TrimSpace(quantityStr)
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		fmt.Println("数量必须是数字")
		return
	}

	fmt.Print("请输入金额: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("金额必须是数字")
		return
	}

	fmt.Print("请输入订单状态: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	order := Order{
		ProductID: productID,
		Quantity:  quantity,
		Amount:    amount,
		Status:    status,
	}

	sendPutRequest(baseURL+"/api/orders/"+id, order)
}

// 删除订单
func deleteOrderByID(baseURL string, reader *bufio.Reader) {
	fmt.Println("=== 删除订单 ===")

	fmt.Print("请输入订单ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	sendDeleteRequest(baseURL + "/api/orders/" + id)
}
