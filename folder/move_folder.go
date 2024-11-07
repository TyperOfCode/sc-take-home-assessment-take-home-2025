package folder

import (
	"strings"
)

func (d *driver) MoveFolder(
	name string,
	dst string,
) ([]Folder, error) {

	if name == "" || dst == "" {
		return nil, ErrInvalidArguments
	}

	if name == dst {
		return nil, ErrMoveToSource
	}

	srcNode, srcExists := d.nameToNode[name]
	dstNode, dstExists := d.nameToNode[dst]

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
	for _, names := range d.folderNames {
		node, ok := d.nameToNode[names]
		if !ok {
			return nil, ErrUnexpectedError
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
