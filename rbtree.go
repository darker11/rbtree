package rbtree

/********************* rbtree 性质 ************************************
1. 每个节点是红色或者黑色
2. 根节点是黑色的
3. 每个叶节点(nil)是黑色的，实际上 nil 只有一个
4. 如果一个节点是红色的，则它的两个子节点都是黑色的(一个红色节点的父节点不可能是红色)
5. 对每个节点，从该节点到其所有后代叶结点的简单路径(不包含重复点的路径）上，均包含相同数目的黑色节点
参考：
	http://www.freecls.com/a/2712/d5
	https://blog.csdn.net/oyw5201314ck/article/details/78329607
	https://zh.wikipedia.org/wiki/%E7%BA%A2%E9%BB%91%E6%A0%91?spm=ata.13261165.0.0.351b2ed1lVP8Tq
	https://www.cs.usfca.edu/~galles/visualization/RedBlack.html?spm=ata.13261165.0.0.351b2ed1lVP8Tq
**********************************************************************/

// InsertCallBack is function type, use need to implement it.
type InsertCallBack func(root, node, sentinel *Node)

func compare(result bool, l, r **Node) **Node {
	if result {
		return l
	}
	return r
}

// implement examples:

//InsertValue used to insert value to RBTree.
func InsertValue(tmp, node, sentinel *Node) {
	var p **Node
	for {

		p = compare(node.key < tmp.key, &(tmp.left), &(tmp.right))
		if *p == sentinel {
			break
		}
		tmp = *p
	}
	// 将 node 插入 RBTree
	*p = node
	node.parent = tmp
	node.left, node.right = sentinel, sentinel
	node.setColor(red)
}

// RBTree defines red-black-tree structure.
type RBTree struct {
	size           uint32
	root           *Node
	sentinel       *Node
	insertCallBack InsertCallBack
}

// Init used to initialize RBTree structure.
func Init(sentinel *Node, icb InsertCallBack) *RBTree {
	sentinel.setColor(black)
	return &RBTree{
		root:           sentinel,
		sentinel:       sentinel,
		insertCallBack: icb,
	}
}

// leftRotate used to left rotate red-black-tree.
/********************************************************
       * (n)                      * (tmp)
	    \                       /   \
		 \                     /     \
		  *  (tmp)  -->       *  (n)  *
		 /  \                   \
        /    \                   \
       *      *                   *
*********************************************************/
func (rb *RBTree) leftRotate(n *Node) {
	tmp := n.right
	n.right = tmp.left

	if tmp.left != rb.sentinel {
		tmp.left.parent = n
	}
	tmp.parent = n.parent

	if n == rb.root {
		rb.root = tmp
	} else if n == n.parent.left {
		n.parent.left = tmp
	} else {
		n.parent.right = tmp
	}

	tmp.left = n
	n.parent = tmp
}

// rightRotate used to right rotate red-black-tree.
/********************************************************
		   * (n)              * (tmp)
		 /                   /  \
	    /                   /    \
	   * (tmp)     -->     *      * (n)
	  / \                        /
	 /   \                      /
    *     * 				   *
*********************************************************/
func (rb *RBTree) rightRotate(n *Node) {
	tmp := n.left
	n.left = tmp.right
	if tmp.right != rb.sentinel {
		tmp.right.parent = n
	}
	tmp.parent = n.parent
	if n == rb.root {
		rb.root = tmp
	} else if n.parent.right == n {
		n.parent.right = tmp
	} else {
		n.parent.left = tmp
	}
	tmp.right = n
	n.parent = tmp
}

// Insert used to inset node to RBTree structure.
func (rb *RBTree) Insert(n *Node) {
	var (
		root     = &(rb.root)
		sentinel = rb.sentinel
		temp     *Node
	)
	rb.size++
	// insert root node
	if *root == sentinel {
		n.parent = nil
		n.left, n.right = sentinel, sentinel
		n.setColor(black)
		*root = n
		return
	}
	// insert internal node
	rb.insertCallBack(*root, n, sentinel)

	// reblance rbtree
	for n != *root && (n.parent).judgeColor(red) {
		if n.parent == n.parent.parent.left {
			temp = n.parent.parent.right
			// 无需调整，只需要重新上色即可
			if temp.judgeColor(red) {
				n.parent.setColor(black)
				temp.setColor(black)
				n.parent.parent.setColor(red)
				n = n.parent.parent
			} else {
				// n 是右节点，需要先左旋转
				if n == n.parent.right {
					// n.parent 和 n 均为红色的节点
					n = n.parent
					// 左旋，并未重新上色
					rb.leftRotate(n)
				}
				n.parent.setColor(black)
				// 重置 n.parent.parent 为红色节点，可能违反性质 4
				n.parent.parent.setColor(red)
				// 右旋调整，局部平衡
				rb.rightRotate(n.parent.parent)
			}
		} else {
			// 镜像代码，跟上面的理解成反的就可以啦
			temp = n.parent.parent.left
			if temp.judgeColor(red) {
				n.parent.setColor(black)
				temp.setColor(black)
				n.parent.parent.setColor(red)
				n = n.parent.parent
			} else {
				if n == n.parent.left {
					n = n.parent
					rb.rightRotate(n)
				}
				n.parent.setColor(black)
				n.parent.parent.setColor(red)
				rb.leftRotate(n.parent.parent)
			}
		}
	}
	(*root).setColor(black)
}

