package main

import (
	"fmt"

	"github.com/werbenhu/bptree"
)

func main() {
	// 创建B+树实例
	tree := bptree.NewTree()

	// 插入测试数据
	for i := 1; i <= 16; i++ {
		tree.Insert(i, fmt.Sprintf("值%d", i))
	}

	// 打印树结构
	fmt.Println("B+树结构：")
	tree.Print()
}
