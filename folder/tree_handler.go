package folder

import "strings"

func BuildFolderTree(folders []Folder, nameOfFolders []string) (map[string]*FolderNode, error) {
	nameToNode := make(map[string]*FolderNode)

	// assumption: the names of the files are unique.
	// populate map with nodes
	for i, f := range folders {
		nameToNode[f.Name] = &FolderNode{Folder: &folders[i]}
	}

	// populate the existing nodes with children
	for _, name := range nameOfFolders {
		node, ok := nameToNode[name]
		if !ok {
			continue
		}

		parts := strings.Split(node.Folder.Paths, ".")

		if len(parts) <= 1 {
			continue
		}

		parentName := parts[len(parts)-2]
		parentNode, ok := nameToNode[parentName]
		if ok {
			parentNode.Children = append(parentNode.Children, node)
		}
	}

	return nameToNode, nil
}
