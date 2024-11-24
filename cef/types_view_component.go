//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https//www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package cef

import (
	"github.com/energye/energy/v2/cef/internal/def"
	"github.com/energye/energy/v2/common/imports"
	"github.com/energye/energy/v2/types"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/api"
	"unsafe"
)

// TCEFViewComponent
type TCEFViewComponent struct {
	lcl.IComponent
	instance unsafe.Pointer
}

// ViewComponentRef -> TCEFViewComponent
var ViewComponentRef viewComponent

type viewComponent uintptr

func (*viewComponent) New(AOwner lcl.IComponent) *TCEFViewComponent {
	var result uintptr
	imports.Proc(def.ViewComponent_Create).Call(AOwner.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &TCEFViewComponent{instance: getInstance(result)}
	}
	return nil
}

// Instance 实例
func (m *TCEFViewComponent) Instance() uintptr {
	if m == nil {
		return 0
	}
	return uintptr(m.instance)
}

func (m *TCEFViewComponent) Free() {
	if m.instance != nil {
		imports.Proc(def.ViewComponent_Free).Call(m.Instance())
		m.instance = nil
	}
}

func (m *TCEFViewComponent) IsValid() bool {
	if m == nil || m.instance == nil {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetValid).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) ToStringEx(includeChildren bool) string {
	if !m.IsValid() {
		return ""
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ToStringEx).Call(m.Instance(), api.PascalBool(includeChildren))
	return api.GoStr(r1)
}

func (m *TCEFViewComponent) IsSame(that *ICefView) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_IsSame).Call(m.Instance(), that.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) SizeToPreferredSize() {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SizeToPreferredSize).Call(m.Instance())
}

func (m *TCEFViewComponent) InvalidateLayout() {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_InvalidateLayout).Call(m.Instance())
}

func (m *TCEFViewComponent) RequestFocus() {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_RequestFocus).Call(m.Instance())
}

func (m *TCEFViewComponent) ConvertPointToScreen(point TCefPoint) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ConvertPointToScreen).Call(m.Instance(), uintptr(unsafe.Pointer(&point)))
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) ConvertPointFromScreen(point TCefPoint) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ConvertPointFromScreen).Call(m.Instance(), uintptr(unsafe.Pointer(&point)))
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) ConvertPointToWindow(point TCefPoint) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ConvertPointToWindow).Call(m.Instance(), uintptr(unsafe.Pointer(&point)))
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) ConvertPointFromWindow(point TCefPoint) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ConvertPointFromWindow).Call(m.Instance(), uintptr(unsafe.Pointer(&point)))
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) ConvertPointToView(view *ICefView, point TCefPoint) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ConvertPointToView).Call(m.Instance(), view.Instance(), uintptr(unsafe.Pointer(&point)))
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) ConvertPointFromView(view *ICefView, point TCefPoint) bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_ConvertPointFromView).Call(m.Instance(), view.Instance(), uintptr(unsafe.Pointer(&point)))
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) Initialized() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetInitialized).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) AsView() *ICefView {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetAsView).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefView{instance: getInstance(result)}
	}
	return nil
}

func (m *TCEFViewComponent) AsBrowserView() *ICefBrowserView {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetAsBrowserView).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefBrowserView{&ICefView{instance: getInstance(result)}}
	}
	return nil
}

func (m *TCEFViewComponent) AsButton() *ICefButton {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetAsButton).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefButton{&ICefView{instance: getInstance(result)}}
	}
	return nil
}

func (m *TCEFViewComponent) AsPanel() *ICefPanel {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetAsPanel).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefPanel{&ICefView{instance: getInstance(result)}}
	}
	return nil
}

func (m *TCEFViewComponent) AsScrollView() *ICefScrollView {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetAsScrollView).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefScrollView{&ICefView{instance: getInstance(result)}}
	}
	return nil
}

func (m *TCEFViewComponent) AsTextfield() *ICefTextfield {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetAsTextfield).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefTextfield{&ICefView{instance: getInstance(result)}}
	}
	return nil
}

func (m *TCEFViewComponent) GetViewForID(id int32) *ICefView {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetViewForID).Call(m.Instance(), uintptr(id), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefView{instance: getInstance(result)}
	}
	return nil
}

func (m *TCEFViewComponent) IsAttached() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetAttached).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) GetDelegate() *ICefViewDelegate {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetDelegate).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefViewDelegate{instance: getInstance(result)}
	}
	return nil
}

func (m *TCEFViewComponent) GetWindow() *ICefWindow {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetWindow).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefWindow{&ICefPanel{&ICefView{instance: getInstance(result)}}}
	}
	return nil
}

func (m *TCEFViewComponent) GetParentView() *ICefView {
	if !m.IsValid() {
		return nil
	}
	var result uintptr
	imports.Proc(def.ViewComponent_GetParentView).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefView{instance: getInstance(result)}
	}
	return nil
}

func (m *TCEFViewComponent) BoundsInScreen() (bounds TCefRect) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetBoundsInScreen).Call(m.Instance(), uintptr(unsafe.Pointer(&bounds)))
	return
}

func (m *TCEFViewComponent) GetPreferredSize() (size TCefSize) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetPreferredSize).Call(m.Instance(), uintptr(unsafe.Pointer(&size)))
	return
}

