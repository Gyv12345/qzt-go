// Package main 数据库检查工具
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// 连接 PostgreSQL（不指定数据库）
	connStr := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("连接 PostgreSQL 失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("PostgreSQL 未运行或无法连接: %v", err)
	}
	fmt.Println("✅ PostgreSQL 连接成功")

	// 检查数据库是否存在
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'qzt')`
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		log.Fatalf("查询数据库失败: %v", err)
	}

	if exists {
		fmt.Println("✅ 数据库 'qzt' 已存在")
	} else {
		fmt.Println("⚠️  数据库 'qzt' 不存在，正在创建...")
		_, err := db.Exec("CREATE DATABASE qzt WITH ENCODING='UTF8'")
		if err != nil {
			log.Fatalf("创建数据库失败: %v", err)
		}
		fmt.Println("✅ 数据库 'qzt' 创建成功")
	}

	// 连接到 qzt 数据库并检查表
	connStr = "host=localhost port=5432 user=postgres password=postgres dbname=qzt sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("连接 qzt 数据库失败: %v", err)
	}
	defer db.Close()

	// 检查表是否存在
	tables := []string{"users", "customers", "contacts", "contracts", "api_keys"}
	fmt.Println("\n📊 检查表结构:")
	for _, table := range tables {
		var tableExists bool
		query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = '%s')`, table)
		if err := db.QueryRow(query).Scan(&tableExists); err != nil {
			log.Printf("检查表 %s 失败: %v", table, err)
			continue
		}
		if tableExists {
			fmt.Printf("  ✅ %s\n", table)
		} else {
			fmt.Printf("  ⚠️  %s (将在启动时自动创建)\n", table)
		}
	}

	fmt.Println("\n✅ 数据库检查完成！")
}
