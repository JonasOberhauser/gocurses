package curses

// struct ldat{};
// struct _win_st{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <curses.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"os"
	"unsafe"
	"strings"
)

type void unsafe.Pointer

type Window C.WINDOW

type CursesError struct {
	message string
}

func (ce CursesError) String() string {
	return ce.message
}

// Cursor options.
const (
	CURS_HIDE = iota
	CURS_NORM
	CURS_HIGH
)

// Goroutine safeness channels
var (
	inc  = make(chan bool)
	outc = make(chan bool)
)

// Pointers to the values in curses, which may change values.
var Cols *int = nil
var Rows *int = nil

var Colors *int = nil
var ColorPairs *int = nil

var Tabsize *int = nil

// The window returned from C.initscr()
var Stdwin *Window = nil

// Initializes gocurses
func init() {
	Cols = (*int)(void(&C.COLS))
	Rows = (*int)(void(&C.LINES))

	Colors = (*int)(void(&C.COLORS))
	ColorPairs = (*int)(void(&C.COLOR_PAIRS))

	Tabsize = (*int)(void(&C.TABSIZE))
}

// Manage goroutine safety
func in()  { inc <- true }
func out() { <-outc }
// When something is sent to inc, block until there's a recv on outc
func init() {go func() {for { outc <- <-inc }}()}


func Initscr() (*Window, os.Error) {
	in()
	defer out()
	Stdwin = (*Window)(C.initscr())

	if Stdwin == nil {
		return nil, CursesError{"Initscr failed"}
	}

	return Stdwin, nil
}

func Newwin(cols int, rows int, startx int, starty int) (*Window, os.Error) {
	in()
	defer out()
	nw := (*Window)(C.newwin(C.int(rows), C.int(cols), C.int(starty), C.int(startx)))

	if nw == nil {
		return nil, CursesError{"Failed to create window"}
	}

	return nw, nil
}


func (win *Window) Subwin(cols int, rows int, startx int, starty int) (*Window, os.Error) {
	in()
	defer out()
	sw := (*Window)(C.subwin((*C.WINDOW)(win), C.int(rows), C.int(cols), C.int(starty), C.int(startx)))

	if sw == nil {
		return nil, CursesError{"Failed to create window"}
	}

	return sw, nil
}

func (win *Window) Derwin(cols int, rows int, startx int, starty int) (*Window, os.Error) {
	in()
	defer out()
	dw := (*Window)(C.derwin((*C.WINDOW)(win), C.int(rows), C.int(cols), C.int(starty), C.int(startx)))

	if dw == nil {
		return nil, CursesError{"Failed to create window"}
	}

	return dw, nil
}

func Start_color() os.Error {
	in()
	defer out()
	if int(C.has_colors()) == 0 {
		return CursesError{"terminal does not support color"}
	}
	C.start_color()

	return nil
}

func Init_pair(pair int, fg int, bg int) os.Error {
	in()
	defer out()
	if C.init_pair(C.short(pair), C.short(fg), C.short(bg)) == 0 {
		return CursesError{"Init_pair failed"}
	}
	return nil
}

func Color_pair(pair int) int32 {
	in()
	defer out()
	return int32(C.COLOR_PAIR(C.int(pair)))
}

func Noecho() os.Error {
	in()
	defer out()
	if int(C.noecho()) == 0 {
		return CursesError{"Noecho failed"}
	}
	return nil
}

func Echo() os.Error {
	in()
	defer out()
	if int(C.echo()) == 0 {
		return CursesError{"Echo failed"}
	}
	return nil
}

func Curs_set(c int) os.Error {
	in()
	defer out()
	if C.curs_set(C.int(c)) == 0 {
		return CursesError{"Curs_set failed"}
	}
	return nil
}

func Nocbreak() os.Error {
	in()
	defer out()
	if C.nocbreak() == 0 {
		return CursesError{"Nocbreak failed"}
	}
	return nil
}

func Cbreak() os.Error {
	in()
	defer out()
	if C.cbreak() == 0 {
		return CursesError{"Cbreak failed"}
	}
	return nil
}

func Endwin() os.Error {
	in()
	defer out()
	if C.endwin() == 0 {
		return CursesError{"Endwin failed"}
	}
	return nil
}

func (win *Window) Getch() int {
	return int(C.wgetch((*C.WINDOW)(win)))
}

func (win *Window) Addch(x, y int, c int32, flags int32) {
	in()
	defer out()
	C.mvwaddch((*C.WINDOW)(win), C.int(y), C.int(x), C.chtype(c)|C.chtype(flags))
}


// Internal non-blocking addstr
func (win *Window) addstr ( x,y int, str string, flags int32 ) {
	win.move(x, y)
	for _, ch := range str {
		C.waddch((*C.WINDOW)(win), C.chtype(ch)|C.chtype(flags))
	}
}

