package curses

// struct ldat{};
// struct _win_st{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <curses.h>
import "C"

type Window C.WINDOW;

func Initscr() *Window {
	return (*Window)(C.initscr());
}

func Getch() int {
	return int(C.getch());
}

func Endwin() {
	C.endwin();
}
