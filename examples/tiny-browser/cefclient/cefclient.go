package main

import (
	"embed"
	"fmt"
	"github.com/energye/energy/v2/cef"
	"github.com/energye/energy/v2/cef/process"
	"github.com/energye/energy/v2/consts"
	. "github.com/energye/energy/v2/examples/tiny-browser/cefclient/browse"
	"github.com/energye/energy/v2/examples/tiny-browser/cefclient/views_style"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/api"
	"os"
	"path/filepath"
)

var (
	app    *cef.TCEFApplication
	vf     *ViewsFramework
	width  int32 = 800
	height int32 = 600
)

type ViewsFramework struct {
	*lcl.TComponent
	chromium    cef.IChromium
	window      *cef.TCEFWindowComponent
	browserView *cef.TCEFBrowserViewComponent
	homePage    string
}

//go:embed assets
var assets embed.FS

func main() {
	Assets = assets
	cef.GlobalInit(nil, nil)
	rootCache := filepath.Join(consts.CurrentExecuteDir, "rootcache")
	app = cef.CreateApplication()
	app.SetFrameworkDirPath(os.Getenv("ENERGY_HOME"))
	app.SetMultiThreadedMessageLoop(false)
	app.SetExternalMessagePump(false)
	app.SetRootCache(rootCache)
	app.SetCache(filepath.Join(rootCache, "cache"))
	app.SetLocale(consts.LANGUAGE_zh_CN)
	app.SetTouchEvents(consts.STATE_ENABLED)
	app.SetDisableZygote(true)
	//fmt.Println("libname.LibName:", libname.LibName)
	fmt.Println("WidgetUI:", api.WidgetUI(), "ChromeVersion:", app.ChromeVersion(), "LibCefVersion:", app.LibCefVersion())

	kPrefWindowRestore := "cefclient.window_restore"
	app.SetOnRegisterCustomPreferences(func(type_ consts.TCefPreferencesType, registrar *cef.TCefPreferenceRegistrarRef) {
		fmt.Println("OnRegisterCustomPreferences ProcessType:", process.Args.ProcessType())
		if type_ == consts.CEF_PREFERENCES_TYPE_GLOBAL {
			dict := cef.DictionaryValueRef.New()
			dict.SetInt(kPrefWindowRestore, int32(consts.CEF_SHOW_STATE_NORMAL))
			value := cef.ValueRef.New()
			value.SetDictionary(dict)
			registrar.AddPreference(kPrefWindowRestore, value)
		}
	})
	app.SetOnAlreadyRunningAppRelaunch(func(commandLine *cef.ICefCommandLine, currentDirectory string) bool {
		fmt.Println("OnAlreadyRunningAppRelaunch ProcessType:", process.Args.ProcessType())
		// 在此处创建一个新窗口

		// 重新启动处理好了
		return true
	})
	app.SetOnContextInitialized(func() {
		fmt.Println("OnContextInitialized ProcessType:", process.Args.ProcessType())
		fmt.Println("  GetScreenDPI:", cef.GetScreenDPI(), "GetDeviceScaleFactor:", cef.GetDeviceScaleFactor())
	})
	if app.StartMainProcess() {
		fmt.Println("StartMainProcess Success")
		StartServer()
		// 创建窗口
		vf = new(ViewsFramework)
		vf.homePage = "http://localhost:22022"
		vf.Create()
		vf.window.CreateTopLevelWindow()
		app.RunMessageLoop()
	}
}

