package curses

// struct ldat{};
// struct _win_st{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <curses.h>
import "C"

import "unsafe"

type void unsafe.Pointer;

type Window C.WINDOW;

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

func Getch() int {
	return int(C.getch());
}

func Endwin() {
	C.endwin();
}

func (w *Window) Addch(y, x int, c int32) {
	C.mvwaddch((*C.WINDOW)(w), C.int(y), C.int(x), C.chtype(c));
}

func Keypad(win *Window, tf bool) {
	var outint int;
	if tf == true {outint = 1;}
	if tf == false {outint = 0;}
	C.keypad((*C.WINDOW)(win), C.int(outint));
}