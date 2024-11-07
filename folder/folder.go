package folder

import "github.com/gofrs/uuid"

type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) ([]Folder, error)
	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]Folder, error)
}

type FolderNode struct {
	Folder   *Folder
	Children []*FolderNode
}

type driver struct {
	folderNames []string
	nameToNode  map[string]*FolderNode
}

func NewDriver(folders []Folder) (IDriver, error) {
	var folderNames []string
	for _, folder := range folders {
		folderNames = append(folderNames, folder.Name)
	}

	nameToNode, err := BuildFolderTree(folders, folderNames)
	if err != nil {
		return nil, err
	}

	return &driver{
		folderNames: folderNames,
		nameToNode:  nameToNode,
	}, nil
}
