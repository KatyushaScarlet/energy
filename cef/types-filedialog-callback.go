//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package cef

import (
	"github.com/energye/energy/v2/common/imports"
	"github.com/energye/golcl/lcl"
)

// Instance 实例
func (m *ICefFileDialogCallback) Instance() uintptr {
	if m == nil {
		return 0
	}
	return uintptr(m.instance)
}

func (m *ICefFileDialogCallback) Free() {
	if m.instance != nil {
		m.base.Free(m.Instance())
		m.instance = nil
	}
}

func (m *ICefFileDialogCallback) IsValid() bool {
	if m == nil || m.instance == nil {
		return false
	}
	return m.instance != nil
}

func (m *ICefFileDialogCallback) Cont(filePaths []string) {
	if !m.IsValid() {
		return
	}
	fps := lcl.NewStrings()
	for _, fp := range filePaths {
		fps.Add(fp)
	}
	imports.Proc(internale_FileDialogCallback_Cont).Call(m.Instance(), fps.Instance())
	fps.Free()
}

func (m *ICefFileDialogCallback) Cancel() {
	if !m.IsValid() {
		return
	}
	imports.Proc(internale_FileDialogCallback_Cancel).Call(m.Instance())
}
