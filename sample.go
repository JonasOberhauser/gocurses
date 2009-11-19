package main

import . "curses"

func main() {
	Initscr();
	Noecho();
	Stdwin.Addch(10, 10, '@');
	for Getch() != 'q' {  }
	Endwin();
}
