package folder

import (
	"errors"
	"strings"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {

	if name == "" || dst == "" {
		return nil, errors.New("empty folder name in source or destination")
	}

	if name == dst {
		return nil, errors.New("cannot move a folder to itself")
	}

	folders := f.folders

	// find source & dstFolder folder
	srcFolder := Folder{}
	dstFolder := Folder{}
	for _, f := range folders {
		if f.Name == name {
			srcFolder = f
		}

		if f.Name == dst {
			dstFolder = f
		}

		if srcFolder.Name != "" && dstFolder.Name != "" {
			break
		}
	}

	if srcFolder.Name == "" {
		return nil, errors.New("source folder doesn't exist")
	}

	if dstFolder.Name == "" {
		return nil, errors.New("destination folder doesn't exist")
	}

	if srcFolder.OrgId != dstFolder.OrgId {
		return nil, errors.New("cannot move a folder to a different organization")
	}

	children, err := f.GetAllChildFolders(srcFolder.OrgId, srcFolder.Name)
	if err != nil {
		return nil, err
	}

	// check if dstFolder is a child of srcFolder
	for _, c := range children {
		if c.Name == dstFolder.Name {
			return nil, errors.New("cannot move a folder to a child of itself")
		}
	}

	// move folders and update children paths
	res := []Folder{}

	for _, f := range folders {
		if f.Paths == srcFolder.Paths {
			f.Paths = dstFolder.Paths + "." + f.Name
		} else if strings.HasPrefix(f.Paths, srcFolder.Paths) {
			f.Paths = dstFolder.Paths + "." + srcFolder.Name + f.Paths[len(srcFolder.Paths):]
		}

		res = append(res, f)
	}

	return res, nil
}
