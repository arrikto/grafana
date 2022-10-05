package accesscontrol

import (
	"strings"
)

const (
	pathDelim  = "/"
	scopeDelim = ":"
)

func TrieFromPermissions(permissions []Permission) *Trie {
	t := newTrie()
	for _, p := range permissions {
		t.Actions[p.Action] = true
		t.Root.addNode(p.Action, p.Scope, scopeDelim)
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
	Root    *Node           `json:"root"`
	Actions map[string]bool `json:"actions"`
}

func (t *Trie) HasAccess(action, scope string) bool {
	if scope == "" {
		return t.Actions[action]
	}
	var hasAccess bool
	t.Root.walkPath(scope, func(n *Node) bool {
		if n.Actions[action] {
			hasAccess = true
			return true
		}
		return false
	})
	return hasAccess
}

func (t *Trie) Identifiers(action, prefix string) (bool, []string) {
	var hasWildcard bool
	t.Root.walkPath(prefix, func(n *Node) bool {
		if n.Actions[action] {
			hasWildcard = true
			return true
		}
		return false
	})

	if hasWildcard {
		return true, nil
	}

	var identifiers []string
	t.Root.walkPrefix(prefix, func(n *Node) bool {
		identifiers = append(identifiers, n.Path)
		return false
	})
	return false, identifiers
}

func (t *Trie) Metadata(scope string) map[string]bool {
	metadata := Metadata{}
	t.Root.walkPath(scope, func(n *Node) bool {
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
		c = Node{
			Path:     prefix,
			Actions:  map[string]bool{},
			Children: map[string]Node{},
		}
		n.Children[prefix] = c
	}

	c.addNode(action, path[idx+1:], delim)
}

func (n *Node) walkPath(path string, walkFn func(n *Node) bool) {
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
		c.walkPath(path[idx+1:], walkFn)
	}
}

// walkPrefix walks every node under prefix
func (n *Node) walkPrefix(prefix string, walkFn func(n *Node) bool) {
	path := prefix
	idx := strings.Index(path, ":")
	if idx > 0 {
		path = prefix[:idx]
	}

	if c, ok := n.Children[path]; ok {
		c.walkPrefix(prefix[idx+1:], walkFn)
	}

	if path == "" {
		for _, c := range n.Children {
			walkFn(&c)
		}
	}
}
