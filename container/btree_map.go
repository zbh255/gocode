package container

import (
	"cmp"
	"math"
	"slices"
	"sync"
	"sync/atomic"
)

type stackElement[K cmp.Ordered, V any] struct {
	node *btNode[K, V]
	tag  uint64
}

type stack[K cmp.Ordered, V any] struct {
	list []stackElement[K, V]
}

func (s *stack[K, V]) push(e stackElement[K, V]) {
	s.list = append(s.list, e)
}

func (s *stack[K, V]) pop() stackElement[K, V] {
	if len(s.list) == 0 {
		return stackElement[K, V]{
			node: nil,
		}
	}
	v := s.list[len(s.list)-1]
	s.list = s.list[:len(s.list)-1]
	return v
}

type keyword[K cmp.Ordered, V any] struct {
	key K
	val V
}

type btNode[K cmp.Ordered, V any] struct {
	keywords []keyword[K, V]
	subNodes []*btNode[K, V]
	flags    uint64
	// padding
	_ [128 - 24 - 24 - 8]byte
}

func (n *btNode[K, V]) isLeaf() bool {
	return len(n.subNodes) == 0
}

func (n *btNode[K, V]) subNodeSize() int {
	return len(n.subNodes)
}

func (n *btNode[K, V]) keywordSize() int {
	return len(n.keywords)
}

func (n *btNode[K, V]) delLastKeyword() keyword[K, V] {
	v := n.keywords[len(n.keywords)-1]
	n.keywords = n.keywords[:len(n.keywords)-1]
	return v
}

func (n *btNode[K, V]) delFirstKeyword() keyword[K, V] {
	v := n.keywords[0]
	n.keywords = n.keywords[1:]
	return v
}

type BTreeMap[K cmp.Ordered, V any] struct {
	root *btNode[K, V]
	m    uint32
	_    [128 - 4]byte
	size atomic.Uint64
	rw   sync.RWMutex
}

func NewBtreeMap[K cmp.Ordered, V any](m uint32) *BTreeMap[K, V] {
	return &BTreeMap[K, V]{
		m: m,
	}
}

func (bt *BTreeMap[K, V]) Store(key K, val V) {
	bt.put(key, val)
}

func (bt *BTreeMap[K, V]) StoreOk(key K, val V) bool {
	return !bt.put(key, val)
}

func (bt *BTreeMap[K, V]) LoadOk(key K) (V, bool) {
	return bt.get(key)
}

func (bt *BTreeMap[K, V]) Delete(key K) {
	bt.del(key)
}

func (bt *BTreeMap[K, V]) DeleteOk(key K) (V, bool) {
	return bt.del(key)
}

func (bt *BTreeMap[K, V]) put(key K, val V) bool {
	bt.rw.Lock()
	defer bt.rw.Unlock()
	if bt.root == nil {
		bt.root = &btNode[K, V]{}
		bt.root.keywords = append(bt.root.keywords, keyword[K, V]{key: key, val: val})
		bt.size.Add(1)
		return false
	}
	isReplace, isFull := bt.doPut(bt.root, key, val)
	if isFull {
		mediumElem, left, right := bt.splitNode(bt.root)
		bt.root.keywords = make([]keyword[K, V], 0, bt.m)
		bt.root.keywords = append(bt.root.keywords, mediumElem)
		bt.root.subNodes = make([]*btNode[K, V], 0, bt.m+1)
		bt.root.subNodes = append(bt.root.subNodes, left, right)
	}
	if !isReplace {
		bt.size.Add(1)
	}
	return isReplace
}

func (bt *BTreeMap[K, V]) doPut(root *btNode[K, V], key K, val V) (bool, bool) {
	index, found := slices.BinarySearchFunc(root.keywords, key, func(a keyword[K, V], b K) int {
		return cmp.Compare(a.key, b)
	})
	if found {
		root.keywords[index].val = val
		return true, false
	}
	if root.isLeaf() {
		root.keywords = slices.Insert(root.keywords, index, keyword[K, V]{key: key, val: val})
		return false, bt.nodeEQMax(root)
	} else {
		subNode := root.subNodes[index]
		isReplace, isFull := bt.doPut(subNode, key, val)
		// do split
		if isFull {
			mediumElem, left, right := bt.splitNode(subNode)
			root.keywords = slices.Insert(root.keywords, index, mediumElem)
			root.subNodes[index] = left
			root.subNodes = slices.Insert(root.subNodes, index+1, right)
		}
		return isReplace, bt.nodeEQMax(root)
	}
}

