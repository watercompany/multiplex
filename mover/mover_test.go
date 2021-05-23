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
