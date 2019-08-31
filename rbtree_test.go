package rbtree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// todo : test left and right rotate
type testDataTyp uint8

const (
	customTestData testDataTyp = iota
	randomTestData
)

func randNumber(key int) uint32 {
	return uint32(rand.Intn(key))
}

func initTestData(typ testDataTyp, key, nodeCount int) ([]*Node, []uint32) {
	nodes := make([]*Node, 0)
	values := make([]uint32, 0)
	switch typ {
	case randomTestData:
		for i := 0; i < nodeCount; i++ {
			value := randNumber(key + i)
			values = append(values, value)
			nodes = append(nodes, NewNode(value, red))
		}
	case customTestData:
		//                    黑 黑 黑 黑  红  红  红 黑  红  红 -> 调整后各个节点的颜色
		tmpValues := []uint32{1, 9, 9, 3, 1, 18, 1, 26, 4, 12}
		values = append(values, tmpValues...)
		for i := 0; i < len(values); i++ {
			nodes = append(nodes, NewNode(values[i], red))
		}
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return nodes, values
}

// 期待生成的红黑树结构如下图：
// 这是一种非常不明智的测试方法，但是套用之前了解过的极大似然函数的概念在这里的话，应该也是能说的通的
// 定制一组数据然后插入到红黑树中，并且结果符合预期，则这颗红黑树的正确的可能性要比非正确可能性更高
/******************************************************************************
						   9 (黑)
						/	    \
 			     	   /	     \
 					  1 (红）     18 (红)
 				    /   \	     /   \
                   /     \      /     \
               (黑)1  （黑）3    9 (黑) 26(黑)
                          / \    \
                         /   \    \
                        1(红) 4(红) 12(红)
******************************************************************************/
func customTestCaseCheck(t *testing.T, rbtree *RBTree) {
	check := func(root *Node, c color, left, right uint32) {
		assert.Equal(t, c, root.color)
		assert.Equal(t, left, root.left.key)
		assert.Equal(t, right, root.right.key)
	}
	check(rbtree.root, black, 1, 18)
	check(rbtree.root.left, red, 1, 3)
	check(rbtree.root.left.right, black, 1, 4)
	check(rbtree.root.right, red, 9, 26)
	check(rbtree.root.right.left, black, 0, 12)

}

func getLeafNodes(sentinel, node *Node, nodes *[]*Node) {
	if node == sentinel {
		return
	}
	if node.left == sentinel && node.right == sentinel {
		*nodes = append(*nodes, node)
		return
	}
	if node.right != sentinel {
		getLeafNodes(sentinel, node.right, nodes)
	}
	if node.left != sentinel {
		getLeafNodes(sentinel, node.left, nodes)
	}
}

// rbtreePropritiesCheck used to check rbtree's proprities
func rbtreePropritiesCheck(t *testing.T, rbt *RBTree) {
	// 空的红黑树，直接退出
	if rbt.root == rbt.sentinel {
		return
	}
	// 1: 根节点一定是黑色
	assert.Equal(t, black, rbt.root.color)
	// 查找所有叶节点
	leafNodes := make([]*Node, 0)
	getLeafNodes(rbt.sentinel, rbt.root, &leafNodes)
	blackNum := -1
	for _, leafN := range leafNodes {
		// 所有页节点的左右孩子均指向哨兵节点
		assert.Equal(t, rbt.sentinel, leafN.left)
		assert.Equal(t, rbt.sentinel, leafN.right)
		tmpBlackNum := 0
		for leafN != rbt.root {
			if leafN.judgeColor(black) {
				tmpBlackNum++
			} else {
				// 4: 一个红色节点的父节点不可能是红色
				assert.Equal(t, false, leafN.parent.judgeColor(red))
			}
			leafN = leafN.parent
		}
		if blackNum == -1 {
			// 5：初始化从根节点到叶节点的黑色节点的个数
			blackNum = tmpBlackNum
		} else {
			// 判断所有到所有叶结点的黑色节点的个数是相等的
			assert.Equal(t, blackNum, tmpBlackNum)
		}

	}

}

func orderCheck(t *testing.T, expectedOrder []uint32, rbTree *RBTree) {
	actualOrder := make([]uint32, 0)
	PreOrder(rbTree.root, rbTree.sentinel, &actualOrder)
	assert.Equal(t, expectedOrder, actualOrder)
}

func insert(typ testDataTyp, key, nodeCount int) (*RBTree, []uint32, []*Node) {
	sentinel := NewNode(0, black)
	rbTree := Init(sentinel, InsertValue)
	nodes, expectedOrder := initTestData(typ, key, nodeCount)
	for _, node := range nodes {
		rbTree.Insert(node)
	}
	return rbTree, expectedOrder, nodes
}

func initInsertTestCase(t *testing.T, typ testDataTyp, desc string, key, nodeCount int) {
	fmt.Printf("[TestCase] desc %s, key %d, node count %d\n", desc, key, nodeCount)
	rbTree, expectedOrder, _ := insert(typ, key, nodeCount)
	orderCheck(t, expectedOrder, rbTree)
	rbtreePropritiesCheck(t, rbTree)
	if typ == customTestData {
		customTestCaseCheck(t, rbTree)
	}
}

func remove(arr *[]uint32, e uint32) {
	l := len(*arr)
	for i := 0; i < l; i++ {
		if (*arr)[i] == e {
			*arr = append((*arr)[0:i], (*arr)[i+1:]...)
			break
		}
	}
}

func initDeleteTestCase(t *testing.T, dc deleteCases) {
	fmt.Printf("[TestCase] desc %s, key %d, node count %d\n", dc.desc, dc.ics.key, dc.ics.nodeCount)
	rbTree, expectedOrder, nodes := insert(dc.ics.typ, dc.ics.key, dc.ics.nodeCount)
	for i := dc.deleteIndex; i < dc.ics.nodeCount; i = i * dc.ics.factor {
		node := nodes[i]
		remove(&expectedOrder, node.key)
		rbTree.Delete(node)
		orderCheck(t, expectedOrder, rbTree)
		rbtreePropritiesCheck(t, rbTree)
	}
}

type insertCases struct {
	typ       testDataTyp
	desc      string
	key       int
	factor    int // factor used to adjust insert node count
	nodeCount int
	repeTimes int
}

type deleteCases struct {
	ics         *insertCases
	desc        string
	deleteIndex int
	repeTimes   int
}

func TestRBTreeDelete(t *testing.T) {
	cs := []deleteCases{
		{
			desc: "test rand insert and delete data in red black tree",
			ics: &insertCases{
				typ:       randomTestData,
				key:       100,
				factor:    4,
				nodeCount: 11,
			},
			repeTimes:   1,
			deleteIndex: 2,
		},
	}
	for _, c := range cs {
		for i := 0; i < c.repeTimes; i++ {
			initDeleteTestCase(t, c)
			c.ics.key *= c.ics.factor
			c.ics.nodeCount *= c.ics.factor
		}
	}
}

func TestRBTreeInsert(t *testing.T) {
	cs := []insertCases{
		{
			desc:      "test rand data insert to red black tree",
			typ:       randomTestData,
			key:       10,
			factor:    10,
			nodeCount: 1,
			repeTimes: 3,
		}, {

			desc:      "test custom data insert to red black tree",
			typ:       customTestData,
			repeTimes: 1,
			nodeCount: 10,
		},
	}
	for _, c := range cs {
		for i := 0; i < c.repeTimes; i++ {
			initInsertTestCase(t, c.typ, c.desc, c.key, c.nodeCount)
			c.key *= c.factor
			c.nodeCount *= c.factor
		}
	}
}