func (bt *BTreeMap[K, V]) splitNode(root *btNode[K, V]) (medium keyword[K, V], s1, s2 *btNode[K, V]) {
	medium = root.keywords[len(root.keywords)/2]
	s1 = &btNode[K, V]{
		keywords: root.keywords[:len(root.keywords)/2],
	}
	s2 = &btNode[K, V]{
		keywords: root.keywords[len(root.keywords)/2+1:],
	}
	if len(root.subNodes) > 0 {
		s1.subNodes = root.subNodes[:len(root.subNodes)/2]
		s2.subNodes = root.subNodes[len(root.subNodes)/2:]
	}
	return
}

func (bt *BTreeMap[K, V]) get(key K) (val V, found bool) {
	bt.rw.RLock()
	defer bt.rw.RUnlock()
	if bt.root == nil {
		return
	}
	node, idx := bt.findNode(bt.root, nil, key)
	if idx == -1 {
		return
	}
	found = true
	val = node.keywords[idx].val
	return
}

func (bt *BTreeMap[K, V]) findNode(root *btNode[K, V], s *stack[K, V], key K) (*btNode[K, V], int) {
	index, found := slices.BinarySearchFunc(root.keywords, key, func(e keyword[K, V], k K) int {
		return cmp.Compare(e.key, key)
	})
	if s != nil {
		s.push(stackElement[K, V]{
			node: root,
			tag:  uint64(index),
		})
	}
	if found {
		return root, index
	} else {
		if root.isLeaf() {
			return nil, -1
		} else {
			return bt.findNode(root.subNodes[index], s, key)
		}
	}
}

func (bt *BTreeMap[K, V]) del(key K) (val V, found bool) {
	bt.rw.Lock()
	defer bt.rw.Unlock()
	if bt.root == nil {
		return
	}
	s := new(stack[K, V])
	node, idx := bt.findNode(bt.root, s, key)
	if idx == -1 {
		return
	}
	val = node.keywords[idx].val
	found = true
	bt.size.Add(math.MaxUint64)
	// 在叶节点发生删除, 直接删除关键字后检查是否需要做下溢或者连接
	if node.isLeaf() {
		node.keywords = slices.Delete(node.keywords, idx, idx+1)
		s.pop()
		bt.del2(node, s)
		return
	}
	node2 := node.subNodes[idx]
	for !node2.isLeaf() {
		s.push(stackElement[K, V]{
			node: node2,
			tag:  uint64(len(node.subNodes) - 1),
		})
		node2 = node2.subNodes[len(node.subNodes)-1]
	}
	node.keywords[idx] = node2.keywords[len(node2.keywords)-1]
	node2.keywords = node2.keywords[:len(node2.keywords)-1]
	bt.del2(node2, s)
	return
}

// 处理下溢和连接
func (bt *BTreeMap[K, V]) del2(leafNode *btNode[K, V], s *stack[K, V]) {
	if !(uint32(len(leafNode.keywords)) < bt.m/2) {
		return
	}
	node := leafNode
	for {
		parent := s.pop()
		if parent.node == nil {
			break
		}
		parentNode, parentIdx := parent.node, parent.tag
		// 查看兄弟节点是否有多余的关键字, 可以借一个过来, 可能是左兄弟或者右兄弟
		if parentIdx+1 >= uint64(parentNode.subNodeSize()) && bt.nodeGEQMin(parentNode.subNodes[parentIdx-1]) {
			// 在最右侧的节点, 只能借左兄弟了
			leftNode := parentNode.subNodes[parentIdx-1]
			node.keywords = slices.Insert(node.keywords, 0, parentNode.keywords[parentIdx])
			parentNode.keywords[parentIdx] = leftNode.delLastKeyword()
		} else if parentIdx+1 < uint64(parentNode.subNodeSize()) && bt.nodeGEQMin(parentNode.subNodes[parentIdx+1]) {
			// 可以借右兄弟的节点
			rightNode := parentNode.subNodes[parentIdx+1]
			node.keywords = append(node.keywords, parentNode.keywords[parentIdx])
			parentNode.keywords[parentIdx] = rightNode.delFirstKeyword()
			break
		} else {
			// parentNode.keywords[parentIdx] = parentNode.subNodes[parentIdx+1].keywords[len(parentNode.subNodes[parentIdx+1].keywords)-1]
			// 中间节点下推
			node.keywords = append(node.keywords, parentNode.keywords[parentIdx])
			// 合并右兄弟节点
			node.keywords = append(node.keywords, parentNode.subNodes[parentIdx+1].keywords...)
			node.subNodes = append(node.subNodes, parentNode.subNodes[parentIdx+1].subNodes...)
			// 合并完成之后删除父节点中的元素和多余的子节点
			parentNode.keywords = slices.Delete(parentNode.keywords, int(parentIdx), int(parentIdx)+1)
			parentNode.subNodes = slices.Delete(parentNode.subNodes, int(parentIdx+1), int(parentIdx+1+1))
			if bt.nodeGEQMin(parentNode) {
				break
			}
			// 合并操作导致根节点没有元素了, 那就将子节点作为根节点
			if len(parentNode.keywords) == 0 && bt.root == parentNode {
				bt.root = bt.root.subNodes[0]
				break
			}
		}
		node = parentNode
	}
}

