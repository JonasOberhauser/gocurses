package main

import . "curses"

func main() {
	x := 10;
	y := 10;
	Initscr();
	Noecho();
	Stdwin.Addch(x, y, '@');
	for Getch() != 'q' {  }
	Endwin();
}
