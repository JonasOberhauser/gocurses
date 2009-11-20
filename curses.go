package curses

// struct ldat{};
// struct _win_st{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <curses.h>
import "C"

import (
	//"os";
	"unsafe";
)

type void unsafe.Pointer;

type Window C.WINDOW;

type CursesError struct {
	message string;
}

func (ce CursesError) String() string {
	return ce.message;
}

// Cursor options.
const (
	CURS_HIDE = iota;
	CURS_NORM;
	CURS_HIGH;
)

// Pointers to the values in curses, which may change.
var Cols *int = nil;
var Rows *int = nil;

var Colors *int = nil;
var ColorPairs *int = nil;

var Tabsize *int = nil;

// The window returned from C.initscr()
var Stdwin *Window = nil;

// Initializes gocurses
func init() {
	Cols = (*int)(void(&C.COLS));
	Rows = (*int)(void(&C.LINES));
	
	Colors = (*int)(void(&C.COLORS));
	ColorPairs = (*int)(void(&C.COLOR_PAIRS));
	
	Tabsize = (*int)(void(&C.TABSIZE));
}

func Initscr() *Window {
	Stdwin = (*Window)(C.initscr());
	
	return Stdwin;
}

func Noecho() {
	C.noecho();
}

func Echo() {
	C.echo();
}

func Curs_set(c int) {
	C.curs_set(C.int(c));
}

func Cbreak() {
	C.cbreak();
}

func Endwin() {
	C.endwin();
}

func (win *Window) Getch() int {
	return int(C.wgetch((*C.WINDOW)(win)));
}

func (win *Window) Addch(x, y int, c int32) {
	C.mvwaddch((*C.WINDOW)(win), C.int(y), C.int(x), C.chtype(c));
}

// Normally Y is the first parameter passed in curses.
func (win *Window) Move(x, y int) {
	C.wmove((*C.WINDOW)(win), C.int(y), C.int(x));
}

func (w *Window) Keypad(tf bool) {
	var outint int;
	if tf == true {outint = 1;}
	if tf == false {outint = 0;}
	C.keypad((*C.WINDOW)(w), C.int(outint));
}

func (win *Window) Refresh() {
	C.wrefresh((*C.WINDOW)(win));
}

func (win *Window) Redrawln(beg_line, num_lines int) {
	C.wredrawln((*C.WINDOW)(win), C.int(beg_line), C.int(num_lines));
}

func (win *Window) Redraw() {
	C.redrawwin((*C.WINDOW)(win));
}

func (win *Window) Clear() {
	C.wclear((*C.WINDOW)(win));
}

func (win *Window) Erase() {
	C.werase((*C.WINDOW)(win));
}

func (win *Window) Clrtobot() {
	C.wclrtobot((*C.WINDOW)(win));
}

func (win *Window) Clrtoeol() {
	C.wclrtoeol((*C.WINDOW)(win));
}
