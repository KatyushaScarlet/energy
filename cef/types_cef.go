//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

// cef -> energy 所有结构类型定义
// 每个结构对象创建 XXXRef.New() 创建并返回CEF对象, 创建后的对象是
// 引用CEF指针在不使用时,使用Free函数合理的释放掉该对象

package cef

import (
	"github.com/energye/energy/v2/common"
	"github.com/energye/energy/v2/consts"
	. "github.com/energye/energy/v2/types"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/api"
	"github.com/energye/golcl/lcl/types"
	"time"
	"unsafe"
)

// ICefCookie CEF Cookie
type ICefCookie struct {
	Url, Name, Value, Domain, Path string
	Secure, Httponly, HasExpires   bool
	Creation, LastAccess, Expires  time.Time
	Count, Total, ID               int32
	SameSite                       consts.TCefCookieSameSite
	Priority                       consts.TCefCookiePriority
	SetImmediately                 bool
	DeleteCookie                   bool
	Result                         bool
}

// TCefKeyEvent CEF 键盘事件
type TCefKeyEvent struct {
	Kind                 consts.TCefKeyEventType // called 'type' in the original CEF source code
	Modifiers            consts.TCefEventFlags
	WindowsKeyCode       Int32
	NativeKeyCode        Int32
	IsSystemKey          Int32
	Character            UInt16
	UnmodifiedCharacter  UInt16
	FocusOnEditableField Int32
}

// TCefRequestContextSettings CEF 请求上下文配置
type TCefRequestContextSettings struct {
	CachePath                        TCefString
	PersistSessionCookies            Int32
	AcceptLanguageList               TCefString // Remove CEF 118
	CookieableSchemesList            TCefString
	CookieableSchemesExcludeDefaults Int32
}

func (m *TCefRequestContextSettings) ToPtr() *tCefRequestContextSettingsPtr {
	return &tCefRequestContextSettingsPtr{
		CachePath:                        api.PascalStr(string(m.CachePath)),
		PersistSessionCookies:            uintptr(m.PersistSessionCookies),
		AcceptLanguageList:               api.PascalStr(string(m.AcceptLanguageList)), // Remove CEF 118
		CookieableSchemesList:            api.PascalStr(string(m.CookieableSchemesList)),
		CookieableSchemesExcludeDefaults: uintptr(m.CookieableSchemesExcludeDefaults),
	}
}

// TCefBrowserSettings CEF Browser配置
type TCefBrowserSettings struct {
	instance                   *tCefBrowserSettingsPtr
	WindowlessFrameRate        Integer
	StandardFontFamily         TCefString
	FixedFontFamily            TCefString
	SerifFontFamily            TCefString
	SansSerifFontFamily        TCefString
	CursiveFontFamily          TCefString
	FantasyFontFamily          TCefString
	DefaultFontSize            Integer
	DefaultFixedFontSize       Integer
	MinimumFontSize            Integer
	MinimumLogicalFontSize     Integer
	DefaultEncoding            TCefString
	RemoteFonts                consts.TCefState
	Javascript                 consts.TCefState
	JavascriptCloseWindows     consts.TCefState
	JavascriptAccessClipboard  consts.TCefState
	JavascriptDomPaste         consts.TCefState
	ImageLoading               consts.TCefState
	ImageShrinkStandaLonetoFit consts.TCefState
	TextAreaResize             consts.TCefState
	TabToLinks                 consts.TCefState
	LocalStorage               consts.TCefState
	Databases                  consts.TCefState
	Webgl                      consts.TCefState
	BackgroundColor            TCefColor
	AcceptLanguageList         TCefString // Remove CEF 118
	ChromeStatusBubble         consts.TCefState
	ChromeZoomBubble           consts.TCefState // Use[118]
}

// TCefCommandLine 进程启动命令行参数设置
type TCefCommandLine struct {
	commandLines map[string]string
}

