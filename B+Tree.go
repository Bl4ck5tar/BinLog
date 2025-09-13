package main

import (
	"sort"
	"fmt"
)
const (
	order = 3
)
//B+ Tree节点
type Node struct {
	keys		[]int		//键数组
	children	[]*Node		//子节点数组
	isLeaf		bool		//是否为叶子节点
	next		*Node		//叶子节点指针，连接到下一个叶子节点
}
//B+ Tree结构体
type BPlusTree struct {
	root *Node
}
//创建节点
func newNode(isLeaf bool) *Node {
	return &Node{
		keys:		make([]int, 0),
		children:	make([]*Node, 0),
		isLeaf: 	isLeaf,
		next: 		nil,
	}
}
//创建一颗新的B+ Tree
func newBPlusTree() *BPlusTree {
	root := newNode(true)
	return &BPlusTree{
		root: root,
	}
}
//插入一个键值到 B+ Tree
func (tree *BPlusTree) insert(key int) {
	root := tree.root
	if len(root.keys) == order-1 { //根节点满了，需要分裂
		newRoot := newNode(false)
		newRoot.children = append(newRoot.children, tree.root)
		tree.splitChild(newRoot, 0)
		tree.root = newRoot
	}
	tree.insertNonFull(tree.root, key)
}
//插入到非满节点
func (tree *BPlusTree) insertNonFull(node *Node, key int) {
	if node.isLeaf {
		//叶子节点，插入键值
		node.keys = append(node.keys, key)
		sort.Ints(node.keys)
	}else {
		//非叶子节点，寻找合适的子节点插入
		i := len(node.keys) - 1
		for i>=0 && key<node.keys[i] {
			i--
		}
		i++
		if len(node.children[i].keys) == order-1 {
			tree.splitChild(node, i)
			if key > node.keys[i] {
				i++
			}
		}
		tree.insertNonFull(node.children[i], key)
	}
}
//分裂子节点
func (tree *BPlusTree) splitChild(parent *Node, index int) {
	child := parent.children[index]
	midIndex := order / 2

	//创建一个新的节点
	newNode := newNode(child.isLeaf)
	parent.keys = append(parent.keys, 0)
	copy(parent.keys[index+1:], parent.keys[index:])
	parent.keys[index] = child.keys[midIndex]

	//将右半部分键和子节点移到新节点
	newNode.keys = append(newNode.keys, child.keys[midIndex+1:]...)
	child.keys = child.keys[:midIndex]
	if !child.isLeaf {
		newNode.children = append(newNode.children, child.children[midIndex+1:]...)
		child.children = child.children[:midIndex+1]
	}else {
		newNode.next = child.next
		child.next = newNode
	}

	//将新节点添加到父节点
	parent.children = append(parent.children, nil)
	copy(parent.children[index+2:], parent.children[index+1:])
	parent.children[index+1] = newNode

}
//查找键值是否存在
func (tree *BPlusTree) search(key int) bool {
	return tree.searchNode(tree.root, key)
}
//递归查找子节点
func (tree *BPlusTree) searchNode(node *Node, key int) bool {
	if node.isLeaf {
		//在叶子节点中查找
		for _,k := range node.keys {
			if k == key {
				return true
			}
		}
		return false
	}

	//在非叶子节点中查找
	i := len(node.keys) - 1
	for i>=0 && key<node.keys[i] {
		i--
	}
	i++
	return tree.searchNode(node.children[i], key)
}
//打印树结构
func (tree *BPlusTree) print() {
	tree.printNode(tree.root, 0)
}

func (tree *BPlusTree) printNode(node *Node, level int) {
	for range level {
		fmt.Print("\t")
	}
	fmt.Println(node.keys)
	if !node.isLeaf {
		for _, child := range node.children {
			tree.printNode(child, level+1)
		}
	}
}

func main() {
	tree := newBPlusTree()

	keys := []int{10,20,5,6,45,30,15,25,45,40,35}
	for _, key := range keys {
		tree.insert(key)
	}

	tree.print()

	fmt.Println("查找15", tree.search(15))
	fmt.Println("查找70", tree.search(70))
}