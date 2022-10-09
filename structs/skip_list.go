/*
A skip list is a data structure that allows fast search within an ordered
sequence of elements, O(log n) complexity.
See https://en.wikipedia.org/wiki/Skip_list for more details.

How to sort the elements in the skip list?
A skip list Node contains a string value and a float64 score.
When comparing two nodes A and B, we consider both the score and the value.
That means, when A < B, it has to suffice that:
1. if A.score < B.score -> A < B
2. else if A.score == B.score && A.val < B.val -> A < B
*/

package structs

import "math/rand"

const sklMaxLevel = 32 /* Should be enough for 2^32 elements */
const sklLevelP = 0.25

type SkipListNode struct {
	val   string
	score float64
	back  *SkipListNode
	level []skipListLevel
}

type skipListLevel struct {
	forward *SkipListNode
	span    int
}

type SkipList struct {
	header, tail *SkipListNode // header is a dummy node, tail is the last node
	length       int           // number of the skip list nodes
	level        int           // current max level of the skip list
}

func NewSklNode(level int, score float64, val string) *SkipListNode {
	return &SkipListNode{
		val:   val,
		score: score,
		level: make([]skipListLevel, level),
	}
}

func NewSkipList() *SkipList {
	skl := &SkipList{
		header: NewSklNode(sklMaxLevel, 0, ""),
		level:  1,
		length: 0,
		tail:   nil,
	}
	for i := 0; i < sklMaxLevel; i++ {
		skl.header.level[i].forward = nil
		skl.header.level[i].span = 0
	}
	skl.header.back = nil
	return skl
}

// sklRandomLevel returns a random level for the new skip list node
func sklRandomLevel() int {
	level := 1
	for float64(rand.Int31()&0xFFFF) < (sklLevelP * 0xFFFF) {
		level += 1
	}
	if level > sklMaxLevel {
		return sklMaxLevel
	}
	return level
}

// Insert a new element in the skip list and return the new node
func (skl *SkipList) Insert(score float64, val string) *SkipListNode {
	update := make([]*SkipListNode, sklMaxLevel)
	rank := make([]int, sklMaxLevel)
	var i, level int

	// find the insert position and store the update path and rank
	x := skl.header
	for i = skl.level - 1; i >= 0; i-- {
		if i == skl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil && (x.level[i].forward.score < score ||
			(x.level[i].forward.score == score && x.level[i].forward.val < val)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	// get a random level for the new node
	// initial unused levels of header if the new lever is larger than the current level
	level = sklRandomLevel()
	if level > skl.level {
		for i = skl.level; i < level; i++ {
			rank[i] = 0
			update[i] = skl.header
			update[i].level[i].span = skl.length
		}
		skl.level = level
	}

	// create the new node
	x = NewSklNode(level, score, val)

	// update the forward and span of the nodes in updated path
	for i = 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	//	increment the span for the nodes which go across the new node
	for i = level; i < skl.level; i++ {
		update[i].level[i].span++
	}

	// update the back pointer
	if update[0] == skl.header {
		x.back = nil
	} else {
		x.back = update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.back = x
	} else {
		skl.tail = x
	}

	skl.length++

	return x
}

// Delete the node when the score and val are both equal to the given one
func (skl *SkipList) Delete(score float64, val string) int {
	update := make([]*SkipListNode, sklMaxLevel)
	var i int

	// find the deleting position and store the update path
	x := skl.header
	for i = skl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && (x.level[i].forward.score < score ||
			(x.level[i].forward.score == score && x.level[i].forward.val < val)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward

	// delete the node if found
	if x != nil && score == x.score && val == x.val {
		skl.deleteNode(x, update)
		return 1
	}
	return 0
}

func (skl *SkipList) deleteNode(x *SkipListNode, update []*SkipListNode) {
	var i int

	// update the forward and span of the nodes in updated path
	for i = 0; i < skl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}

	// update the back pointer
	if x.level[0].forward != nil {
		x.level[0].forward.back = x.back
	} else {
		skl.tail = x.back
	}

	// remove the unused levels of header
	for skl.level > 1 && skl.header.level[skl.level-1].forward == nil {
		skl.level--
	}

	skl.length--
}