// ICefCommandLine
type ICefCommandLine struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefProxy 代理配置
type TCefProxy struct {
	ProxyType              consts.TCefProxyType
	ProxyScheme            consts.TCefProxyScheme
	ProxyServer            string
	ProxyPort              int32
	ProxyUsername          string
	ProxyPassword          string
	ProxyScriptURL         string
	ProxyByPassList        string
	MaxConnectionsPerProxy int32
}

// TCefTouchEvent 触摸事件
type TCefTouchEvent struct {
	Id            int32
	X             float32
	Y             float32
	RadiusX       float32
	RadiusY       float32
	RotationAngle float32
	Pressure      float32
	Type          consts.TCefTouchEeventType
	Modifiers     consts.TCefEventFlags
	PointerType   consts.TCefPointerType
}

// TCustomHeader 自定义请求头
type TCustomHeader struct {
	CustomHeaderName  string
	CustomHeaderValue string
}

// TCefMouseEvent 鼠标事件
type TCefMouseEvent struct {
	X         int32
	Y         int32
	Modifiers consts.TCefEventFlags
}

// BeforePopupInfo 弹出子窗口信息
type BeforePopupInfo struct {
	TargetUrl         string
	TargetFrameName   string
	TargetDisposition consts.TCefWindowOpenDisposition
	UserGesture       bool
}

// TCefRect
//
//	/include/internal/cef_types_geometry.h (cef_rect_t)
type TCefRect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

type TCefRectArray struct {
	ptr    uintptr
	sizeOf uintptr
	count  uint32
}

type TRGBQuad struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

// NewTCefRectArray
//
//	TCefRect 动态数组结构, 通过指针引用取值
func NewTCefRectArray(ptr uintptr, count uint32) *TCefRectArray {
	return &TCefRectArray{
		ptr:    ptr,
		sizeOf: unsafe.Sizeof(TCefRect{}),
		count:  count,
	}
}

func (m *TCefRectArray) Count() int {
	return int(m.count)
}

func (m *TCefRectArray) Get(index int) *TCefRect {
	if m.count == 0 || index < 0 || index >= int(m.count) {
		return nil
	}
	return (*TCefRect)(common.GetParamPtr(m.ptr, index*int(m.sizeOf)))
}

// TCefSize
//
//	/include/internal/cef_types_geometry.h (cef_size_t)
type TCefSize struct {
	Width  int32
	Height int32
}

// TCefPoint
//
//	/include/internal/cef_types_geometry.h (cef_point_t)
type TCefPoint struct {
	X int32
	Y int32
}

// TCefCursorInfo
//
//	/include/internal/cef_types.h (cef_cursor_info_t)
type TCefCursorInfo struct {
	Hotspot          TCefPoint
	ImageScaleFactor Single
	Buffer           uintptr
	Size             TCefSize
}

// TCefSchemeRegistrarRef
type TCefSchemeRegistrarRef struct {
	instance unsafe.Pointer
}

// TCefBaseRefCounted
type TCefBaseRefCounted struct {
	instance unsafe.Pointer
}

// ICefRequestContext
type ICefRequestContext struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefCookieManager
type ICefCookieManager struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefCookieVisitor
type ICefCookieVisitor struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefBrowser
type ICefBrowser struct {
	base           TCefBaseRefCounted
	instance       unsafe.Pointer
	mainFrame      *ICefFrame
	requestContext *ICefRequestContext
	windowHandle   types.HWND
	idFrames       map[string]*ICefFrame
	nameFrames     map[string]*ICefFrame
}

// ICefFrame
// Html <iframe></iframe>
type ICefFrame struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefImage
type ICefImage struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefDraggableRegions 拖拽区域集合
type TCefDraggableRegions struct {
	regions      []TCefDraggableRegion
	regionsCount int
}

// TCefDraggableRegion 拖拽区域集
type TCefDraggableRegion struct {
	Bounds    TCefRect
	Draggable bool
}

// ICefProcessMessage
type ICefProcessMessage struct {
	base         TCefBaseRefCounted
	instance     unsafe.Pointer
	argumentList *ICefListValue
	name         string
}

