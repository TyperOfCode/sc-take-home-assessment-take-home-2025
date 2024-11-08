package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	orgIdToFetch := uuid.FromStringOrNil(folder.DefaultOrgID)
	orgId1 := uuid.Must(uuid.NewV4())
	orgId2 := uuid.Must(uuid.NewV4())
	orgId3 := uuid.Must(uuid.NewV4())

	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		expect  []folder.Folder
	}{

		//-------- non-error cases
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
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "beta", OrgId: orgId2, Paths: "beta"},
				{Name: "c", OrgId: orgId3, Paths: "alpha.c"},
			},
			expect: []folder.Folder{},
		},

		{
			name:  "Correctly fetches folders by orgID.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "beta.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "alpha.e"},
			},
			expect: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "beta.c.d"},
			},
		},

		{
			name:  "Nested folders of different organisation.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "beta"},
				{Name: "c", OrgId: orgId2, Paths: "alpha.c"},
				{Name: "d", OrgId: orgId2, Paths: "beta.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "alpha.e"},
			},
			expect: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "beta"},
			},
		},

		{
			name:  "Should fetch children folders of non-matching orgID.",
			orgID: orgIdToFetch,
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "beta", OrgId: orgId2, Paths: "beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "beta.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "alpha.e"},
			},
			expect: []folder.Folder{
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "beta.c.d"},
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f, err := folder.NewDriver(tc.folders)
			assert.NoError(t, err, "unexpected error")

			result, err := f.GetFoldersByOrgID(tc.orgID)

			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expect, result, "unexpected result")
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	orgIdToFetch := uuid.FromStringOrNil(folder.DefaultOrgID)
	orgId1 := uuid.Must(uuid.NewV4())

	tests := [...]struct {
		testName    string
		orgID       uuid.UUID
		folderName  string
		folders     []folder.Folder
		expect      []folder.Folder
		expectError bool
	}{

		//-------- non-error cases

		{
			testName:   "No children of the folder.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "beta.c"},
			},
			expect:      []folder.Folder{},
			expectError: false,
		},

		{
			testName:   "No children of the folder with the necessary orgId.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgId1, Paths: "alpha.c"},
			},
			expect:      []folder.Folder{},
			expectError: false,
		},

		{
			testName:   "Shouldnt match with subset of folder name.",
			orgID:      orgIdToFetch,
			folderName: "alp",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgId1, Paths: "alpha.c"},
				{Name: "alp", OrgId: orgIdToFetch, Paths: "alp"},
			},
			expect:      []folder.Folder{},
			expectError: false,
		},

		{
			testName:   "Fetches all direct children of the folder.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
			},
			expect: []folder.Folder{
				{Name: "beta", OrgId: orgIdToFetch, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
			},
			expectError: false,
		},

		{
			testName:   "case sensitive when fetching folder name.",
			orgID:      orgIdToFetch,
			folderName: "Alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "Alpha", OrgId: orgIdToFetch, Paths: "Alpha"},
				{Name: "Beta", OrgId: orgIdToFetch, Paths: "Alpha.Beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "Alpha.c"},
			},
			expect: []folder.Folder{
				{Name: "Beta", OrgId: orgIdToFetch, Paths: "Alpha.Beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "Alpha.c"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all direct children of non-root folder",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "parent", OrgId: orgIdToFetch, Paths: "parent"},
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "parent.alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "parent.alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.alpha.c"},
			},
			expect: []folder.Folder{
				{Name: "beta", OrgId: orgIdToFetch, Paths: "parent.alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.alpha.c"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all direct children of the folder, mixed OrgIDs.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "alpha.d"},
			},
			expect: []folder.Folder{
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "alpha.d"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all descendants, not just children.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.beta.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "alpha.beta.c.d"},
			},
			expect: []folder.Folder{
				{Name: "beta", OrgId: orgIdToFetch, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.beta.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "alpha.beta.c.d"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all descendants, not just children, non-root.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "parent", OrgId: orgIdToFetch, Paths: "parent"},
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "parent.alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "parent.alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.alpha.beta.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "parent.alpha.beta.c.d"},
			},
			expect: []folder.Folder{
				{Name: "beta", OrgId: orgIdToFetch, Paths: "parent.alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "parent.alpha.beta.c"},
				{Name: "d", OrgId: orgIdToFetch, Paths: "parent.alpha.beta.c.d"},
			},
			expectError: false,
		},

		{
			testName:   "Fetches all descendants, not just children, mixed orgIDs.",
			orgID:      orgIdToFetch,
			folderName: "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.beta.c"},
				{Name: "d", OrgId: orgId1, Paths: "alpha.beta.c.d"},
			},
			expect: []folder.Folder{
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.beta.c"},
			},
			expectError: false,
		},

		//-------- errorful behaviour tests
		{
			testName:    "Error on empty folders.",
			orgID:       orgIdToFetch,
			folderName:  "alpha",
			folders:     []folder.Folder{},
			expect:      nil,
			expectError: true,
		},

		{
			testName:   "Error empty folder name.",
			orgID:      orgIdToFetch,
			folderName: "",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgIdToFetch, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "alpha.c"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName:   "Error on folder does not exist.",
			orgID:      orgIdToFetch,
			folderName: "x",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgIdToFetch, Paths: "c"},
				{Name: "d", OrgId: orgId1, Paths: "d"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName:   "Error on folder isnt in org",
			orgID:      orgIdToFetch,
			folderName: "beta",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgIdToFetch, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
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

			f, err := folder.NewDriver(tc.folders)
			assert.NoError(t, err, "unexpected error")

			result, err := f.GetAllChildFolders(tc.orgID, tc.folderName)

			if tc.expectError {
				assert.Error(t, err, "expected error")
				return
			} else {
				assert.NoError(t, err, "unexpected error")
			}

			assert.Equal(t, tc.expect, result, "unexpected result")
		})
	}

}
