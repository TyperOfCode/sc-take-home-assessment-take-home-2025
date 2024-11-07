package folder

import "strings"

type FolderNode struct {
	Folder   *Folder
	Children []*FolderNode
}

func BuildFolderTree(folders []Folder) (map[string]*FolderNode, error) {
	nameToNode := make(map[string]*FolderNode)

	for i, f := range folders {
		nameToNode[f.Name] = &FolderNode{Folder: &folders[i]}
	}

	for _, node := range nameToNode {
		parts := strings.Split(node.Folder.Paths, ".")

		if len(parts) <= 1 {
			continue
		}

		parentName := strings.Join(parts[:len(parts)-1], ".")
		parentNode, ok := nameToNode[parentName]
		if ok {
			parentNode.Children = append(parentNode.Children, node)
		}
	}

	return nameToNode, nil
}
