// Copyright (c) 2021 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tools
// +build tools

package tools

import (
	_ "github.com/google/addlicense"
	_ "github.com/incu6us/goimports-reviser"
	_ "github.com/kisielk/errcheck"
	_ "golang.org/x/lint/golint"
)
