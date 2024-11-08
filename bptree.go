package bptree

import (
	"fmt"
	"sort"
)

// 定义B+树的阶数
const M = 4

// Tree B+树结构体
type Tree struct {
	root *Node // 根节点
	head *Node // 最左侧叶子节点（用于范围查询）
}

// Node B+树节点结构体
type Node struct {
	keys   []int    // 存储的键列表
	kids   []*Node  // 子节点指针列表（内部节点使用）
	vals   []string // 存储的值列表（叶子节点使用）
	leaf   bool     // 是否为叶子节点
	next   *Node    // 指向下一个叶子节点（叶子节点使用）
	parent *Node    // 父节点指针
}

// newNode 创建新节点
func newNode() *Node {
	return &Node{
		keys: make([]int, 0),
		kids: make([]*Node, 0),
		vals: make([]string, 0),
		leaf: true,
	}
}

// NewTree 创建新的B+树
func NewTree() *Tree {
	root := newNode()
	return &Tree{
		root: root,
		head: root,
	}
}

// Insert 插入键值对到B+树
func (t *Tree) Insert(key int, val string) {
	cur := t.root
	// 查找适合插入的叶子节点
	for !cur.leaf {
		idx := sort.SearchInts(cur.keys, key)
		cur = cur.kids[idx]
	}
	t.insertLeaf(cur, key, val)
}

// insertLeaf 向叶子节点插入键值对
func (t *Tree) insertLeaf(n *Node, key int, val string) {
	pos := sort.SearchInts(n.keys, key)

	// 在正确位置插入新的键值对
	if pos == len(n.keys) {
		// 插入到末尾
		n.keys = append(n.keys, key)
		n.vals = append(n.vals, val)
	} else {
		// 插入到中间位置
		n.keys = append(n.keys, 0)
		n.vals = append(n.vals, "")
		copy(n.keys[pos+1:], n.keys[pos:])
		copy(n.vals[pos+1:], n.vals[pos:])
		n.keys[pos] = key
		n.vals[pos] = val
	}

	// 节点已满，需要分裂
	if len(n.keys) >= M {
		t.splitLeaf(n)
	}
}

// splitLeaf 分裂叶子节点
func (t *Tree) splitLeaf(n *Node) {
	// 分裂点
	mid := M / 2

	// 创建新的右侧节点
	right := &Node{
		keys:   n.keys[mid:],
		vals:   n.vals[mid:],
		leaf:   true,
		next:   n.next,
		parent: n.parent,
	}

	// 更新左侧节点
	n.keys = n.keys[:mid]
	n.vals = n.vals[:mid]
	n.next = right

	// 处理根节点分裂的特殊情况
	if n.parent == nil {
		// 创建新的根节点
		root := &Node{
			keys: []int{right.keys[0]},
			kids: []*Node{n, right},
			leaf: false,
		}
		n.parent = root
		right.parent = root
		t.root = root
		return
	}

	// 将新节点插入到父节点
	t.insertParent(n, right.keys[0], right)
}

// insertParent 将节点插入到父节点
func (t *Tree) insertParent(left *Node, key int, right *Node) {
	p := left.parent
	pos := sort.SearchInts(p.keys, key)

	// 在父节点中插入新键和子节点
	if pos == len(p.keys) {
		// 插入到末尾
		p.keys = append(p.keys, key)
		p.kids = append(p.kids, right)
	} else {
		// 插入到中间位置
		p.keys = append(p.keys, 0)
		p.kids = append(p.kids, nil)
		copy(p.keys[pos+1:], p.keys[pos:])
		copy(p.kids[pos+1:], p.kids[pos:])
		p.keys[pos] = key
		p.kids[pos] = right
	}

	// 父节点已满，需要分裂
	if len(p.keys) >= M {
		t.splitInternal(p)
	}
}

// splitInternal 分裂内部节点
func (t *Tree) splitInternal(n *Node) {
	// 分裂点
	mid := M / 2
	key := n.keys[mid]

	// 创建新的右侧节点
	right := &Node{
		keys:   n.keys[mid+1:],
		kids:   n.kids[mid+1:],
		leaf:   false,
		parent: n.parent,
	}

	// 更新左侧节点
	n.keys = n.keys[:mid]
	n.kids = n.kids[:mid+1]

	// 更新子节点的父指针
	for _, child := range right.kids {
		child.parent = right
	}

	// 处理根节点分裂的特殊情况
	if n.parent == nil {
		// 创建新的根节点
		root := &Node{
			keys: []int{key},
			kids: []*Node{n, right},
			leaf: false,
		}
		n.parent = root
		right.parent = root
		t.root = root
		return
	}

	// 将新节点插入到父节点
	t.insertParent(n, key, right)
}

// Print 打印整个B+树结构（用于调试）
func (t *Tree) Print() {
	t.printNode(t.root, 0)
}

// printNode 递归打印节点及其子树
func (t *Tree) printNode(n *Node, level int) {
	if n == nil {
		return
	}
	// 打印当前节点信息
	fmt.Printf("第%d层: ", level)
	if n.leaf {
		fmt.Printf("叶子节点 键:%v 值:%v\n", n.keys, n.vals)
	} else {
		fmt.Printf("内部节点 键:%v\n", n.keys)
	}
	// 递归打印子节点
	for _, child := range n.kids {
		t.printNode(child, level+1)
	}
}
