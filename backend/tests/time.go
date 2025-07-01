package main

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前时间
	t := time.Now()

	// 打印当前时间（需要格式化或转换为字符串）
	fmt.Println("当前时间:", t)

	// 或者使用格式化输出
	fmt.Printf("当前时间: %v\n", t)

	// 或者使用特定的时间格式
	fmt.Println("当前时间(格式化):", t.Format("2006-01-02 15:04:05"))
}
