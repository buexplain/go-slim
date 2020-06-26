package tree

import (
	"errors"
	"regexp"
	"strings"
)

//节点
type Node struct {
	//父级
	parent *Node
	//子节点
	child map[rune]*Node
	//节点名称
	name rune
	//动态参数名称
	paramName string
	//节点中存储的数据
	//不能是nil
	data interface{}
}

//新建一个节点
func NewNode(parent *Node, name rune) *Node {
	tmp := &Node{
		parent:    parent,
		child:     make(map[rune]*Node),
		name:      name,
		paramName: "",
		data:      nil,
	}
	return tmp
}

//节点树
type Tree struct {
	//节点
	root *Node
	//是否严格区分大小写
	strict bool
}

//新建一颗节点树
func New(strict bool) *Tree {
	return &Tree{
		root:   NewNode(nil, 0),
		strict: strict,
	}
}

//给当前节点树添加一条路径
func (this *Tree) Add(path string, data interface{}) error {
	//检查路径
	if len(path) == 0 {
		return errors.New("path not allow empty")
	}
	//检查存储值
	if data == nil {
		return errors.New("data not allow nil")
	}
	var currPath = path
	//提取参数名
	re := regexp.MustCompile(`:[0-9A-Za-z_\-\.]+`)
	param := re.FindAllString(path, -1)
	if len(param) > 0 {
		path = re.ReplaceAllString(path, ":")
	}
	//判断是否严格大小写
	if !this.strict {
		path = strings.ToLower(path)
	}
	//开始插入节点
	var i int = 0
	var currNode *Node = this.root
	for _, v := range path {
		if node, ok := currNode.child[v]; ok {
			//节点存在继续遍历
			currNode = node
			//动态参数索引自增
			if v == ':' {
				i++
			}
		} else {
			//检查节点是否冲突
			//冲突规则：静态节点上插入动态节点或者是动态节点上插入静态节点
			if _, ok := currNode.child[':']; (v == ':' && len(currNode.child) > 0) || (v != ':' && ok) {
				for k, _ := range currNode.child {
					if v == ':' {
						if k != ':' {
							currNode = currNode.child[k]
							break
						}
					} else {
						if k == ':' {
							currNode = currNode.child[k]
							break
						}
					}
				}
				conflict := ""
				for currNode.parent != nil {
					if currNode.paramName == "" {
						conflict = string(currNode.name) + conflict
					} else {
						conflict = ":" + currNode.paramName + conflict
					}
					currNode = currNode.parent
				}
				return errors.New("path conflict " + conflict + " " + currPath)
			}
			//新增子节点
			currNode.child[v] = NewNode(currNode, v)
			//新增的子节点为动态节点，提取参数名称
			if v == ':' {
				//检查参数名称是否正确
				if i >= len(param) || param[i][1:] == "" {
					return errors.New("path format error")
				}
				currNode.child[v].paramName = param[i][1:]
				i++
			}
			//继续向下遍历
			currNode = currNode.child[v]
		}
	}
	//存储data
	currNode.data = data
	//返回成功
	return nil
}

//搜索到的参数的存储接口
type Store interface {
	Set(key string, v interface{})
}

//查找当前路径是否在节点树中
func (this *Tree) Search(path string, param Store) (interface{}, bool) {
	if !this.strict {
		path = strings.ToLower(path)
	}
	var currNode *Node = this.root
loop:
	for k, v := range path {
		if node, ok := currNode.child[v]; ok {
			//静态节点，继续遍历
			currNode = node
		} else {
			//判断是否存在动态节点
			if node, ok := currNode.child[':']; ok {
				//找到了动态的节点
				currNode = node
				//提取参数到当前的请求
				index := strings.IndexByte(path[k:], '/')
				if index == -1 {
					param.Set(node.paramName, path[k:])
					break
				} else {
					param.Set(node.paramName, path[k:k+index])
					path = path[k+index:]
					goto loop
				}
			} else {
				//没找到动态节点，直接返回
				return nil, false
			}
		}
	}
	//当前节点并非挂载数据的节点
	if currNode.data == nil {
		return nil, false
	}
	return currNode.data, true
}
