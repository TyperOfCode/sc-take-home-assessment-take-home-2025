package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_folder_MoveFolder(t *testing.T) {
	t.Parallel()

	orgId1 := uuid.FromStringOrNil(folder.DefaultOrgID)
	orgId2 := uuid.Must(uuid.NewV4())

	tests := [...]struct {
		testName    string
		src         string
		dst         string
		folders     []folder.Folder
		expect      []folder.Folder
		expectError bool
	}{

		//-------- non-error cases

		{
			testName: "Move to new folder without children",
			src:      "alpha",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "beta"},
			},
			expect: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "beta.alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "beta"},
			},
			expectError: false,
		},

		{
			testName: "Case sensitivity",
			src:      "Alpha",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "Alpha", OrgId: orgId1, Paths: "Alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "Beta", OrgId: orgId1, Paths: "Alpha.Beta"},
			},
			expect: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "Alpha", OrgId: orgId1, Paths: "alpha.beta.Alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "Beta", OrgId: orgId1, Paths: "alpha.beta.Alpha.Beta"},
			},
			expectError: false,
		},

		{
			testName: "Move out of folder to parent",
			src:      "beta",
			dst:      "parent",
			folders: []folder.Folder{
				{Name: "parent", OrgId: orgId1, Paths: "parent"},
				{Name: "alpha", OrgId: orgId1, Paths: "parent.alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "parent.alpha.beta"},
			},
			expect: []folder.Folder{
				{Name: "parent", OrgId: orgId1, Paths: "parent"},
				{Name: "alpha", OrgId: orgId1, Paths: "parent.alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "parent.beta"},
			},
			expectError: false,
		},

		{
			testName: "Move to sibling folder",
			src:      "beta",
			dst:      "alpha",
			folders: []folder.Folder{
				{Name: "parent", OrgId: orgId1, Paths: "parent"},
				{Name: "alpha", OrgId: orgId1, Paths: "parent.alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "parent.beta"},
			},
			expect: []folder.Folder{
				{Name: "parent", OrgId: orgId1, Paths: "parent"},
				{Name: "alpha", OrgId: orgId1, Paths: "parent.alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "parent.alpha.beta"},
			},
			expectError: false,
		},

		{
			testName: "Move to new folder with children",
			src:      "alpha",
			dst:      "forest",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgId1, Paths: "alpha.beta.c"},
				{Name: "d", OrgId: orgId1, Paths: "alpha.beta.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "alpha.beta.e"},
				{Name: "forest", OrgId: orgId1, Paths: "forest"},
			},
			expect: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "forest.alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "forest.alpha.beta"},
				{Name: "c", OrgId: orgId1, Paths: "forest.alpha.beta.c"},
				{Name: "d", OrgId: orgId1, Paths: "forest.alpha.beta.c.d"},
				{Name: "e", OrgId: orgId1, Paths: "forest.alpha.beta.e"},
				{Name: "forest", OrgId: orgId1, Paths: "forest"},
			},
			expectError: false,
		},

		{
			testName: "Shouldn't match with subset of folder name",
			src:      "alp",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "alp", OrgId: orgId1, Paths: "alp"},
				{Name: "betaone", OrgId: orgId1, Paths: "alpha.beta.betaone"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
			},
			expect: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "alp", OrgId: orgId1, Paths: "alpha.beta.alp"},
				{Name: "betaone", OrgId: orgId1, Paths: "alpha.beta.betaone"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
			},
			expectError: false,
		},

		//-------- errorful cases

		{
			testName:    "Empty folders.",
			src:         "alpha",
			dst:         "beta",
			folders:     []folder.Folder{},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Source cannot be empty.",
			src:      "",
			dst:      "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
			},
			expect:      nil,
			expectError: true,
		},

		// (this is assuming that dst cannot be empty)
		// Should it move to the root instead?
		{
			testName: "Destination cannot be empty.",
			src:      "alpha",
			dst:      "",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Source doesnt exist.",
			src:      "alpha",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "beta", OrgId: orgId1, Paths: "beta"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Destination doesnt exist.",
			src:      "alpha",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Both source and destination dont exist.",
			src:      "alpha",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "c", OrgId: orgId1, Paths: "c"},
				{Name: "d", OrgId: orgId1, Paths: "d"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Can't copy across organisations.",
			src:      "alpha",
			dst:      "beta",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "beta", OrgId: orgId2, Paths: "beta"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Can't move a folder to itself",
			src:      "alpha",
			dst:      "alpha",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
			},
			expect:      nil,
			expectError: true,
		},

		{
			testName: "Can't move a folder to a child of itself",
			src:      "alpha",
			dst:      "c",
			folders: []folder.Folder{
				{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
				{Name: "beta", OrgId: orgId1, Paths: "alpha.beta"},
				{Name: "c", OrgId: orgId1, Paths: "alpha.beta.c"},
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

			result, err := f.MoveFolder(tc.src, tc.dst)

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