// Delete used to delete the node from red-black-tree.
func (rb *RBTree) Delete(n *Node) {
	var (
		// subst 是待删除节点
		// tmp   是删除后的替代节点
		// w     是待删除节点的兄弟
		tmp, subst, w *Node
	)
	// 找到删除节点和替代节点
	if n.left == rb.sentinel {
		tmp = n.right
		subst = n
	} else if n.right == rb.sentinel {
		tmp = n.left
		subst = n
	} else {
		subst = mix(n.right, rb.sentinel)
		if subst.left != rb.sentinel {
			tmp = subst.left
		} else {
			tmp = subst.right
		}
	}

	if subst == rb.root {
		rb.root = tmp
		rb.root.setColor(black)
		n.left, n.right, n.parent, n.key = nil, nil, nil, 0
		return
	}
	isRed := subst.judgeColor(red)

	// 将删除节点的 parent 左孩子或者右孩子置为替代节点 tmp
	if subst == subst.parent.left {
		subst.parent.left = tmp
	} else {
		subst.parent.right = tmp
	}

	// 将删除节点的 parent 设置为替代节点的 parent
	if subst == n {
		tmp.parent = subst.parent
	} else {
		if subst.parent == n {
			// ？总觉得这里有问题
			tmp.parent = subst
		} else {
			tmp.parent = subst.parent
		}
		subst.left = n.left
		subst.right = n.right
		subst.parent = n.parent
		subst.setColor(n.color)

		// 将待删除节点用 subst 替换
		// 建立 subst 节点与 node parent 的连接关系
		if n == rb.root {
			rb.root = subst
		} else {
			if n == n.parent.left {
				n.parent.left = subst
			} else {
				n.parent.right = subst
			}
		}

		// 建立 subst 节点与 node child 的连接关系
		if subst.left != rb.sentinel {
			subst.left.parent = subst
		}
		if subst.right != rb.sentinel {
			subst.right.parent = subst
		}
	}
	n.left, n.right, n.parent, n.key = nil, nil, nil, 0

	// 删除红色节点不会破坏红黑树性质，直接返回即可
	if isRed {
		return
	}

	// 1. tmp 是根节点，即需删除节点 node 为根节点，将 tmp 置为黑色即可
	// 2. tmp 不是根节点，且是红色，将 tmp 置为黑色即可
	// 3. tmp 不是根节点，且是黑色，违背性质 5 需要调整
	for tmp != rb.root && tmp.judgeColor(black) {
		if tmp == tmp.parent.left {
			w = tmp.parent.right
			if w.judgeColor(red) {
				w.setColor(black)
				tmp.parent.setColor(red)
				rb.leftRotate(tmp.parent)
				w = tmp.parent.right
			}
			// if w == rb.sentinel {
			// 	break
			// }
			if w.left.judgeColor(black) && w.right.judgeColor(black) {
				w.setColor(red)
				tmp = tmp.parent
			} else {
				if w.right.judgeColor(black) {
					w.left.setColor(black)
					w.setColor(red)
					rb.rightRotate(w)
					w = tmp.parent.right
				}
				// 交换 parent 和 右孩子的节点颜色, w 成为 tmp 新的 parent
				w.setColor(tmp.parent.color)
				tmp.parent.setColor(black)
				w.right.setColor(black)
				rb.leftRotate(tmp.parent)
				tmp = rb.root
			}
		} else {
			w = tmp.parent.left
			if w.judgeColor(red) {
				w.setColor(black)
				tmp.parent.setColor(red)
				rb.rightRotate(tmp.parent)
				w = tmp.parent.left
			}
			// if w == rb.sentinel {
			// 	break
			// }
			if w.left.judgeColor(black) && w.right.judgeColor(black) {
				w.setColor(red)
				tmp = tmp.parent
			} else {
				if w.left.judgeColor(black) {
					w.right.setColor(black)
					w.setColor(red)
					rb.leftRotate(w)
					w = tmp.parent.left
				}
				w.setColor(tmp.parent.color)
				tmp.parent.setColor(black)
				w.left.setColor(black)
				rb.rightRotate(tmp.parent)
				tmp = rb.root
			}
		}
	}
	tmp.setColor(black)
}

// PreOrder used to previous order red-black-tree structure.
func PreOrder(node *Node, sentinel *Node, result *[]uint32) {
	if node != sentinel {
		PreOrder(node.left, sentinel, result)
		*result = append(*result, node.key)
		PreOrder(node.right, sentinel, result)
	}
}
