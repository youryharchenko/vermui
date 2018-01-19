package render

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/maruel/panicparse/stack"
	"github.com/pkg/errors"

	tm "github.com/verdverm/termbox-go"
)

var renderJobs chan []Bufferer
var renderLock sync.Mutex

var termWidth int
var termHeight int

func Init() error {
	if err := tm.Init(); err != nil {
		return err
	}
	tm.SetInputMode(tm.InputAlt | tm.InputMouse)

	termSync()

	renderJobs = make(chan []Bufferer, 3)
	return nil
}

func Start() error {
	return startRenderLoop()
}

// Close finalizes vermui library,
// should be called after successful initialization when vermui's functionality isn't required anymore.
func Stop() error {
	stopRenderLoop()
	tm.Close()
	return nil
}

func TermRect() image.Rectangle {
	return image.Rect(0, 0, termWidth, termHeight)
}

// TermWidth returns the current terminal's width.
func TermWidth() int {
	return termWidth
}

// TermHeight returns the current terminal's height.
func TermHeight() int {
	return termHeight
}

func Render(bs ...Bufferer) {
	//go func() { renderJobs <- bs }()
	renderJobs <- bs
}

func Clear() error {
	renderLock.Lock()
	defer renderLock.Unlock()

	err := tm.Clear(tm.ColorDefault, ToTmAttr(ThemeAttr("bg")))
	if err != nil {
		return errors.Wrapf(err, "in lib.Clear() - tm.Clear()\n")
	}

	return nil
}

func ClearArea(r image.Rectangle, bg Attribute) error {
	renderLock.Lock()
	defer renderLock.Unlock()

	clearArea(r, bg)

	err := tm.Flush()
	if err != nil {
		return errors.Wrapf(err, "in vermui.render() - tm.Flush()\n")
	}

	return nil
}

func startRenderLoop() error {
	go func() {
		for bs := range renderJobs {
			err := doRender(bs...)
			if err != nil {
				perr := errors.Wrapf(err, "in lib.Start() - render goroutine\n")
				panic(perr)
			}
		}
	}()

	return nil
}

func stopRenderLoop() error {
	close(renderJobs)

	return nil
}

func termSync() (err error) {
	renderLock.Lock()
	defer renderLock.Unlock()
	err = tm.Sync()
	if err != nil {
		return errors.Wrapf(err, "in lib.Clear()\n")
	}
	termWidth, termHeight = tm.Size()
	return nil
}

// Render renders all Bufferer in the given order from left to right,
// right could overlap on left ones.
func doRender(bs ...Bufferer) error {
	defer func() {
		if e := recover(); e != nil {
			// Stop VermUI
			Stop()

			// Print a formatted panic output
			fmt.Fprintf(os.Stderr, "Captured a panic(value=%v) when rendering Bufferer. Exit vermui and clean terminal...\nPrint stack trace:\n\n", e)
			//debug.PrintStack()
			gs, err := stack.ParseDump(bytes.NewReader(debug.Stack()), os.Stderr)
			if err != nil {
				debug.PrintStack()
				os.Exit(1)
			}
			p := &stack.Palette{}
			buckets := stack.SortBuckets(stack.Bucketize(gs, stack.AnyValue))
			srcLen, pkgLen := stack.CalcLengths(buckets, false)
			for _, bucket := range buckets {
				io.WriteString(os.Stdout, p.BucketHeader(&bucket, false, len(buckets) > 1))
				io.WriteString(os.Stdout, p.StackLines(&bucket.Signature, srcLen, pkgLen, false))
			}
			os.Exit(1)
		}
	}()

	renderLock.Lock()
	defer renderLock.Unlock()

	for _, b := range bs {

		buf := b.Buffer()
		// set cels in buf
		for p, c := range buf.CellMap {
			if p.In(buf.Area) {
				tm.SetCell(p.X, p.Y, c.Ch, ToTmAttr(c.Fg), ToTmAttr(c.Bg))
			}
		}
	}

	// render
	err := tm.Flush()
	if err != nil {
		return errors.Wrapf(err, "in vermui.render() - tm.Flush()\n")
	}
	termWidth, termHeight = tm.Size()

	return nil
}

func clearArea(r image.Rectangle, bg Attribute) {
	for i := r.Min.X; i < r.Max.X; i++ {
		for j := r.Min.Y; j < r.Max.Y; j++ {
			tm.SetCell(i, j, ' ', tm.ColorDefault, ToTmAttr(bg))
		}
	}
}
