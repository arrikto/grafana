package accesscontrol

import "strings"

const (
	pathDelim  = "/"
	scopeDelim = ":"
)

func TrieFromPermissions(permissions []Permission) *Trie {
	t := newTrie()
	for _, p := range permissions {
		t.Actions[p.Action] = true
		t.Root.addNode(p.Action, p.Scope, ":")
	}
	return t
}

func TrieFromMap(permissions map[string][]string) *Trie {
	t := newTrie()
	for action, scopes := range permissions {
		t.Actions[action] = true
		for _, scope := range scopes {
			t.Root.addNode(action, scope, ":")
		}
	}
	return t
}

func newTrie() *Trie {
	return &Trie{
		Root: &Node{
			Root:     true,
			Path:     "",
			Children: map[string]Node{},
			Actions:  map[string]bool{},
		},
		Actions: map[string]bool{},
	}
}

type Trie struct {
	Root    *Node `json:"root"`
	Actions map[string]bool
}

func (t *Trie) HasAccess(action, scope string) bool {
	if scope == "" {
		return t.Actions[action]
	}
	var hasAccess bool
	t.Root.walk(scope, func(n *Node) bool {
		if n.Actions[action] {
			hasAccess = true
			return true
		}
		return false
	})
	return hasAccess
}

func (t *Trie) Metadata(scope string) Metadata {
	metadata := make(Metadata, 0)
	t.Root.walk(scope, func(n *Node) bool {
		for action := range n.Actions {
			metadata[action] = true
		}
		return false
	})
	return metadata
}

type Node struct {
	Root     bool            `json:"root"`
	Path     string          `json:"path"`
	Actions  map[string]bool `json:"actions"`
	Children map[string]Node `json:"children"`
}

func (n *Node) addNode(action, path, delim string) {
	if path == "" {
		return
	}

	if n.Path == path || path == "*" {
		n.Actions[action] = true
		return
	}

	// no need to child Node when parent has action
	if n.Actions[action] {
		return
	}

	prefix := path
	idx := strings.Index(prefix, delim)
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

	c.addNode(action, path[idx+1:], delim)
}

func (n *Node) walk(path string, walkFn func(n *Node) bool) {
	stop := walkFn(n)
	if stop {
		return
	}

	if n.Path == path {
		return
	}

	prefix := path
	idx := strings.Index(prefix, ":")
	if idx > 0 {
		prefix = path[:idx]
	}

	if c, ok := n.Children[prefix]; ok {
		c.walk(path[idx+1:], walkFn)
	}
}