// Since CGO currently can't handle varg C functions we'll mimic the
// ncurses addstr functions.
// Per Issue 635 the variadic function definition needs to end with
// 'v ... interface {}' instead of 'v ...'.
func (win *Window) Addstr(x0, y0 int, str string, flags int32, v ...interface{}) {
	in()
	defer out()
	newstr := fmt.Sprintf(str, v)
	win.addstr( x0, y0, newstr, flags )
}

// Like AddStr, just that it tries to nicely align the string within the window.
// A very simplistic algorithm. Whenever a newline starts or the next word is too long, we go
// to the next line and perform a carriage return. 
func (win *Window) AddStrAlign(x0, y0 int, str string, flags int32, v ...interface{}) {
	in()
	defer out()
	lines := strings.Split( fmt.Sprintf(str, v), "\n", -1 )
	
	x,y := x0,y0
	maxx,_ := win.getmax()
	if x0 >= 1 { 
		maxx -= 1 
	}

	for _, line := range lines {
		words := strings.Split( line, " ", -1 )
		for _, word := range words {
			l := len( word )
			if x + l > maxx {
				y += 1 // add one line when word is too long
				x = x0 // carriage return
			}
			win.move( x, y )
			win.addstr( x, y, word, flags )
			x += l + 1 // +1 for space
		}
		y += 1 // add one line after newlines
		x = x0 // carriage return
	}
	
}


// Non-atomic move, needs to be called by atomic funcs in order to prevent deadlocks.
func (win *Window) move(x, y int) {
	// Y normally is the first parameter
	C.wmove((*C.WINDOW)(win), C.int(y), C.int(x))
}

// Normally Y is the first parameter passed in curses.
func (win *Window) Move(x, y int) {
	in()
	defer out()
	win.move(x, y)
}

func (w *Window) Keypad(tf bool) os.Error {
	in()
	defer out()
	var outint int
	if tf == true {
		outint = 1
	}
	if tf == false {
		outint = 0
	}
	if C.keypad((*C.WINDOW)(w), C.int(outint)) == 0 {
		return CursesError{"Keypad failed"}
	}
	return nil
}

func (win *Window) Refresh() os.Error {
	in()
	defer out()
	if C.wrefresh((*C.WINDOW)(win)) == 0 {
		return CursesError{"refresh failed"}
	}
	return nil
}

func (win *Window) Redrawln(beg_line, num_lines int) {
	in()
	defer out()
	C.wredrawln((*C.WINDOW)(win), C.int(beg_line), C.int(num_lines))
}

func (win *Window) Redraw() {
	in()
	defer out()
	C.redrawwin((*C.WINDOW)(win))
}

func (win *Window) Clear() {
	in()
	defer out()
	C.wclear((*C.WINDOW)(win))
}

func (win *Window) Erase() {
	in()
	defer out()
	C.werase((*C.WINDOW)(win))
}

func (win *Window) Clrtobot(x,y int) {
	in()
	defer out()
	win.move(x,y)
	C.wclrtobot((*C.WINDOW)(win))
}

func (win *Window) Clrtoeol(x,y int) {
	in()
	defer out()
	win.move(x,y)
	C.wclrtoeol((*C.WINDOW)(win))
}

func (win *Window) Box(verch, horch int) {
	in()
	defer out()
	C.box((*C.WINDOW)(win), C.chtype(verch), C.chtype(horch))
}

func (win *Window) Background(colour int32) {
	in()
	defer out()
	C.wbkgd((*C.WINDOW)(win), C.chtype(colour))
}

func (win *Window) Noutrefresh() {
	in()
	defer out()
	C.wnoutrefresh((*C.WINDOW)(win))
}
func (win *Window) Touchwin() {
	in()
	defer out()
	C.touchwin((*C.WINDOW)(win))
}

func (win *Window) Getbeg() (x, y int) {
	in()
	defer out()
	x = int(C.getbegx((*C.WINDOW)(win)))
	y = int(C.getbegy((*C.WINDOW)(win)))
	return x, y
}


// non-atomic getmax. Doesn't block and has to be used as part of other atomic functions to prevent deadlocks.
func (win *Window) getmax() (x, y int) {
	x = int(C.getmaxx((*C.WINDOW)(win)))
	y = int(C.getmaxy((*C.WINDOW)(win)))
	return x, y
}

func (win *Window) Getmax() (x, y int) {
	in()
	defer out()
	return win.getmax()
}

func (win *Window) Getstr() (str string, err os.Error) {
	cstr := C.CString(str)
	//defer C.free(unsafe.Pointer(cstr))
	if C.wgetstr((*C.WINDOW)(win), cstr) == -1 {
		return "", CursesError{"wgetstr failed"}
	}
	s := C.GoString(cstr)
	return s, nil
}

func DoUpdate() {
	in()
	defer out()
	C.doupdate()
}
