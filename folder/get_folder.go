package folder

import (
	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	// folders := f.folders

	// res := []Folder{}
	// for _, f := range folders {
	// 	if f.OrgId == orgID {
	// 		res = append(res, f)
	// 	}
	// }

	return []Folder{}

}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// folders := f.folders

	// // (assuming all the folder names are unique)
	// // finds the folder path with matching name and orgid.
	// path := ""
	// folderExists := false
	// for _, f := range folders {
	// 	if f.Name == name {
	// 		folderExists = true

	// 		if f.OrgId == orgID {
	// 			path = f.Paths
	// 		}
	// 		break
	// 	}
	// }

	// if !folderExists {
	// 	return nil, errors.New("folder doesn't exist")
	// }

	// if path == "" {
	// 	return nil, errors.New("folder doesn't exist in the specified organization")
	// }

	// // finds all the children/descendants of parent, ignoring the parent.
	// res := []Folder{}
	// for _, f := range folders {
	// 	if f.OrgId == orgID && strings.HasPrefix(f.Paths, path) && f.Paths != path {
	// 		res = append(res, f)
	// 	}
	// }

	return []Folder{}, nil
}
