package folder

import (
	"errors"

	"github.com/gofrs/uuid"
)

// errors
var ErrFolderDoesNotExist = errors.New("folder doesn't exist")
var ErrFolderDoesNotExistInOrg = errors.New("folder doesn't exist in the specified organization")

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {

	res := []Folder{}
	for _, names := range f.folderNames {
		f := f.nameToNode[names].Folder

		if f.OrgId == orgID {
			res = append(res, *f)
		}
	}

	return res

}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {

	node, folderExists := f.nameToNode[name]
	folder := node.Folder

	if !folderExists {
		return nil, ErrFolderDoesNotExist
	}

	if folder.OrgId != orgID {
		return nil, ErrFolderDoesNotExistInOrg
	}

	res := []Folder{}
	getDescendants(&res, node, orgID)

	return res, nil
}

// adds all descendants of node to res
func getDescendants(res *[]Folder, node *FolderNode, orgID uuid.UUID) {
	for _, child := range node.Children {
		if child.Folder.OrgId == orgID {
			*res = append(*res, *child.Folder)
		}

		getDescendants(res, child, orgID)
	}
}
