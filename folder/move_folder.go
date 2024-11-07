package folder

import (
	"errors"
	"strings"
)

// errors
var ErrInvalidArguments = errors.New("empty folder name in source or destination")
var ErrMoveToSource = errors.New("cannot move a folder to itself")
var ErrMoveToDescendant = errors.New("cannot move a folder to its descendant")
var ErrMoveToDifferentOrg = errors.New("cannot move a folder to a different organization")

var ErrSourceDoesNotExist = errors.New("source folder doesn't exist")
var ErrDestDoesNotExist = errors.New("destination folder doesn't exist")


func (f *driver) MoveFolder(
	name string,
	dst string,
) ([]Folder, error) {

	if name == "" || dst == "" {
		return nil, ErrInvalidArguments
	}

	if name == dst {
		return nil, ErrMoveToSource
	}

	srcNode, srcExists := f.nameToNode[name]
	dstNode, dstExists := f.nameToNode[dst]

	if !srcExists {
		return nil, ErrSourceDoesNotExist
	}

	if !dstExists {
		return nil, ErrDestDoesNotExist
	}

	srcFolder, dstFolder := srcNode.Folder, dstNode.Folder

	if srcFolder.OrgId != dstFolder.OrgId {
		return nil, ErrMoveToDifferentOrg
	}

	if strings.HasPrefix(dstFolder.Paths, srcFolder.Paths) {
		return nil, ErrMoveToDescendant
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