// TCefBinaryValueArray
//
//	[]ICefBinaryValue
type TCefBinaryValueArray struct {
	instance     unsafe.Pointer
	binaryValues []*ICefBinaryValue
	count        uint32
}

// ICefBinaryValue -> ArgumentList
type ICefBinaryValue struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefValue -> ArgumentList
type ICefValue struct {
	base            TCefBaseRefCounted
	instance        unsafe.Pointer
	binaryValue     *ICefBinaryValue
	dictionaryValue *ICefDictionaryValue
	listValue       *ICefListValue
}

// ICefListValue -> ArgumentList
type ICefListValue struct {
	base             TCefBaseRefCounted
	instance         unsafe.Pointer
	values           map[int]*ICefValue
	binaryValues     map[int]*ICefBinaryValue
	dictionaryValues map[int]*ICefDictionaryValue
	listValues       map[int]*ICefListValue
}

// ICefDictionaryValue -> ArgumentList
type ICefDictionaryValue struct {
	base             TCefBaseRefCounted
	instance         unsafe.Pointer
	values           map[string]*ICefValue
	binaryValues     map[string]*ICefBinaryValue
	dictionaryValues map[string]*ICefDictionaryValue
	listValues       map[string]*ICefListValue
}

// ICefDisplayArray
//
//	[]ICefDisplayArray
type ICefDisplayArray struct {
	instance     unsafe.Pointer
	binaryValues []*ICefDisplayArray
	count        uint32
}

// ICefDisplay
type ICefDisplay struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefExtensionHandler
type ICefExtensionHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCustomExtensionHandler
type TCustomExtensionHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefExtension
type ICefExtension struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefSchemeHandlerFactory
//
//	/include/capi/cef_scheme_capi.h (cef_scheme_handler_factory_t)
type ICefSchemeHandlerFactory struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefResourceHandlerClass
type TCefResourceHandlerClass uintptr

// ICefRequest
type ICefRequest struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefResponse
type ICefResponse struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDomVisitor
type ICefDomVisitor struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDomDocument
type ICefDomDocument struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefScreenInfo
//
//	/include/internal/cef_types.h (cef_screen_info_t)
type TCefScreenInfo struct {
	DeviceScaleFactor Single
	Depth             int32
	DepthPerComponent int32
	IsMonochrome      int32
	Rect              TCefRect
	AvailableRect     TCefRect
}

// TCefTouchHandleState
//
//	/include/internal/cef_types.h (cef_touch_handle_state_t)
type TCefTouchHandleState struct {
	TouchHandleId    int32
	Flags            uint32
	Enabled          int32
	Orientation      consts.TCefHorizontalAlignment
	MirrorVertical   int32
	MirrorHorizontal int32
	Origin           TCefPoint
	Alpha            float32
}

// ICefRequestContextHandler
type ICefRequestContextHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefMenuModel 菜单
type ICefMenuModel struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	CefMis   *keyEventAccelerator
}

// ICefMenuModelDelegate
// /include/capi/cef_menu_model_delegate_capi.h (cef_menu_model_delegate_t)
type ICefMenuModelDelegate struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefStringMultiMap 实例
type ICefStringMultiMap struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPostData
type ICefPostData struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPostDataElement
type ICefPostDataElement struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefPostDataElementArray
type TCefPostDataElementArray struct {
	instance              unsafe.Pointer
	postDataElement       uintptr
	postDataElementLength uint32
}

// ICefBrowserView
// /include/capi/views/cef_browser_view_capi.h (cef_browser_view_t)
type ICefBrowserView struct {
	*ICefView
}

// ICefButton
// /include/capi/views/cef_button_capi.h (cef_button_t)
type ICefButton struct {
	*ICefView
}

// ICefLabelButton
// /include/capi/views/cef_label_button_capi.h (cef_label_button_t)
type ICefLabelButton struct {
	*ICefButton
}

// ICefMenuButton
// /include/capi/views/cef_menu_button_capi.h (cef_menu_button_t)
type ICefMenuButton struct {
	*ICefLabelButton
}

