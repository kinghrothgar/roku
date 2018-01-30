package main

import (
	"fmt"
	"log"
	"os"

	"github.com/huandu/xstrings"
	"github.com/jroimartin/gocui"
	"github.com/kinghrothgar/roku/roku"
)

var ip string
var rokuClient *roku.Roku

var instructWidth = 35

func castRune(s string) rune {
	return rune(s[0])
}

func rokuKeyPress(key string) error {
	return rokuClient.KeyPress(key)
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	if name == "insert" || name == "apps" {
		g.Cursor = true
	} else {
		g.Cursor = false
	}
	return g.SetViewOnTop(name)
}

func setRectangleView(g *gocui.Gui, name string, x int, y int) (*gocui.View, error) {
	maxX, maxY := g.Size()
	return g.SetView(name, maxX/2-x/2, maxY/2-y/2, maxX/2+x/2, maxY/2+y/2)
}

func output(g *gocui.Gui, v *gocui.View) error {
	aOut, err := g.View("action")
	if err != nil {
		return err
	}
	fmt.Fprint(aOut, "output")
	vOut, err := g.View("output")
	if err != nil {
		return err
	}
	fmt.Fprint(vOut, v.ViewBuffer())
	return nil
}

func selectViewHandler(viewStr string) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		_, err := setCurrentViewOnTop(g, viewStr)
		return err
	}
}

func sendInsert(g *gocui.Gui, v *gocui.View) error {
	// send text to roku
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	_, err := setCurrentViewOnTop(g, "remote")
	return err
}

func launchApp(g *gocui.Gui, v *gocui.View) error {
	// send text to roku
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	_, err := setCurrentViewOnTop(g, "remote")
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keyPressHandler(key string) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return rokuKeyPress(key)
	}
}

func setRemoteKeybindRoku(g *gocui.Gui, key interface{}, rokuKey string) {
	if err := g.SetKeybinding("remote", key, gocui.ModNone, keyPressHandler(rokuKey)); err != nil {
		log.Panicln(err)
	}
}

func centerInstruct(cmd string, desc string) string {
	shiftLeft := 10
	cmdPadLen := (instructWidth - shiftLeft) / 2
	descPadLen := (instructWidth + shiftLeft + 1) / 2
	cmdPadded := xstrings.RightJustify(cmd, cmdPadLen, " ")
	descPadded := xstrings.LeftJustify(desc, descPadLen, " ")
	return cmdPadded + " " + descPadded
}

func layout(g *gocui.Gui) error {
	if v, err := setRectangleView(g, "insert", instructWidth, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "insert"
		v.Editable = true
	}

	if v, err := setRectangleView(g, "apps", instructWidth, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "apps"
		v.Editable = true
	}

	if v, err := setRectangleView(g, "remote", instructWidth, 12); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintln(v, centerInstruct("[ARROWS]", "move"))
		fmt.Fprintln(v, centerInstruct("[ENTER]", "select"))
		fmt.Fprintln(v, centerInstruct("[SPACE]", "play/pause"))
		fmt.Fprintln(v, centerInstruct("[*(8)]", "info (start button)"))
		fmt.Fprintln(v, centerInstruct("[h]", "home"))
		fmt.Fprintln(v, centerInstruct("[b]", "back"))
		fmt.Fprintln(v, centerInstruct("[+(=)/-]", "volume up/down"))
		fmt.Fprintln(v, centerInstruct("[<(,)/>(.)]", "reverse/forward"))
		fmt.Fprintln(v, centerInstruct("[i]", "insert text"))
		fmt.Fprintln(v, centerInstruct("[a]", "launch app"))

		v.Title = "remote"

		if _, err = setCurrentViewOnTop(g, "remote"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	ip = os.Getenv("ROKU")
	rokuClient = roku.New(ip)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	setRemoteKeybindRoku(g, gocui.KeyArrowDown, "down")
	setRemoteKeybindRoku(g, gocui.KeyArrowLeft, "left")
	setRemoteKeybindRoku(g, gocui.KeyArrowRight, "right")
	setRemoteKeybindRoku(g, gocui.KeyArrowUp, "up")
	setRemoteKeybindRoku(g, gocui.KeyEnter, "select")
	setRemoteKeybindRoku(g, gocui.KeySpace, "play")
	setRemoteKeybindRoku(g, castRune("*"), "info")
	setRemoteKeybindRoku(g, castRune("8"), "info")
	setRemoteKeybindRoku(g, castRune("+"), "volume_up")
	setRemoteKeybindRoku(g, castRune("="), "volume_up")
	setRemoteKeybindRoku(g, castRune("-"), "volume_down")
	setRemoteKeybindRoku(g, castRune("<"), "reverse")
	setRemoteKeybindRoku(g, castRune(","), "reverse")
	setRemoteKeybindRoku(g, castRune(">"), "forward")
	setRemoteKeybindRoku(g, castRune("."), "forward")
	if err := g.SetKeybinding("remote", castRune("i"), gocui.ModNone, selectViewHandler("insert")); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("remote", castRune("a"), gocui.ModNone, selectViewHandler("apps")); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("insert", gocui.KeyEnter, gocui.ModNone, sendInsert); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("apps", gocui.KeyEnter, gocui.ModNone, launchApp); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
