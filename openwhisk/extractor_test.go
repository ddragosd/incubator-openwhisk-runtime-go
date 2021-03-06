/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openwhisk

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/h2non/filetype"
	"github.com/stretchr/testify/assert"
)

func sys(cli string, args ...string) {
	os.Chmod(cli, 0755)
	cmd := exec.Command(cli, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	} else {
		fmt.Print(string(out))
	}
}

func exists(dir, filename string) error {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	_, err := os.Stat(path)
	return err
}

func detect(dir, filename string) string {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	file, _ := ioutil.ReadFile(path)
	kind, _ := filetype.Match(file)
	return kind.Extension
}

func TestExtractActionTest_exec(t *testing.T) {
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/x1", "", log)
	sys("_test/build.sh")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action/x1"))
	file, _ := ioutil.ReadFile("_test/exec")
	ap.ExtractAction(&file, "exec")
	assert.Nil(t, exists("./action/x1", "exec"))
}

func TestExtractActionTest_exe(t *testing.T) {
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/x2", "", log)
	sys("_test/build.sh")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action/x2"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec")
	ap.ExtractAction(&file, "exec")
	assert.Equal(t, detect("./action/x2", "exec"), "elf")
}

func TestExtractActionTest_zip(t *testing.T) {
	log, _ := ioutil.TempFile("", "log")
	sys("_test/build.sh")
	ap := NewActionProxy("./action/x3", "", log)
	// cleanup
	assert.Nil(t, os.RemoveAll("./action/x3"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec.zip")
	ap.ExtractAction(&file, "exec")
	assert.Equal(t, detect("./action/x3", "exec"), "elf")
	assert.Nil(t, exists("./action/x3", "etc"))
	assert.Nil(t, exists("./action/x3", "dir/etc"))
}

func TestExtractAction_script(t *testing.T) {
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/x4", "", log)
	buf := []byte("#!/bin/sh\necho ok")
	_, err := ap.ExtractAction(&buf, "exec")
	//fmt.Print(err)
	assert.Nil(t, err)
}

func TestHighestDir(t *testing.T) {
	assert.Equal(t, highestDir("./_test"), 0)
	assert.Equal(t, highestDir("./_test/first"), 3)
	assert.Equal(t, highestDir("./_test/second"), 17)
}
