package structs

type ZSet struct {
	dict map[string]*SkipListNode
	skl  *SkipList
}
