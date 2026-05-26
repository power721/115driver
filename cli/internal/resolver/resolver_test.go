package resolver

import (
	"testing"

	"github.com/power721/115driver/pkg/driver"
)

func TestResolvePath_FallsBackToFileWhenDirLookupReturnsRootID(t *testing.T) {
	client := fakeResolverClient{
		dirIDs: map[string]string{
			"q9tVD1jYR8e626EteJ0qDQ.mp4": RootID,
		},
		filesByDir: map[string][]driver.File{
			RootID: {
				{
					FileID:      "123456",
					Name:        "q9tVD1jYR8e626EteJ0qDQ.mp4",
					IsDirectory: false,
				},
			},
		},
	}

	fileID, isDir, err := ResolvePath(client, "q9tVD1jYR8e626EteJ0qDQ.mp4")
	if err != nil {
		t.Fatalf("ResolvePath returned error: %v", err)
	}
	if isDir {
		t.Fatalf("ResolvePath should treat file path as file")
	}
	if fileID != "123456" {
		t.Fatalf("unexpected file ID: %s", fileID)
	}
}

func TestResolvePath_RootStillResolvesToDirectory(t *testing.T) {
	fileID, isDir, err := ResolvePath(fakeResolverClient{}, "/")
	if err != nil {
		t.Fatalf("ResolvePath returned error: %v", err)
	}
	if !isDir {
		t.Fatalf("root path should resolve as directory")
	}
	if fileID != RootID {
		t.Fatalf("unexpected root ID: %s", fileID)
	}
}

type fakeResolverClient struct {
	dirIDs     map[string]string
	filesByDir map[string][]driver.File
}

func (f fakeResolverClient) DirName2CID(dir string) (*driver.APIGetDirIDResp, error) {
	id := f.dirIDs[dir]
	return &driver.APIGetDirIDResp{
		CategoryID: driver.IntString(id),
	}, nil
}

func (f fakeResolverClient) List(dirID string, _ ...driver.ListOption) (*[]driver.File, error) {
	files := f.filesByDir[dirID]
	return &files, nil
}