func (m *ViewsFramework) Create() {
	m.TComponent = lcl.NewComponent(nil)
	m.chromium = cef.NewChromium(m, nil)
	m.chromium.SetOnBeforeClose(func(sender lcl.IObject, browser *cef.ICefBrowser) {
		fmt.Println("OnBeforeClose")
		app.QuitMessageLoop()
	})
	var (
		minimumWindowSize = cef.TCefSize{Width: 100, Height: 100}
		canGoBack_        bool
		canGoForward_     bool
		isLoading_        bool
	)
	m.browserView = cef.BrowserViewComponentRef.New(m)
	m.browserView.SetID(ID_BROWSER_VIEW)
	//m.browserView.SetPreferAccelerators(true)
	m.window = cef.WindowComponentRef.New(m)
	m.window.SetID(ID_WINDOW)
	m.chromium.SetOnAutoResize(func(sender lcl.IObject, browser *cef.ICefBrowser, newSize *cef.TCefSize) bool {
		fmt.Println("OnAutoResize", newSize)
		return true
	})
	m.chromium.SetOnLoadingStateChange(func(sender lcl.IObject, browser *cef.ICefBrowser, isLoading, canGoBack, canGoForward bool) {
		isLoading_ = isLoading
		canGoBack_ = canGoBack
		canGoForward_ = canGoForward
		fmt.Println("OnLoadingStateChange:", isLoading_, canGoBack_, canGoForward_)
	})
	m.chromium.SetOnLoadEnd(func(sender lcl.IObject, browser *cef.ICefBrowser, frame *cef.ICefFrame, httpStatusCode int32) {
		fmt.Println("OnLoadEnd httpStatusCode:", httpStatusCode)
	})
	m.chromium.SetOnDraggableRegionsChanged(func(sender lcl.IObject, browser *cef.ICefBrowser, frame *cef.ICefFrame, regions *cef.TCefDraggableRegions) {
		fmt.Println("OnDraggableRegionsChanged RegionsCount:", regions.RegionsCount())
		for i := 0; i < regions.RegionsCount(); i++ {
			fmt.Println("Region:", regions.Region(i))
			regions.Region(i).Bounds.X = 0
			regions.Region(i).Bounds.Y = 0
			regions.Region(i).Bounds.Height = 30
		}
		m.window.SetDraggableRegions(regions.Regions())
	})
	m.window.SetOnGetMinimumSize(func(view *cef.ICefView, result *cef.TCefSize) {
		if view.GetID() == ID_WINDOW {
			//fmt.Println("OnGetMinimumSize", result)
			*result = minimumWindowSize
		}
	})
	m.window.SetOnGetPreferredSize(func(view *cef.ICefView, result *cef.TCefSize) {
		m.window.SetBackgroundColor(cef.CefColorSetARGB(255, 33, 34, 38))
	})
	m.window.SetOnGetTitleBarHeight(func(window *cef.ICefWindow, titleBarHeight *float32, result *bool) {
		fmt.Println("OnGetTitleBarHeight:", *titleBarHeight, *result)
	})
	m.window.SetOnWindowBoundsChanged(func(window *cef.ICefWindow, newBounds cef.TCefRect) {
		//fmt.Println("OnWindowBoundsChanged", newBounds)
		//m.window.SizeToPreferredSize()
		//m.browserView.SetSize(cef.TCefSize{Width: newBounds.Width - minimumWindowSize.Height, Height: newBounds.Height})
	})
	m.window.SetOnCanClose(func(window *cef.ICefWindow, result *bool) {
		fmt.Println("OnCanClose:", *result)
		*result = m.chromium.TryCloseBrowser()
	})
	m.window.SetOnIsFrameless(func(window *cef.ICefWindow, result *bool) {
		//*result = true
	})
	m.window.SetOnGetInitialBounds(func(window *cef.ICefWindow, result *cef.TCefRect) {
		fmt.Println("OnGetInitialBounds")
		result.Width = width
		result.Height = height
	})
	m.window.SetOnWindowActivationChanged(func(window *cef.ICefWindow, active bool) {
		fmt.Println("OnWindowActivationChanged")
	})
	m.window.SetOnGetInitialBounds(func(window *cef.ICefWindow, result *cef.TCefRect) {
		result.Width = 800
		result.Height = 600
		fmt.Println("OnGetInitialBounds")
	})
	m.window.SetOnThemeColorsChanged(func(window *cef.ICefWindow, chromeTheme int32) {
		views_style.WindowApplyTo(window)
	})
	m.window.SetOnAccelerator(func(window *cef.ICefWindow, commandId int32, result *bool) {
		fmt.Println("OnAccelerator commandId:", commandId)
		//表示已处理，否则还会执行多次
		m.chromium.CloseBrowser(true)
		*result = true
	})
	titleBar := NewTitleBar(m.window)
	m.window.SetOnWindowChanged(func(view *cef.ICefView, added bool) {
		fmt.Println("OnWindowChanged added:", added)
		if added {
			//titleBar.CreateTitleBrowser(m)
		}
	})
	//m.window.SetOnKeyEvent(func(window *cef.ICefWindow, event *cef.TCefKeyEvent, result *bool) {
	//	fmt.Println("OnKeyEvent")
	//})
	//m.window.SetOnWithStandardWindowButtons()
	m.window.SetOnWindowCreated(func(window *cef.ICefWindow) {
		fmt.Println("OnWindowCreated")
		m.window.SetID(ID_WINDOW)
		//window.SetThemeColor()
		m.window.ThemeChanged()
		m.window.SetWindowIcon(LoadImage("app-icon.png"))
		m.window.SetWindowAppIcon(LoadImage("app-icon.png"))
		m.window.SetTitle("Go ENERGY Client")

		menuBar := NewMenuBar(m.window) // 顶部菜单栏
		toolBar := NewToolBar(m.window) // 顶部工具栏
		if m.chromium.CreateBrowserByBrowserViewComponent(m.homePage, m.browserView, nil, nil) {
			//regions := []cef.TCefDraggableRegion{}
			//regions = append(regions, cef.TCefDraggableRegion{Bounds: cef.TCefRect{X: 130, Y: 0, Width: 100, Height: 30}})
			//m.window.SetDraggableRegions(regions)
			fmt.Println("ChromeToolbar:", m.browserView.ChromeToolbar().IsValid())
			// 允许|browser_view_|增长并填充任何剩余空间。
			windowLayout := m.window.SetToBoxLayout(cef.TCefBoxLayoutSettings{
				BetweenChildSpacing: 5,
				CrossAxisAlignment:  consts.CEF_AXIS_ALIGNMENT_STRETCH,
			})

			// 菜单栏, 创建菜单，并添加到菜单栏中
			menuBar.CreateFileMenuItems()
			menuBar.CreateTestMenuItems()

			// 工具栏, 创建工具组件，并添加到工具栏中
			toolBar.CreateToolComponent()

			//var minWidth int32 = toolBar.AllButtonWidth()
			//var minHeight int32 = toolBar.EnsureToolPanel().GetBounds().Height + 100
			//minimumWindowSize = cef.TCefSize{Width: minWidth, Height: minHeight}
			//fmt.Println("minWidth:", minWidth, "minHeight:", minHeight)

			// 标题栏添加到窗口
			m.window.AddChildView(titleBar.EnsureTitlePanel().AsView())

			// 菜单栏添加到窗口
			m.window.AddChildView(menuBar.EnsureMenuPanel().AsView())

			// 工具栏添加到窗口
			m.window.AddChildView(toolBar.EnsureToolPanel().AsView())

			// 浏览器view添加到窗口
			m.window.AddChildView(m.browserView.AsView())
			windowLayout.SetFlexForView(m.browserView.AsView(), 1)

			m.window.Layout()
			// 窗口居中
			display := m.window.Display()
			if display.IsValid() {
				workArea := display.WorkArea()
				position := &cef.TCefPoint{
					X: ((workArea.Width - width) / 2) + workArea.X,
					Y: ((workArea.Height - height) / 2) + workArea.Y,
				}
				m.window.SetBounds(cef.NewCefRect(position.X, position.Y, width, height))
			}
			m.window.Show()
			//m.window.ShowAsBrowserModalDialog(m.browserView.BrowserView())
			m.browserView.RequestFocus()
		}
	})
}
