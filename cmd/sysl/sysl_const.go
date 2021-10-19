package main

import "path/filepath"

var (
	projDir = filepath.Join("..", "..")
	syslDir = filepath.Join(projDir, "pkg")
	testDir = filepath.Join(projDir, "tests")
)
