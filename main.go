package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

	res := folder.GetAllFolders()

	// example usage
	folderDriver, err := folder.NewDriver(res)
	if err != nil {
		panic(err)
	}

	_, err = folderDriver.GetFoldersByOrgID(orgID)
	if err != nil {
		panic(err)
	}

	newFolder, err := folderDriver.MoveFolder("fast-watchmen", "nearby-secret")

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nFolders after moving folder")
	folder.PrettyPrint(newFolder)

	// folder.PrettyPrint(res)
	// fmt.Printf("\n Folders for orgID: %s", orgID)
	// folder.PrettyPrint(orgFolder)

	// fmt.Printf("\nNumber: %d\n", len(orgFolder))
}
