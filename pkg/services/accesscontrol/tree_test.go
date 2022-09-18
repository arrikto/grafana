package accesscontrol

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	content, err := os.ReadFile("./data2.json")
	require.NoError(t, err)

	var permissions []Permission
	err = json.Unmarshal(content, &permissions)
	require.NoError(t, err)

	tree := BuildPermissionTrie(permissions)
	data, err := json.MarshalIndent(tree, "", " ")
	require.NoError(t, err)
	fmt.Println(string(data))
}

func TestTree2(t *testing.T) {
	permissions := []Permission{
		{Action: "datasources:read", Scope: "datasources:*"},
		{Action: "datasources:read", Scope: "datasources:uid:123"},
		{Action: "datasources:write", Scope: "datasources:uid:123"},
	}

	tree := BuildPermissionTrie(permissions)
	data, err := json.MarshalIndent(tree, "", " ")
	require.NoError(t, err)
	fmt.Println(string(data))
}

func TestTrie
