package main

import (
	"embed"
	"fmt"
	"github.com/energye/energy/v2/cef"
	"github.com/energye/energy/v2/cef/ipc"
	"github.com/energye/energy/v2/cef/ipc/target"
	"github.com/energye/energy/v2/cef/ipc/types"
	"github.com/energye/energy/v2/consts"
	t "github.com/energye/energy/v2/types"
	lcltypes "github.com/energye/golcl/lcl/types"
)

//go:embed resources
var resources embed.FS

func main() {
	//全局初始化 每个应用都必须调用的
	cef.GlobalInit(nil, &resources)
	//创建应用
	var app = cef.NewApplication()
	cef.BrowserWindow.Config.Title = "Energy - mockevent"
	cef.BrowserWindow.Config.Url = "fs://energy" // 设置默认
	cef.BrowserWindow.Config.LocalResource(cef.LocalLoadConfig{
		ResRootDir: "resources",
		FS:         &resources, //静态资源所在的 embed.FS
	}.Build())

	/*
		模拟: 按钮点击, 文本输入, 滚动条
			1. 仅在主进程中使用ipc监听renderLoadEnd事件
			2. 在渲染进程中的app.SetOnRenderLoadEnd监听渲染进程页面加载结束事件，在这个事件里获取页面html元素位置和大小
			3. html元素位置和大小获取之后触发主进程监听事件, 把要获取到的html元素位置和大小发送到主进程renderLoadEnd事件
			4. 主进程renderLoadEnd事件被触发后使用渲染进程传递的元素数据模拟事件操作
				模拟事件使用当前窗口chromium或者browser提供的函数
				chromium.SendXXX 只能在主进程中使用
				browser.SendXXX  可在主进程或渲染进程中直接使用
			5. 最后主进程回复ipc消息给渲染进程
	*/

	// 在渲染进程中处理结束事件
	// 通过VisitDom获取html元素的位置和大小
	// SetOnVisit 函数只能在渲染进程中执行
	app.SetOnRenderLoadEnd(func(browser *cef.ICefBrowser, frame *cef.ICefFrame, httpStatusCode int32) {
		// 浏览器ID和通道ID, 用于回复给渲染进程使用
		var (
			browserId = browser.Identifier()
			channelId = frame.Identifier()
		)
		// 创建 dom visitor
		visitor := cef.DomVisitorRef.New()
		// 监听事件
		// 这个事件在渲染进程中才会执行
		visitor.SetOnVisit(func(document *cef.ICefDomDocument) {
			// html > body
			body := document.GetBody()
			// 使用元素id获取
			btn1 := body.GetDocument().GetElementById("btn1")
			btn2 := body.GetDocument().GetElementById("btn2")
			btn3 := body.GetDocument().GetElementById("btn3")
			inpText := body.GetDocument().GetElementById("inpText")
			inp2Text := body.GetDocument().GetElementById("inp2Text")
			// dom元素Rect集合数据发送到主进程
			var doms = make(map[string]cef.TCefRect)
			doms["btn1"] = btn1.GetElementBounds()
			doms["btn2"] = btn2.GetElementBounds()
			doms["btn3"] = btn3.GetElementBounds()
			doms["inpText"] = inpText.GetElementBounds()
			doms["inp2Text"] = inp2Text.GetElementBounds()
			// 触发主进程ipc监听事件
			ipc.EmitTarget("renderLoadEnd", target.NewTargetMain(), browserId, channelId, doms)
		})
		// 调用该函数后, 执行SetOnVisit回调函数
		frame.VisitDom(visitor)
	})

	// 丰富一下示例, 在Go中主子进程消息传递
	ipc.On("repayMockIsSuccess", func() {
		fmt.Println("mock success")
	}, // OtSub 仅子进程监听该事件
		types.OnOptions{OnType: types.OtSub})

	cef.BrowserWindow.SetBrowserInit(func(event *cef.BrowserEvent, window cef.IBrowserWindow) {
		var domXYCenter = func(bound cef.TCefRect) (int32, int32) {
			return bound.X + bound.Width/2, bound.Y + bound.Height/2
		}
		// 模拟按钮点击事件
		var buttonClickEvent = func(domRect cef.TCefRect) {
			window.RunOnMainThread(func() { //在UI主线程执行
				chromium := window.Chromium()
				// 鼠标事件
				me := &cef.TCefMouseEvent{}
				// 设置元素坐标，元素坐标相对于窗口，这里取元素中间位置
				me.X, me.Y = domXYCenter(domRect)
				// 模拟鼠标到指定位置
				chromium.SendMouseMoveEvent(me, false)
				// 模拟鼠标双击事件
				//   左键点击按下1次
				chromium.SendMouseClickEvent(me, consts.MBT_LEFT, false, 1)
				//   左键点击抬起1次
				chromium.SendMouseClickEvent(me, consts.MBT_LEFT, true, 1)
			})
		}
		// 模拟文本输入
		var inputTextEvent = func(value string, domRect cef.TCefRect) {
			// 文本框区域点击, 获取焦点
			buttonClickEvent(domRect)
			window.RunOnMainThread(func() { //在UI主线程执行
				chromium := window.Chromium()
				// 鼠标事件
				me := &cef.TCefMouseEvent{}
				// 设置元素坐标，元素坐标相对于窗口，这里取元素中间位置
				me.X, me.Y = domXYCenter(domRect)
				// 一个一个字符设置
				for _, v := range value {
					chromium.SendKeyEvent(keyPress(string(v)))
				}
			})
		}
		ipc.On("renderLoadEnd", func(browserId int32, channelId int64, doms map[string]cef.TCefRect) {
			fmt.Println("doms", doms)
			// 按钮
			buttonClickEvent(doms["btn1"])
			buttonClickEvent(doms["btn2"])
			buttonClickEvent(doms["btn3"])
			// 文本框
			inputTextEvent("我爱中国！", doms["inp2Text"])           //中文
			inputTextEvent("energy.yanghy.cn", doms["inpText"]) //英文
			// 滚动条
			// 回复到渲染进程执行成功, 触发是Go的事件.
			ipc.EmitTarget("repayMockIsSuccess", target.NewTarget(browserId, channelId, target.TgGoSub))
		})
	})
	//运行应用
	cef.Run(app)
}

func keyPress(key string) *cef.TCefKeyEvent {
	utf8Key := &lcltypes.TUTF8Char{}
	utf8Key.SetString(key)
	event := &cef.TCefKeyEvent{}
	var asciiCode int
	fmt.Sscanf(utf8Key.ToString(), "%c", &asciiCode)
	event.Kind = consts.KEYEVENT_CHAR
	event.WindowsKeyCode = t.Int32(asciiCode)
	event.FocusOnEditableField = 1 // 0=false, 1=true
	return event
}
