package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	appx "gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	gwidget "gioui.org/widget"
	"gioui.org/widget/material"
	"nenki.ui/app"
	"nenki.ui/widget"
	"nenki.ui/widget/axis"
)

func printMemStat(ms runtime.MemStats) {
	runtime.ReadMemStats(&ms)
	fmt.Println("--------------------------------------")
	fmt.Println("Memory Statistics Reporting time: ", time.Now())
	fmt.Println("--------------------------------------")
	fmt.Println("Bytes of allocated heap objects: ", ms.Alloc)
	fmt.Println("Total bytes of Heap object: ", ms.TotalAlloc)
	fmt.Println("Bytes of memory obtained from OS: ", ms.Sys)
	fmt.Println("Count of heap objects: ", ms.Mallocs)
	fmt.Println("Count of heap objects freed: ", ms.Frees)
	fmt.Println("Count of live heap objects", ms.Mallocs-ms.Frees)
	fmt.Println("Number of completed GC cycles: ", ms.NumGC)
	fmt.Println("--------------------------------------")
}

var buttonPool = sync.Pool{
	New: func() interface{} {
		return new(widget.Button)
	},
}

func main() {
	var ms runtime.MemStats
	printMemStat(ms)
	//Test()
	app.NewApp("测试").Size(1024, 1024).Title("你好").DragFiles(true).
		Then(func(self *app.App, root *widget.ContainerLayout) {
			h := widget.NewRowLayout()
			root.AppendChild(h)
			v := widget.NewListLayout(axis.Vertical).ScrollToEnd(true)
			h.AppendFlexChild(0.5, widget.NewButton("float").FontSize(60).CornerRadius(10))

			h.AppendFlexChild(1, v)

			for i := 0; i < 50; i++ {
				v.AppendChild(widget.NewButton(fmt.Sprintf("item %d", i)).FontSize(20).CornerRadius(0))
			}

		})

	// 阻塞
	app.Run()
}

func Test() {
	go func() {
		w := appx.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	appx.Main()
}

func loop(w *appx.Window) error {
	var a gwidget.Bool = gwidget.Bool{}
	var b gwidget.Bool = gwidget.Bool{}
	th := material.NewTheme()

	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {

		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(material.CheckBox(th, &a, "测试1").Layout),
					layout.Rigid(material.CheckBox(th, &b, "测试2").Layout),
				)
			})
			e.Frame(gtx.Ops)
		}
	}
}
