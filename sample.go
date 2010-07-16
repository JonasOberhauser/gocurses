package main

import . "curses"
import "os"
import "fmt"

func main() {
	x := 10
	y := 10
	startGoCurses()
	defer stopGoCurses()
	Init_pair(1, COLOR_RED, COLOR_BLACK)

	loop(x, y)
}

func startGoCurses() {
	Initscr()
	if Stdwin == nil {
		stopGoCurses()
		os.Exit(1)
	}

	Noecho()

	Curs_set(CURS_HIDE)
	Stdwin.Keypad(true)

	if err := Start_color(); err != nil {
		fmt.Printf("%s\n", err.String())
		stopGoCurses()
		os.Exit(1)
	}

}

func stopGoCurses() {
	Endwin()
}

func loop(x, y int) {
	for {
		Stdwin.Addstr(0, 0, "Hello,\nworld!", 0)
		inp := Stdwin.Getch()
		if inp == 'q' {
			break
		}
		if inp == KEY_LEFT {
			x = x - 1
		}
		if inp == KEY_RIGHT {
			x = x + 1
		}
		if inp == KEY_UP {
			y = y - 1
		}
		if inp == KEY_DOWN {
			y = y + 1
		}
		Stdwin.Clear()
		Stdwin.Addch(x, y, '@', Color_pair(1))
		Stdwin.Refresh()
	}
}
