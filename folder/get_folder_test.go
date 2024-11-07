package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	orgIdToFetch := uuid.Must(uuid.NewV4())
	orgId1 := uuid.Must(uuid.NewV4())
	orgId2 := uuid.Must(uuid.NewV4())
	orgId3 := uuid.Must(uuid.NewV4())

	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		expect  []folder.Folder
	}{
		{
			name:    "Empty folders.",
			orgID:   orgIdToFetch,
			folders: []folder.Folder{},
			expect:  []folder.Folder{},
		},
		{
			name:  "No folders with the orgID.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "a", OrgId: orgId1, Paths: "a"},
				{Name: "b", OrgId: orgId2, Paths: "b"},
				{Name: "c", OrgId: orgId3, Paths: "a.c"},
			},
			expect: []folder.Folder{},
		},
		{
			name:  "Correctly fetches folders by orgID.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "b.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "a.e"},
			},
			expect: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "b.c.d"},
			},
		},
		{
			name:  "Nested folders of different organisation.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "b"},
				{Name: "c", OrgId: orgId2, Paths: "a.c"},
				{Name: "d", OrgId: orgId2, Paths: "b.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "a.e"},
			},
			expect: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "b"},
			},
		},
		{
			name:  "Should fetch children folders of non-matching orgID.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "a", OrgId: orgId1, Paths: "a"},
				{Name: "b", OrgId: orgId2, Paths: "b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "b.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "a.e"},
			},
			expect: []folder.Folder{
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "b.c.d"},
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := folder.NewDriver(tc.folders)

			result := f.GetFoldersByOrgID(tc.orgID)

			assert.Equal(t, tc.expect, result, "unexpected result")
		})
	}
}