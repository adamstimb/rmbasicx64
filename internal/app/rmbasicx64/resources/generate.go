// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate file2byteslice -package=audio -input=./audio/jab.wav -output=./audio/jab.go -var=Jab_wav
//go:generate file2byteslice -package=audio -input=./audio/jump.ogg -output=./audio/jump.go -var=Jump_ogg
//go:generate file2byteslice -package=audio -input=./audio/ragtime.mp3 -output=./audio/ragtime.mp3.go -var=Ragtime_mp3
//go:generate file2byteslice -package=audio -input=./audio/ragtime.ogg -output=./audio/ragtime.ogg.go -var=Ragtime_ogg
//go:generate gofmt -s -w .

package resources

import (
	// Dummy imports for go.mod for some Go files with 'ignore' tags. For example, `go mod tidy` does not
	// recognize Go files with 'ignore' build tag.
	//
	// Note that this affects only importing this package, but not 'file2byteslice' commands in //go:generate.
	_ "github.com/hajimehoshi/file2byteslice"
)
