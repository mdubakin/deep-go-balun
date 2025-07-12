package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	В домашнем задании нужно реализовать упорядоченный словарь.

	Идея упорядоченного словаря заключается в том, что он будет реализован на основе бинарного дерева поиска (BST).
	Дерево будет строиться только по ключам элементов, значения элементов при построении дерева не учитываются.
	Элементы с одинаковыми ключами в упорядоченном словаре хранить нельзя.

	API для упорядоченного словаря будет выглядеть следующим образом:
		type OrderedMap struct { ... }
		func NewOrderedMap() OrderedMap                      // создать упорядоченный словарь
		func (m *OrderedMap) Insert(key, value int)          // добавить элемент в словарь
		func (m *OrderedMap) Erase(key int)                  // удалить элемент из словари
		func (m *OrderedMap) Contains(key int) bool          // проверить существование элемента в словаре
		func (m *OrderedMap) Size() int                      // получить количество элементов в словаре
		func (m *OrderedMap) ForEach(action func(int, int))  // применить функцию к каждому элементу словаря от меньшего к большему
*/

// go test -v homework_test.go

type OrderedMap struct {
	size int
	root *node
}

type node struct {
	key         int
	value       int
	left, right *node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		size: 0,
		root: nil,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	var x, y *node = m.root, nil

	// update existing key
	for x != nil {
		if x.key == key {
			x.value = value
			return
		} else {
			y = x
			if x.key > key {
				x = x.left
			} else {
				x = x.right
			}
		}
	}

	// insert new one
	newNode := &node{key: key, value: value}
	if nil == y {
		m.root = newNode
	} else {
		if y.key > key {
			y.left = newNode
		} else {
			y.right = newNode
		}
	}
	m.size++
}

func (m *OrderedMap) Erase(key int) {
	var deeper, prev *node = m.root, nil

	// search for a key to delete
	for deeper != nil {
		if deeper.key == key {
			// key was found
			break
		} else {
			prev = deeper
			if deeper.key > key {
				deeper = deeper.left
			} else {
				deeper = deeper.right
			}
		}
	}
	// there is no such key in the map
	if nil == deeper {
		return
	}
	// if found key doesn't have a right node (bigger value)
	if nil == deeper.right {
		if nil == prev {
			// in that case key == root
			// so, we can set root = x.left, because it is either nil or a node
			//                (deeper)
			//               //      \\
			//      (deeper.left)   (nil)
			//           |
			//           V
			//      can be nil also
			m.root = deeper.left
		} else {
			if deeper == prev.left {
				//				 (prev) |	 |           (prev)
				//				//      |    |          //
				// 			(deeper)    | => |    (deeper.left)
				//         //           |	 |
				//     (deeper.left)    |	 |
				prev.left = deeper.left
			} else {
				//		(prev)          |	 |      (prev)
				//			 \\		    |    |           \\
				// 			(deeper)    | => |        (deeper.left)
				//         //           |	 |
				//    (deeper.left)     |	 |
				prev.right = deeper.left
			}
		}
	} else { // if found key HAS a right node (bigger value)
		leftMostNode := deeper.right
		prev = nil
		for leftMostNode.left != nil {
			prev = leftMostNode
			leftMostNode = leftMostNode.left
		}
		if prev != nil {
			prev.left = leftMostNode.right
		} else {
			// only if deeper.right didn't have ANY left node
			//                 (deeper)
			//        		//		    \\
			//            (nil)     (deeper.right)
			//					  //			\\
			//					(nil)		(deeper.right.right) <- this one will choose
			//
			// equivalent => deeper.right = deeper.right.right
			deeper.right = leftMostNode.right
		}
		deeper.key = leftMostNode.key
		deeper.value = leftMostNode.value
	}
	m.size--
}

func (m *OrderedMap) Contains(key int) bool {
	x := m.root
	for x != nil {
		if x.key == key {
			return true
		} else if x.key > key {
			x = x.left
		} else {
			x = x.right
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	inorder(m.root, action)
}

func inorder(n *node, action func(int, int)) {
	if n != nil {
		inorder(n.left, action)
		action(n.key, n.value)
		inorder(n.right, action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
