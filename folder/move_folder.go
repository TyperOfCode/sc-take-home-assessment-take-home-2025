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

	if isDescendant(srcFolder.Paths, dstFolder.Paths) {
		return nil, ErrMoveToDescendant
	}

	res, err := moveFolderAndChildren(d, srcFolder, dstFolder)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Copy all the folders and update the necessary paths to move a folder.
// O(n) complexity as we have to return a copy of all the folders anyway.
func moveFolderAndChildren(d *driver, srcFolder *Folder, dstFolder *Folder) ([]Folder, error) {
	res := []Folder{}
	newPrefix := dstFolder.Paths + "." + srcFolder.Name

	srcSegments := strings.Split(srcFolder.Paths, ".")

	for _, name := range d.folderNames {
		node, ok := d.nameToNode[name]
		if !ok {
			return nil, ErrUnexpectedError
		}

		folder := *node.Folder

		if folder.Paths == srcFolder.Paths {
			folder.Paths = newPrefix
		} else if isDescendant(srcFolder.Paths, folder.Paths) {
			folderSegments := strings.Split(folder.Paths, ".")

			folder.Paths = newPrefix + "." + strings.Join(folderSegments[len(srcSegments):], ".")
		}

		res = append(res, folder)
	}

	return res, nil
}

// determines if childPath is a descendant of parentPath
func isDescendant(
	parentPath string,
	childPath string,
) bool {
	parentSegments := strings.Split(parentPath, ".")
	childSegments := strings.Split(childPath, ".")

	if len(childSegments) <= len(parentSegments) {
		return false
	}

	// checks if parent path is a prefix of the child path
	for i, segment := range parentSegments {
		if segment != childSegments[i] {
			return false
		}
	}

	return true
}