// ICefPanel
// /include/capi/views/cef_panel_capi.h (cef_panel_t)
type ICefPanel struct {
	*ICefView
}

// ICefScrollView
// /include/capi/views/cef_scroll_view_capi.h (cef_scroll_view_t)
type ICefScrollView struct {
	*ICefView
}

// ICefTextfield
// /include/capi/views/cef_textfield_capi.h (cef_textfield_t)
type ICefTextfield struct {
	*ICefView
}

// ICefWindow
// /include/capi/views/cef_window_capi.h (cef_window_t)
type ICefWindow struct {
	*ICefPanel
}

/*
*********************************
************* Views *************
*********************************

(*) Has CEF creation function
(d) Has delegate

----------------          ----------------------
| TCefView (d) | -------> | TCefTextfield (*d) |
----------------    |     ----------------------
					|
					|     ----------------------
					|---> | TCefScrollView (*) |
					|     ----------------------
					|
					|     ------------------          -------------------
					|---> | TCefPanel (*d) | -------> | TCefWindow (*d) |
					|     ------------------          -------------------
					|
					|     ------------------------
					|---> | TCefBrowserView (*d) |
					|     ------------------------
					|
					|     ------------------          -----------------------          -----------------------
					|---> | TCefButton (d) | -------> | TCefLabelButton (*) | -------> | TCefMenuButton (*d) |
						  ------------------          -----------------------          -----------------------


--------------          -----------------
| TCefLayout | -------> | TCefBoxLayout |
--------------    |     -----------------
				  |
				  |     ------------------
				  |---> | TCefFillLayout |
				  		------------------
*/

// ICefView
// /include/capi/views/cef_view_capi.h (cef_view_t)
type ICefView struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefViewDelegate
// /include/capi/views/cef_view_delegate_capi.h (cef_view_delegate_t)
type ICefViewDelegate struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefOverlayController TODO 未实现
// /include/capi/views/cef_overlay_controller_capi.h (cef_overlay_controller_t)
type ICefOverlayController struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefLayout
// /include/capi/views/cef_layout_capi.h (cef_layout_t)
type ICefLayout struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefFillLayout
// /include/capi/views/cef_fill_layout_capi.h (cef_fill_layout_t)
type ICefFillLayout struct {
	*ICefLayout
}

// ICefBoxLayout
// /include/capi/views/cef_box_layout_capi.h (cef_box_layout_t)
type ICefBoxLayout struct {
	*ICefLayout
}

// ICefBrowserViewDelegate
// /include/capi/views/cef_browser_view_delegate_capi.h (cef_browser_view_delegate_t)
type ICefBrowserViewDelegate struct {
	*ICefViewDelegate
}

// ICefButtonDelegate
// /include/capi/views/cef_button_delegate_capi.h (cef_button_delegate_t)
type ICefButtonDelegate struct {
	*ICefViewDelegate
}

// ICefMenuButtonDelegate
// /include/capi/views/cef_menu_button_delegate_capi.h (cef_menu_button_delegate_t)
type ICefMenuButtonDelegate struct {
	*ICefButtonDelegate
}

// ICefPanelDelegate
// /include/capi/views/cef_panel_delegate_capi.h (cef_panel_delegate_t)
type ICefPanelDelegate struct {
	*ICefViewDelegate
}

// ICefWindowDelegate
// /include/capi/views/cef_window_delegate_capi.h (cef_window_delegate_t)
type ICefWindowDelegate struct {
	*ICefPanelDelegate
}

// ICefTextFieldDelegate
// /include/capi/views/cef_textfield_delegate_capi.h (cef_textfield_delegate_t)
type ICefTextFieldDelegate struct {
	*ICefViewDelegate
}

// TCEFViewComponent
type TCEFViewComponent struct {
	lcl.IComponent
	instance unsafe.Pointer
}

// TCEFScrollViewComponent
type TCEFScrollViewComponent struct {
	*TCEFViewComponent
}

// TCEFBrowserViewComponent
type TCEFBrowserViewComponent struct {
	*TCEFViewComponent
}