func (m *TCEFViewComponent) MinimumSize() (size TCefSize) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetMinimumSize).Call(m.Instance(), uintptr(unsafe.Pointer(&size)))
	return
}

func (m *TCEFViewComponent) MaximumSize() (size TCefSize) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetMaximumSize).Call(m.Instance(), uintptr(unsafe.Pointer(&size)))
	return
}

func (m *TCEFViewComponent) IsVisible() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetVisible).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) SetVisible(visible bool) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetVisible).Call(m.Instance(), api.PascalBool(visible))
}

func (m *TCEFViewComponent) IsDrawn() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetDrawn).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) IsEnabled() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetEnabled).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) SetEnabled(enabled bool) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetEnabled).Call(m.Instance(), api.PascalBool(enabled))
}

func (m *TCEFViewComponent) IsFocusable() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetFocusable).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) SetFocusable(focusable bool) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetFocusable).Call(m.Instance(), api.PascalBool(focusable))
}

func (m *TCEFViewComponent) AccessibilityFocusable() bool {
	if !m.IsValid() {
		return false
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetAccessibilityFocusable).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *TCEFViewComponent) GetBackgroundColor() (color types.TCefColor) {
	if !m.IsValid() {
		return 0
	}
	imports.Proc(def.ViewComponent_GetBackgroundColor).Call(m.Instance(), uintptr(unsafe.Pointer(&color)))
	return
}

func (m *TCEFViewComponent) SetBackgroundColor(color types.TCefColor) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetBackgroundColor).Call(m.Instance(), uintptr(color))
}

func (m *TCEFViewComponent) GetID() int32 {
	if !m.IsValid() {
		return 0
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetID).Call(m.Instance())
	return int32(r1)
}

func (m *TCEFViewComponent) SetID(id int32) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetID).Call(m.Instance(), uintptr(id))
}

func (m *TCEFViewComponent) GroupID() int32 {
	if !m.IsValid() {
		return 0
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetGroupID).Call(m.Instance())
	return int32(r1)
}

func (m *TCEFViewComponent) SetGroupID(groupId int32) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetGroupID).Call(m.Instance(), uintptr(groupId))
}

func (m *TCEFViewComponent) Bounds() (bounds TCefRect) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetBounds).Call(m.Instance(), uintptr(unsafe.Pointer(&bounds)))
	return
}

func (m *TCEFViewComponent) SetBounds(bounds TCefRect) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetBounds).Call(m.Instance(), uintptr(unsafe.Pointer(&bounds)))
}

func (m *TCEFViewComponent) GetSize() (size TCefSize) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetSize).Call(m.Instance(), uintptr(unsafe.Pointer(&size)))
	return
}

func (m *TCEFViewComponent) SetSize(size TCefSize) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetSize).Call(m.Instance(), uintptr(unsafe.Pointer(&size)))
}

func (m *TCEFViewComponent) GetPosition() (point TCefPoint) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_GetPosition).Call(m.Instance(), uintptr(unsafe.Pointer(&point)))
	return
}

func (m *TCEFViewComponent) SetPosition(position TCefPoint) {
	if !m.IsValid() {
		return
	}
	imports.Proc(def.ViewComponent_SetPosition).Call(m.Instance(), uintptr(unsafe.Pointer(&position)))
}

func (m *TCEFViewComponent) GetTypeString() string {
	if !m.IsValid() {
		return ""
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetTypeString).Call(m.Instance())
	return api.GoStr(r1)
}

func (m *TCEFViewComponent) GetHeightForWidth(width int32) int32 {
	if !m.IsValid() {
		return 0
	}
	r1, _, _ := imports.Proc(def.ViewComponent_GetHeightForWidth).Call(m.Instance(), uintptr(width))
	return int32(r1)
}

func (m *TCEFViewComponent) SetOnGetPreferredSize(fn viewOnGetPreferredSize) {
	imports.Proc(def.ViewComponent_SetOnGetPreferredSize).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnGetMinimumSize(fn viewOnGetMinimumSize) {
	imports.Proc(def.ViewComponent_SetOnGetMinimumSize).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnGetMaximumSize(fn viewOnGetMaximumSize) {
	imports.Proc(def.ViewComponent_SetOnGetMaximumSize).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnGetHeightForWidth(fn viewOnGetHeightForWidth) {
	imports.Proc(def.ViewComponent_SetOnGetHeightForWidth).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnParentViewChanged(fn viewOnParentViewChanged) {
	imports.Proc(def.ViewComponent_SetOnParentViewChanged).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnChildViewChanged(fn viewOnChildViewChanged) {
	imports.Proc(def.ViewComponent_SetOnChildViewChanged).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnWindowChanged(fn viewOnWindowChanged) {
	imports.Proc(def.ViewComponent_SetOnWindowChanged).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnLayoutChanged(fn viewOnLayoutChanged) {
	imports.Proc(def.ViewComponent_SetOnLayoutChanged).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnFocus(fn viewOnFocus) {
	imports.Proc(def.ViewComponent_SetOnFocus).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnBlur(fn viewOnBlur) {
	imports.Proc(def.ViewComponent_SetOnBlur).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

func (m *TCEFViewComponent) SetOnThemeChanged(fn viewOnThemeChanged) {
	imports.Proc(def.ViewComponent_SetOnThemeChanged).Call(m.Instance(), api.MakeEventDataPtr(fn))
}
