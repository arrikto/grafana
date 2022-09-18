package accesscontrol

import "strings"

func BuildPermissionTrie(permissions []Permission) *Trie {
	t := NewTree()
	for _, p := range permissions {
		t.addNode(t.Root, p.Action, p.Scope)
	}
	return t
}

func NewTree() *Trie {
	return &Trie{
		Root: &Node{
			Root:     true,
			Path:     "",
			Children: map[string]Node{},
			Actions:  map[string]bool{},
		},
	}
}

type Trie struct {
	Root *Node `json:"root"`
}

func (t *Trie) HasAccess(action, scope string) bool {
	return t.Root.hasAccess(action, scope)
}

func (t *Trie) addNode(n *Node, action, path string) {
	if n.Path == path || path == "*" {
		n.Actions[action] = true
		return
	}

	// no need to child Node when parent has action
	if n.Actions[action] {
		return
	}

	prefix := path
	idx := strings.Index(prefix, ":")
	if idx > 0 {
		prefix = path[:idx]
	}

	c, ok := n.Children[prefix]
	if !ok {
		n.Children[prefix] = Node{
			Path:     prefix,
			Actions:  map[string]bool{},
			Children: map[string]Node{},
		}
		c = n.Children[prefix]
	}

	t.addNode(&c, action, path[idx+1:])
}

type Node struct {
	Root     bool            `json:"root"`
	Path     string          `json:"path"`
	Actions  map[string]bool `json:"actions"`
	Children map[string]Node `json:"children"`
}

func (n Node) hasAccess(action, path string) bool {
	if n.Path == path {
		return n.Actions[action]
	}

	if n.Actions[action] {
		return true
	}

	return false
}