type TCEFPanelComponent struct {
	*TCEFViewComponent
}

// TCEFWindowComponent 窗口组件
type TCEFWindowComponent struct {
	*TCEFPanelComponent
}

// TCEFButtonComponent
type TCEFButtonComponent struct {
	*TCEFViewComponent
}

// TCEFTextFieldComponent
type TCEFTextFieldComponent struct {
	*TCEFViewComponent
}

// TCEFLabelButtonComponent
type TCEFLabelButtonComponent struct {
	*TCEFButtonComponent
}

// TCEFMenuButtonComponent
type TCEFMenuButtonComponent struct {
	*TCEFLabelButtonComponent
}

// ICefMenuButtonPressedLock
// /include/capi/views/cef_menu_button_delegate_capi.h (cef_menu_button_pressed_lock_t)
type ICefMenuButtonPressedLock struct {
	base TCefBaseRefCounted
}

// TCefX509CertificateArray
// []ICefX509Certificate
type TCefX509CertificateArray struct {
	instance     unsafe.Pointer
	certificates []*ICefX509Certificate
	count        uint32
}

// ICefX509Certificate
//
//	/include/capi/cef_x509_certificate_capi.h (cef_x509certificate_t)
type ICefX509Certificate struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefX509CertPrincipal
// /include/capi/cef_x509_certificate_capi.h (cef_x509cert_principal_t)
type ICefX509CertPrincipal struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefSslInfo
//
//	/include/capi/cef_ssl_info_capi.h (cef_sslinfo_t)
type ICefSslInfo struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefClient
type ICefClient struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefAudioHandler
type ICefAudioHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefCommandHandler
type ICefCommandHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefContextMenuHandler
type ICefContextMenuHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDialogHandler
type ICefDialogHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDisplayHandler
type ICefDisplayHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDownloadHandler
type ICefDownloadHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDownloadItem 下载项
type ICefDownloadItem struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDragHandler
type ICefDragHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefFindHandler
type ICefFindHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefFocusHandler
type ICefFocusHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefFrameHandler
type ICefFrameHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPermissionHandler
type ICefPermissionHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefJsDialogHandler
type ICefJsDialogHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefKeyboardHandler
type ICefKeyboardHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefLifeSpanHandler
type ICefLifeSpanHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefLoadHandler
type ICefLoadHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPrintHandler
type ICefPrintHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefRenderHandler
type ICefRenderHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefRequestHandler
type ICefRequestHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefAccessibilityHandler
//
//	/include/capi/cef_accessibility_handler_capi.h (cef_accessibility_handler_t)
type ICefAccessibilityHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefResourceRequestHandler
//
//	/include/capi/cef_resource_request_handler_capi.h (cef_resource_request_handler_t)
type ICefResourceRequestHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefCookieAccessFilter
//
//	/include/capi/cef_resource_request_handler_capi.h (cef_cookie_access_filter_t)
type ICefCookieAccessFilter struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefResourceHandler
//
//	/include/capi/cef_resource_handler_capi.h (cef_resource_handler_t)
type ICefResourceHandler struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefResponseFilter
//
//	/include/capi/cef_response_filter_capi.h (cef_response_filter_t)
type ICefResponseFilter struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDragData
type ICefDragData struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefV8Exception
type ICefV8Exception struct {
	instance unsafe.Pointer
}

// ICefStreamWriter
//
//	/include/capi/cef_stream_capi.h (cef_stream_writer_t)
type ICefStreamWriter struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefV8Context
//
// v8上下文对象
//
// 生命周期
//  1. 在回调函数中有效
//  2. 回调函数外使用 cef.V8ContextRef.Current() 获取上下文对象
type ICefV8Context struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	browser  *ICefBrowser
	frame    *ICefFrame
	global   *ICefV8Value
}

