package mover_test

import (
	"testing"
	"time"

	"github.com/watercompany/multiplex/mover"
)

func TestGetFreeDiskSpaceInMB(t *testing.T) {
	sizeInMB := mover.GetFreeDiskSpaceInMB("./")
	t.Logf("Free Space is %v MB", sizeInMB)
}

func TestDirSizeInMB(t *testing.T) {
	sizeInMB, err := mover.DirSizeInMB("./")
	if err != nil {
		t.Fatalf("error %v", err)
	}
	t.Logf("Dir Size is %v MB", sizeInMB)
}

func TestFileSizeInMB(t *testing.T) {
	sizeInMB, err := mover.FileSizeInMB("./testdata/location-a/test1")
	if err != nil {
		t.Fatalf("error %v", err)
	}
	t.Logf("File Size is %v MB", sizeInMB)
}

func TestMoveFile(t *testing.T) {
	src := "./testdata/location-a"
	dest := "./testdata/location-b"
	fileName := "test1"
	err := mover.MoveFile(src, dest, fileName)
	if err != nil {
		t.Fatalf("error moving %v", err)
	}

	time.Sleep(5 * time.Second)
	err = mover.MoveFile(dest, src, fileName)
	if err != nil {
		t.Fatalf("error moving %v", err)
	}
}

func TestFileCount(t *testing.T) {
	path := "./testdata/location-a"

	count, err := mover.FileCount(path)
	if err != nil {
		t.Errorf("want err nil, got %v", err)
	}

	t.Logf("Count=%v", count)
}

func TestFileCountSubString(t *testing.T) {
	path := "./testdata/location-a"

	count, err := mover.FileCountSubString(path, "test")
	if err != nil {
		t.Errorf("want err nil, got %v", err)
	}

	t.Logf("Count=%v", count)
}

func TestGetDirs(t *testing.T) {
	path := "./testdata/location-a"

	dirs, err := mover.GetDirs(path)
	if err != nil {
		t.Errorf("want err nil, got %v", err)
	}

	for i, val := range dirs {
		t.Logf("Dir[%v]=%v", i, val.Name())
		t.Logf("Dir[%v] IsDir=%v", i, val.IsDir())
	}
}

func TestMoveFileToDrive(t *testing.T) {
	src := "./testdata/location-a"
	dest := "/Volumes/KENJEH"
	fileName := "music.mp3"
	err := mover.MoveFile(src, dest, fileName)
	if err != nil {
		t.Fatalf("error moving %v", err)
	}
}

func TestGetFinalDir(t *testing.T) {
	finalDirs := []string{
		"./testdata/location-a",
		"./testdata/location-b",
	}

	mockFileSize := 300
	maxLockFiles := 3
	finalDir, err := mover.GetFinalDir(int64(mockFileSize), finalDirs, maxLockFiles)
	if err != nil {
		t.Errorf("want err nil, got %v", err)
	}

	t.Logf("FinalDir=%v", finalDir)
}
