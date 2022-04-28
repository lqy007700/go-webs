package go_webs

import (
	"strings"
)

type node struct {
	pattern  string  // 全路由
	part     string  // 路由部分
	children []*node // 子节点
	isWild   bool    // 是否精确匹配
}

// 匹配节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part && child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 添加前缀树
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 最后节点存储全路径
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 无匹配节点，创建
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 递归查找路由
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}