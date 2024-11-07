package folder

import (
	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (d *driver) GetFoldersByOrgID(orgID uuid.UUID) ([]Folder, error) {
	
	res := []Folder{}
	for _, names := range d.folderNames {
		node, ok := d.nameToNode[names]
		if !ok {
			return nil, ErrUnexpectedError
		}

		folder := node.Folder

		if folder.OrgId == orgID {
			res = append(res, *folder)
		}
	}

	return res, nil

}

func (d *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {

	node, folderExists := d.nameToNode[name]
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
