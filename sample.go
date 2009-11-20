package main

import . "curses"

func main() {
	x := 10;
	y := 10;
	Initscr();
	Noecho();
	Keypad(Stdwin, true);
	input(x, y);
	Endwin();
}

func input(x, y int) {
	for {
		inp := Getch();
		if inp == 'q' {
			break;
		}
		if inp == KEY_LEFT {
			x = x - 1;
		}
		if inp == KEY_RIGHT {
			x = x + 1;
		}
		if inp == KEY_UP {
			y = y - 1;
		}
		if inp == KEY_DOWN {
			y = y + 1;
		}
		Stdwin.Clear();
		Stdwin.Addch(x, y, '@');
		Stdwin.Refresh();
	}
}