func (bt *BTreeMap[K, V]) Range(start K, fn func(key K, val V) bool) {
	bt.rw.RLock()
	defer bt.rw.RUnlock()
	if bt.root == nil {
		return
	}
	s := new(stack[K, V])
	bt.doRange(bt.root, s, start, fn)
}

func (bt *BTreeMap[K, V]) doRange(root *btNode[K, V], s *stack[K, V], start K, fn func(key K, val V) bool) {
	index, found := slices.BinarySearchFunc(root.keywords, start, func(e keyword[K, V], k K) int {
		return cmp.Compare(e.key, start)
	})
	if !found {
		if root.isLeaf() {
			return
		}
		s.push(stackElement[K, V]{node: root, tag: uint64(index)})
		bt.doRange(root.subNodes[index], s, start, fn)
		return
	} else {
		if !bt.rangeOfCentral(root, index, true, fn) {
			return
		}
		for {
			parent := s.pop()
			if parent.node == nil {
				break
			}
			if !bt.rangeOfCentral(parent.node, int(parent.tag), true, fn) {
				break
			}
		}
	}
}

func (bt *BTreeMap[K, V]) rangeOfCentral(node *btNode[K, V], eIndex int, isFirst bool, fn func(key K, val V) bool) bool {
	if isFirst {
		for i := eIndex; i < len(node.keywords); i++ {
			e := node.keywords[i]
			if !fn(e.key, e.val) {
				return false
			}
			if node.isLeaf() {
				continue
			}
			subNode := node.subNodes[i+1]
			if !bt.rangeOfCentral(subNode, 0, false, fn) {
				return false
			}
		}
	} else {
		for i := eIndex; i < len(node.keywords); i++ {
			e := node.keywords[i]
			if node.isLeaf() {
				if !fn(e.key, e.val) {
					return false
				}
			} else {
				if !bt.rangeOfCentral(node.subNodes[i], 0, false, fn) {
					return false
				}
				if !fn(e.key, e.val) {
					return false
				}
				if !bt.rangeOfCentral(node.subNodes[i+1], 0, false, fn) {
					return false
				}
			}
		}
	}
	return true
}

func (bt *BTreeMap[K, V]) Len() int {
	return int(bt.size.Load())
}

func (bt *BTreeMap[K, V]) High() int {
	return 0
}

func (bt *BTreeMap[K, V]) MaxKey() (key K) {
	bt.rw.RLock()
	defer bt.rw.RUnlock()
	if bt.root == nil || len(bt.root.keywords) == 0 {
		return
	}
	return bt.maxKey(bt.root)
}

func (bt *BTreeMap[K, V]) maxKey(root *btNode[K, V]) (key K) {
	for {
		key = root.keywords[len(root.keywords)-1].key
		if root.isLeaf() {
			break
		}
		root = root.subNodes[len(root.subNodes)-1]
	}
	return
}

func (bt *BTreeMap[K, V]) nodeGEQMin(node *btNode[K, V]) bool {
	return uint32(len(node.keywords)) >= bt.m/2-1
}

func (bt *BTreeMap[K, V]) nodeEQMax(node *btNode[K, V]) bool {
	return uint32(len(node.keywords)) == bt.m-1
}