// ICefV8Value
//
//	CEF V8 值类型, 对应到 JavaScrip 的类型, 使用该对象时需要合理的管理释放
type ICefV8Value struct {
	base         TCefBaseRefCounted
	instance     unsafe.Pointer
	valueType    consts.V8ValueType      // 值类型
	valueByIndex []*ICefV8Value          // 当前对象的所有数组集合
	valueByKeys  map[string]*ICefV8Value // 当前对象的所有key=value子集合
	cantNotFree  bool                    // 是否允许释放, false时允许释放
}

// ICefV8ValueKeys
type ICefV8ValueKeys struct {
	keys     *lcl.TStrings
	count    int
	keyArray []string
}

// TCefV8ValueArray ICefV8Value 数组的替代结构
type TCefV8ValueArray struct {
	instance         unsafe.Pointer
	arguments        uintptr
	argumentsLength  int
	argumentsCollect []*ICefV8Value
}

// ICefV8Handler
type ICefV8Handler struct {
	instance unsafe.Pointer
}

// ICefV8Interceptor
type ICefV8Interceptor struct {
	instance unsafe.Pointer
}

// ICefV8Accessor
type ICefV8Accessor struct {
	instance unsafe.Pointer
}

// ICefStreamReader
type ICefStreamReader struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPrintSettings
type ICefPrintSettings struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefSelectClientCertificateCallback
//
//	/include/capi/cef_request_handler_capi.h (cef_select_client_certificate_callback_t)
type ICefSelectClientCertificateCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefResourceReadCallback
//
//	/include/capi/cef_resource_handler_capi.h (cef_resource_read_callback_t)
type ICefResourceReadCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefResourceSkipCallback
//
//	/include/capi/cef_resource_handler_capi.h (cef_resource_skip_callback_t)
type ICefResourceSkipCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDeleteCookiesCallback
type ICefDeleteCookiesCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefSetCookieCallback
type ICefSetCookieCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
	ct       consts.CefCreateType
}

// ICefPrintDialogCallback
type ICefPrintDialogCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPrintJobCallback
type ICefPrintJobCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefCallback
type ICefCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefAuthCallback 授权回调
type ICefAuthCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefV8ArrayBufferReleaseCallback
type ICefV8ArrayBufferReleaseCallback struct {
	instance unsafe.Pointer
}

// ICefGetExtensionResourceCallback
type ICefGetExtensionResourceCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDownloadItemCallback
//
// 下载中回调
type ICefDownloadItemCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefBeforeDownloadCallback
//
// 下载之前回调
type ICefBeforeDownloadCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPdfPrintCallback
type ICefPdfPrintCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefCompletionCallback
type ICefCompletionCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefMediaRouter
// TODO no impl
type ICefMediaRouter struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefRunContextMenuCallback
//
//	/include/capi/cef_context_menu_handler_capi.h (cef_run_context_menu_callback_t)
type ICefRunContextMenuCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefRunQuickMenuCallback
//
//	/include/capi/cef_context_menu_handler_capi.h (cef_run_quick_menu_callback_t)
type ICefRunQuickMenuCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefFileDialogCallback
//
//	/include/capi/cef_dialog_handler_capi.h (cef_file_dialog_callback_t)
type ICefFileDialogCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefRunFileDialogCallback
// /include/capi/cef_browser_capi.h (cef_run_file_dialog_callback_t)
type ICefRunFileDialogCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefDownloadImageCallback
// /include/capi/cef_browser_capi.h (cef_download_image_callback_t)
type ICefDownloadImageCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefMediaAccessCallback
//
//	This interface is declared twice with almost identical parameters. "allowed_permissions" is defined as int and uint32.
//	/include/capi/cef_media_access_handler_capi.h (cef_media_access_callback_t)
//	/include/capi/cef_permission_handler_capi.h (cef_media_access_callback_t)
type ICefMediaAccessCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefPermissionPromptCallback
//
//	/include/capi/cef_permission_handler_capi.h (cef_permission_prompt_callback_t)
type ICefPermissionPromptCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefJsDialogCallback
//
//	/include/capi/cef_jsdialog_handler_capi.h (cef_jsdialog_callback_t)
type ICefJsDialogCallback struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefV8StackTrace
type ICefV8StackTrace struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefContextMenuParams 菜单显示时参数，当前鼠标右键的frame & html元素参数
type ICefContextMenuParams struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefDomNode
type ICefDomNode struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefMediaRouteArray
//
//	[]ICefMediaRoute
type TCefMediaRouteArray struct {
	instance   unsafe.Pointer
	mediaRoute []*ICefMediaRoute
	count      uint32
}

