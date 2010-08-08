package main

import . "curses"
import "os"
import "fmt"

func main() {
	startGoCurses()
	defer stopGoCurses()
	Init_pair(1, COLOR_RED, COLOR_BLACK)

	x := *Cols / 2
	y := *Rows / 2
	
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
	
	Stdwin.Addstr( 0, 3, "Cols %v / Rows %v ", 0, *Cols, *Rows )

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
	Stdwin.Addstr(0, 0, "Hello,", A_NORMAL)
	Stdwin.Addstr(0, 1, "world!", A_NORMAL)
	Stdwin.Addstr(3, 2, "Press p for panels or the famous 'any' for a moving cursor test.", A_NORMAL)
	if inp := Stdwin.Getch(); inp == 'p' {
		Stdwin.Clear()
		
		// The moving panel
		w, _ := Newwin( 20, 12, x, y )
		w.Box( '|', '-' )
		p, _ := NewPanel( w )
		
		w.AddstrAlign(1, 1, "Press q to quit.", A_NORMAL)
		w.AddstrAlign(1, 2, "Press any of the Arrow-Keys to move this window.\nPress t to toggle if this window is on the top or not.", A_NORMAL)
		
		
		// The other panel
		w2, _ := Newwin( 30, 20, 2, 2 )
		w2.Box( '|', '-' )
		p2, _ := NewPanel( w2 )
		
		w2.AddstrAlign(1, 1, "This is another window for you to look at...", A_NORMAL)
		w2.AddstrAlign(1, 3, "Press h to hide this window", A_NORMAL)
		w2.AddstrAlign(1, 0, "YET-ANOTHER-WINDOW", A_NORMAL)
		w2.AddstrAlign(1, 4, "Below: %v\nAbove: %v", 0,  p2.Below(), p2.Above() )
		
		DoUpdate()
		
		top := false
		for ; inp != 'q'; inp = Stdwin.Getch()  {
			switch inp {
			case KEY_LEFT: x -= 1
			case KEY_RIGHT: x += 1
			case KEY_UP: y -= 1 
			case KEY_DOWN: y += 1
			case 't':
				top = ! top
				if top {
					p.ToBottom()
				} else {
					p.ToTop()
				} 
			case 'h': 
				p2.Hide( !p2.Hidden() )
			case 'r': 
				p2.Delete()
			}
			
			maxx,maxy := Stdwin.Getmax()
			winx,winy := w.Getmax()
			maxx-= winx-1
			maxy-= winy-1 
			x=(x+maxx)%maxx
			y=(y+maxy)%maxy
			
			p.Move( x, y )
			UpdatePanels()
			DoUpdate()
			//Stdwin.Refresh()
		}
	} else {
		Stdwin.Clear()
		for ; inp != 'q'; inp = Stdwin.Getch()  {
			switch inp {
			case KEY_LEFT: x -= 1
			case KEY_RIGHT: x += 1
			case KEY_UP: y -= 1 
			case KEY_DOWN: y += 1
			}
			maxx,maxy := Stdwin.Getmax()
			x=(x+maxx)%maxx
			y=(y+maxy)%maxy
			
			Stdwin.Clear()
			Stdwin.Addstr(10, 1, "Press q to quit.", A_NORMAL)
			Stdwin.Addstr(10, 2, "Press any of the Arrow-Keys\n to move the cursor. (Using Addstr)", A_NORMAL)
			Stdwin.AddstrAlign(10, 5, "Press any of the Arrow-Keys\n to move the cursor. (Using Align)", A_NORMAL)
			Stdwin.Addch(x, y, '@', Color_pair(1))
			Stdwin.Refresh()
		}
	}
}
