package log

import (
	"fmt"
	"strings"
	"sync"
)

//logger的配置类
type LoggerConf struct {
	Name      string
	Levels    map[int]bool
	Appenders []Appender
}

func (self *LoggerConf) SetAppender(appenders ...Appender) {
	self.Appenders = appenders
}

func (self *LoggerConf) SetLevel(level int) {
	for _, l := range LogLevelMap {
		if l <= level {
			self.Levels[l] = true
		} else {
			self.Levels[l] = false
		}
	}
}

func (self *LoggerConf) SetOnlyLevels(levels ...int) {
	if self.Levels == nil {
		self.Levels = make(map[int]bool)
	}
	for _, l := range levels {
		self.Levels[l] = true
	}
}

//在self 的配置为空时 ，复制from配置
func (self *LoggerConf) copynx(from *LoggerConf) {
	if from == nil {
		return
	}
	if len(self.Appenders) == 0 && len(from.Appenders) != 0 {
		self.Appenders = from.Appenders
	}
	if len(self.Levels) == 0 && len(from.Levels) != 0 {
		if self.Levels == nil {
			self.Levels = make(map[int]bool)
		}
		for l, v := range from.Levels {
			self.Levels[l] = v
		}
	}
}

//判断配置是否完整
func (self *LoggerConf) complete() bool {
	return len(self.Appenders) != 0 && len(self.Levels) != 0
}

//日志配置树，用于可以获取可继承的日志配置
//根节点为""
//日志配置在树中的节点与名称有关，如a/b/c则是 ""->a->b->c 的树
type Tree struct {
	Root *node
}

//日志配置数节点
type node struct {
	name     string
	parent   *node
	children map[string]*node
	mutex    sync.RWMutex

	current *LoggerConf
	final   *LoggerConf
}

func NewTree(root *LoggerConf) *Tree {
	if !root.complete() {
		panic("日志Root配置不完整")
	}
	return &Tree{
		Root: newNode("", nil, root),
	}
}
func printNode(tree *node, tabs int) {
	tstr := ""
	for i := 0; i < tabs; i++ {
		tstr += "\t"
	}
	fmt.Println(tstr + tree.name)
	fmt.Println("----------")
	tree.mutex.RLock()
	for _, c := range tree.children {
		printNode(c, tabs+1)
	}
	tree.mutex.RUnlock()
	fmt.Println("----------")
}

// 拷贝一颗树，其中 current,final 配置使用原来的指针对象，从而可以更新Logger使用的配置
func (self *Tree) clone() *Tree {
	newRoot := &node{}
	self.Root.clone(newRoot)
	return &Tree{
		Root: newRoot,
	}
}

//在树中插入一个配置
func (self *Tree) insert(logger *LoggerConf) {
	self.Root.addChild(logger.Name, logger)
}
func (self *Tree) updateConf(logger *LoggerConf) {
	self.Root.updateConf(logger.Name, logger)
}

//通过名称获取一个配置
func (self *Tree) get(name string) *LoggerConf {
	if name == "" {
		return self.Root.current
	}
	child := self.Root.child(name)
	if child != nil {
		return child.current
	}
	return nil
}

//获取name的配置，当name的配置为空时，会继承name上级最接近的非空配置
func (self *Tree) inheritConf(name string) *LoggerConf {
	return self.Root.generate(name).inheritConf()
}

func newNode(name string, parent *node, current *LoggerConf) *node {
	return &node{
		name:     name,
		parent:   parent,
		current:  current,
		children: make(map[string]*node),
	}
}

func (self *node) clone(nNode *node) {
	nNode.name = self.name
	nNode.current = self.current
	nNode.final = self.final
	nNode.children = map[string]*node{}
	self.mutex.RLock()
	for key, child := range self.children {
		sNode := &node{parent: nNode}
		child.clone(sNode)
		nNode.children[key] = sNode
	}
	self.mutex.RUnlock()
}
func (self *node) updateConf(name string, logger *LoggerConf) {
	if name == "" {
		if self.isRoot() {
			self.current = logger
			self.resetFinalConf()
		}
		return
	}
	arr := strings.Split(name, "/")
	var son *node
	self.mutex.Lock()
	if n, ok := self.children[arr[0]]; ok {
		son = n
	} else {
		son = newNode(arr[0], self, nil)
		self.children[arr[0]] = son
	}
	self.mutex.Unlock()

	if len(arr) == 1 {
		son.current = logger
		son.resetFinalConf()
	} else if len(arr) > 1 {
		son.updateConf(strings.Join(arr[1:], "/"), logger)
	}
}

//添加节点的子节点
func (self *node) addChild(name string, logger *LoggerConf) {
	if name == "" {
		return
	}
	arr := strings.Split(name, "/")
	var son *node
	self.mutex.Lock()
	if n, ok := self.children[arr[0]]; ok {
		son = n
	} else {
		son = newNode(arr[0], self, nil)
		self.children[arr[0]] = son
	}
	self.mutex.Unlock()

	if len(arr) == 1 {
		son.current = logger
	} else if len(arr) > 1 {
		son.addChild(strings.Join(arr[1:], "/"), logger)
	}
}

//通过name获取节点的子节点
func (self *node) child(name string) (ret *node) {
	if name == "" {
		return
	}
	arr := strings.Split(name, "/")
	self.mutex.RLock()
	if son, ok := self.children[arr[0]]; ok {
		if len(arr) == 1 {
			ret = son
		} else {
			ret = son.child(strings.Join(arr[1:], "/"))
		}
	}
	self.mutex.RUnlock()
	return
}

// 生成子节点如果不存在并返回
func (self *node) generate(name string) (ret *node) {
	if name == "" {
		if self.name == "" {
			return self
		}
		return
	}
	arr := strings.Split(name, "/")
	self.mutex.Lock()
	son, ok := self.children[arr[0]]
	if !ok {
		son = newNode(arr[0], self, nil)
		self.children[arr[0]] = son
	}
	self.mutex.Unlock()
	if len(arr) == 1 {
		ret = son
	} else {
		ret = son.generate(strings.Join(arr[1:], "/"))
	}
	return ret
}

//判断节点是不是根节点
func (self *node) isRoot() bool {
	return self.name == ""
}

//获取节点的上级中最接近的配置非空的节点
func (self *node) higher() *node {
	if self.parent != nil {
		if self.parent.current != nil {
			return self.parent
		}
		return self.parent.higher()
	}
	return nil
}

//获取当前节点的配置，当配置为空时，会继承name上级最接近的非空配置
func (self *node) inheritConf() *LoggerConf {
	if self.final == nil {
		var cfg = &LoggerConf{}
		var cur = self
		cfg.copynx(cur.current)
		if !cfg.complete() && !cur.isRoot() {
			higher := self.parent.inheritConf()
			cfg.copynx(higher)
		}
		self.final = cfg
	}
	return self.final
}

func (self *node) resetFinalConf() {
	var cfg = &LoggerConf{}
	var cur = self
	cfg.copynx(cur.current)
	if !cfg.complete() && !cur.isRoot() {
		higher := self.parent.inheritConf()
		cfg.copynx(higher)
	}
	if self.final == nil {
		self.final = cfg
	} else {
		self.final.Name = cfg.Name
		self.final.Levels = cfg.Levels
		self.final.Appenders = cfg.Appenders
	}
	self.mutex.Lock()
	for _, c := range self.children {
		c.resetFinalConf()
	}
	self.mutex.Unlock()
}