// ICefMediaRoute
//
//	/include/capi/cef_media_router_capi.h (cef_media_observer_t)
type ICefMediaRoute struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefMediaSinkArray
//
//	[]ICefMediaSink
type TCefMediaSinkArray struct {
	instance  unsafe.Pointer
	mediaSink []*ICefMediaSink
	count     uint32
}

// ICefMediaSink
//
//	/include/capi/cef_media_router_capi.h (cef_media_sink_t)
type ICefMediaSink struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// ICefNavigationEntry
//
//	/include/capi/cef_navigation_entry_capi.h (cef_navigation_entry_t)
type ICefNavigationEntry struct {
	base     TCefBaseRefCounted
	instance unsafe.Pointer
}

// TCefPreferenceRegistrarRef
// Class that manages custom preference registrations.
// <para><see href="https://bitbucket.org/chromiumembedded/cef/src/master/include/capi/cef_preference_capi.h">CEF source file: /include/capi/cef_preference_capi.h (cef_preference_registrar_t)</see></para>
type TCefPreferenceRegistrarRef struct {
	instance unsafe.Pointer
}

// TCefRange
//
//	/include/internal/cef_types_geometry.h (cef_range_t)
type TCefRange struct {
	From int32
	To   int32
}

// /include/internal/cef_types_geometry.h (cef_insets_t)
type TCefInsets struct {
	Top    int32
	Left   int32
	Bottom int32
	Right  int32
}

// TCefAudioParameters
// /include/internal/cef_types.h (cef_audio_parameters_t)
type TCefAudioParameters struct {
	channelLayout   consts.TCefChannelLayout
	sampleRate      int32
	framesPerBuffer int32
}

type CefPdfPrintSettings struct {
	Landscape           int32                         // Integer
	PrintBackground     int32                         // Integer
	Scale               float64                       // double
	PaperWidth          float64                       // double
	PaperHeight         float64                       // double
	PreferCssPageSize   int32                         // Integer
	MarginType          consts.TCefPdfPrintMarginType // TCefPdfPrintMarginType
	MarginTop           float64                       // double
	MarginRight         float64                       // double
	MarginBottom        float64                       // double
	MarginLeft          float64                       // double
	PageRanges          string                        // TCefString
	DisplayHeaderFooter int32                         // Integer
	HeaderTemplate      string                        // TCefString
	FooterTemplate      string                        // TCefString
}

// /include/internal/cef_types.h (cef_popup_features_t)
type TCefPopupFeatures struct {
	X                  Integer
	XSet               Integer
	Y                  Integer
	YSet               Integer
	Width              Integer
	WidthSet           Integer
	Height             Integer
	HeightSet          Integer
	MenuBarVisible     Integer // Use-CEF:[49]
	StatusBarVisible   Integer // Use-CEF:[49]
	ToolBarVisible     Integer // Use-CEF:[49]
	LocationBarVisible Integer
	ScrollbarsVisible  Integer // Use-CEF:[49]
	IsPopup            Integer // CEF 110 ~ Current :True (1) if browser interface elements should be hidden.
	Resizable          Integer
	Fullscreen         Integer
	Dialog             Integer
	AdditionalFeatures TCefStringList // Use-CEF:[49]
}

// /include/internal/cef_types.h (cef_composition_underline_t)
type TCefCompositionUnderline struct {
	Range           TCefRange
	Color           TCefColor
	BackgroundColor TCefColor
	Thick           int32
	Style           consts.TCefCompositionUnderlineStyle
}

type TCefCompositionUnderlineArray struct {
	count  int
	ptr    uintptr
	sizeOf uintptr
}

