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

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	orgIdToFetch := uuid.Must(uuid.NewV4())
	orgId1 := uuid.Must(uuid.NewV4())

	tests := [...]struct {
		testName    string
		orgID       uuid.UUID
		folderName  string
		folders     []folder.Folder
		expect      []folder.Folder
		expectError bool
	}{

		{
			testName:   "No children of the folder.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "b.c"},
			},
			expect:      []folder.Folder{},
			expectError: false,
		},

		{
			testName:   "No children of the folder with the necessary orgId.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgId1, Paths: "a.b"},
				{Name: "c", OrgId: orgId1, Paths: "a.c"},
			},
			expect:      []folder.Folder{},
			expectError: false,
		},

		{
			testName:   "Fetches all direct children of the folder.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
			},
			expect: []folder.Folder{
				{Name: "b", OrgId: orgIdToFetch, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all direct children of non-root folder",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "parent", OrgId: orgIdToFetch, Paths: "parent"},
				{Name: "a", OrgId: orgIdToFetch, Paths: "parent.a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "parent.a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.a.c"},
			},
			expect: []folder.Folder{
				{Name: "b", OrgId: orgIdToFetch, Paths: "parent.a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.a.c"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all direct children of the folder, mixed OrgIDs.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgId1, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "a.d"},
			},
			expect: []folder.Folder{
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "a.d"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all descendents, not just children.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.b.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "a.b.c.d"},
			},
			expect: []folder.Folder{
				{Name: "b", OrgId: orgIdToFetch, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.b.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "a.b.c.d"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all descendents, not just children, non-root.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "parent", OrgId: orgIdToFetch, Paths: "parent"},
				{Name: "a", OrgId: orgIdToFetch, Paths: "parent.a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "parent.a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.a.b.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "parent.a.b.c.d"},
			},
			expect: []folder.Folder{
				{Name: "b", OrgId: orgIdToFetch, Paths: "parent.a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.a.b.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "parent.a.b.c.d"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all descendents, not just children, mixed orgIDs.",
			orgID:      orgIdToFetch,
			folderName: "a",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgId1, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.b.c"},
				{Name: "d", OrgId: orgId1, Paths: "a.b.c.d"},
			},
			expect: []folder.Folder{
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.b.c"},
			},
			expectError: false,
		},

		// errorful behaviour tests
		{
			testName:    "Error on empty folders.",
			orgID:       orgIdToFetch,
			folderName:  "a",
			folders:     []folder.Folder{},
			expect:      nil,
			expectError: true,
		},

		{
			testName:   "Error empty folder name.",
			orgID:      orgIdToFetch,
			folderName: "",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgIdToFetch, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "a.c"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName:   "Error on folder does not exist.",
			orgID:      orgIdToFetch,
			folderName: "x",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgId1, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "c"},
				{Name: "d", OrgId: orgId1, Paths: "d"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName:   "Error on folder isnt in org",
			orgID:      orgIdToFetch,
			folderName: "b",
			folders: []folder.Folder{
				{Name: "a", OrgId: orgIdToFetch, Paths: "a"},
				{Name: "b", OrgId: orgId1, Paths: "a.b"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "c"},
				{Name: "d", OrgId: orgId1, Paths: "d"},
			},
			expect:      nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			f := folder.NewDriver(tc.folders)

			result, err := f.GetAllChildFolders(tc.orgID, tc.folderName)

			if tc.expectError {
				assert.Error(t, err, "expected error")
				return
			}

			assert.Equal(t, tc.expect, result, "unexpected result")
		})
	}

}
