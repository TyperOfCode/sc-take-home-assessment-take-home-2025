package folder

import (
	"errors"
	"strings"
)

func (f *driver) MoveFolder(
	name string,
	dst string,
) ([]Folder, error) {

	if name == "" || dst == "" {
		return nil, errors.New("empty folder name in source or destination")
	}

	if name == dst {
		return nil, errors.New("cannot move a folder to itself")
	}

	srcNode, srcExists := f.nameToNode[name]
	dstNode, dstExists := f.nameToNode[dst]

	if !srcExists {
		return nil, errors.New("source folder doesn't exist")
	}

	if !dstExists {
		return nil, errors.New("destination folder doesn't exist")
	}

	srcFolder, dstFolder := srcNode.Folder, dstNode.Folder

	if srcFolder.OrgId != dstFolder.OrgId {
		return nil, errors.New("cannot move a folder to a different organization")
	}

	if strings.HasPrefix(dstFolder.Paths, srcFolder.Paths) {
		return nil, errors.New("cannot move a folder to its descendant")
	}

	// Copy all the folders and update the necessary paths.
	// O(n) complexity as we have to return a copy of all the folders anyway.
	res := []Folder{}
	for _, names := range f.folderNames {
		node, ok := f.nameToNode[names]
		if !ok {
			continue
		}

		folder := *node.Folder

		if folder.Paths == srcFolder.Paths {
			folder.Paths = dstFolder.Paths + "." + folder.Name
		} else if strings.HasPrefix(folder.Paths, srcFolder.Paths) {
			folder.Paths = dstFolder.Paths + "." + srcFolder.Name + folder.Paths[len(srcFolder.Paths):]
		}

		res = append(res, folder)
	}

	return res, nil
}