// /include/internal/cef_types.h (cef_box_layout_settings_t)
type TCefBoxLayoutSettings struct {
	Horizontal                    Integer
	InsideBorderHorizontalSpacing Integer
	InsideBorderVerticalSpacing   Integer
	InsideBorderInsets            TCefInsets
	BetweenChildSpacing           Integer
	MainAxisAlignment             consts.TCefMainAxisAlignment
	CrossAxisAlignment            consts.TCefCrossAxisAlignment
	MinimumCrossAxisSize          Integer
	DefaultFlex                   Integer
}

// ResultString 字符串返回值
type ResultString struct {
	value string
}

type TChromiumOptions struct {
	Chromium                   IChromium
	javascript                 consts.TCefState
	javascriptCloseWindows     consts.TCefState
	javascriptAccessClipboard  consts.TCefState
	javascriptDomPaste         consts.TCefState
	imageLoading               consts.TCefState
	imageShrinkStandaloneToFit consts.TCefState
	textAreaResize             consts.TCefState
	tabToLinks                 consts.TCefState
	localStorage               consts.TCefState
	databases                  consts.TCefState
	webgl                      consts.TCefState
	backgroundColor            TCefColor
	acceptLanguageList         String // Remove CEF 118
	windowlessFrameRate        Integer
	chromeStatusBubble         consts.TCefState
}

func (m *TCefCompositionUnderlineArray) Count() int {
	return m.count
}

func (m *TCefCompositionUnderlineArray) Get(index int) *TCefCompositionUnderline {
	if index >= 0 && index < m.count {
		return (*TCefCompositionUnderline)(common.GetParamPtr(m.ptr, index*int(m.sizeOf)))
	}
	return nil
}

func (m *ResultString) SetValue(value string) {
	m.value = value
}

func (m *ResultString) Value() string {
	return m.value
}

// ResultBool  bool返回值
type ResultBool struct {
	value bool
}

func (m *ResultBool) SetValue(value bool) {
	m.value = value
}

func (m *ResultBool) Value() bool {
	return m.value
}

// ResultBytes  []byte返回值
type ResultBytes struct {
	value []byte
}

func (m *ResultBytes) SetValue(value []byte) {
	m.value = value
}

func (m *ResultBytes) Value() []byte {
	return m.value
}

// NewCefRect
func NewCefRect(x, y, width, height int32) *TCefRect {
	return &TCefRect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// NewCefSize
func NewCefSize(width, height int32) *TCefSize {
	return &TCefSize{
		Width:  width,
		Height: height,
	}
}

// NewCefPoint
func NewCefPoint(x, y int32) *TCefPoint {
	return &TCefPoint{
		X: x,
		Y: y,
	}
}

func (m *TCefKeyEvent) KeyDown() bool {
	return m.Kind == consts.KEYEVENT_RAW_KEYDOWN || m.Kind == consts.KEYEVENT_KEYDOWN
}

func (m *TCefKeyEvent) KeyUp() bool {
	return m.Kind == consts.KEYEVENT_KEYUP
}

// TLinuxWindowProperties String version
// <para><see href="https://bitbucket.org/chromiumembedded/cef/src/master/include/internal/cef_types.h">CEF source file: /include/internal/cef_types.h (cef_linux_window_properties_t)</see></para>
type TLinuxWindowProperties struct {
	/// Main window's Wayland's app_id
	WaylandAppId string
	/// Main window's WM_CLASS_CLASS in X11
	WmClassClass string
	/// Main window's WM_CLASS_NAME in X11
	WmClassName string
	/// Main window's WM_WINDOW_ROLE in X11
	WmRoleName string
}

func (m *TLinuxWindowProperties) ToPtr() *tLinuxWindowPropertiesPtr {
	return &tLinuxWindowPropertiesPtr{
		WaylandAppId: api.PascalStr(m.WaylandAppId),
		WmClassClass: api.PascalStr(m.WmClassClass),
		WmClassName:  api.PascalStr(m.WmClassName),
		WmRoleName:   api.PascalStr(m.WmRoleName),
	}
}
