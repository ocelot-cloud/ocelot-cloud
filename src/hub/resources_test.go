package main

import (
	"bytes"
	"fmt"
	"os"
)

var (
	sampleUser                    = "myuser"
	sampleApp                     = "myapp"
	sampleTag                     = "v0.0.1"
	singleUserDir                 = usersDir + "/" + sampleUser
	appDir                        = singleUserDir + "/" + sampleApp
	sampleFile                    = appDir + fmt.Sprintf("/%s.tar.gz", sampleTag)
	sampleTaggedFileContentBuffer = bytes.NewBuffer([]byte("hello"))
	sampleFileInfo                = &FileInfo{sampleUser, sampleApp, sampleTag, sampleFile}
	sampleMail                    = "testuser@example.com"
	samplePassword                = "mypassword"
)

func cleanup() {
	err := deleteIfExist(dataDir)
	if err != nil {
		Logger.Error("Cleanup: Could not delete dir: %s", dataDir)
		os.Exit(1)
	}
}
