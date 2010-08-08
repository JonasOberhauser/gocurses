package curses

// struct ldat{};
// struct _win_st{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <ncurses/panel.h>
// #include <stdlib.h>
import "C"

import (
// 	"fmt"
	"os"
// 	
	"unsafe"
)

var panelByPtr = make( map[ uintptr ] *Panel )

type Panel struct {
	panel *C.PANEL
	hidden bool
}



/*

	extern NCURSES_EXPORT(PANEL*)  new_panel (WINDOW *);
	extern NCURSES_EXPORT(PANEL*)  panel_above (const PANEL *);
	extern NCURSES_EXPORT(PANEL*)  panel_below (const PANEL *);
	extern NCURSES_EXPORT(int)     set_panel_userptr (PANEL *, NCURSES_CONST void *);
	extern NCURSES_EXPORT(NCURSES_CONST void*) panel_userptr (const PANEL *);
	extern NCURSES_EXPORT(int)     move_panel (PANEL *, int, int);
	extern NCURSES_EXPORT(int)     replace_panel (PANEL *,WINDOW *);
	extern NCURSES_EXPORT(int)     panel_hidden (const PANEL *);


*/

func (p *Panel) Window () (w *Window, err os.Error) {
	in()
	defer out()
	w=(*Window)(C.panel_window( p.panel ))
	if w==nil { err = CursesError{ "No Window" } }
	return
}

func UpdatePanels() {
	in()
	defer out()
	C.update_panels()
}

func (p *Panel) Hide (b bool) {
	in()
	defer out()
	if b {
		C.hide_panel( p.panel )
	} else {
		C.show_panel( p.panel )
	}
	p.hidden = b
}

func (p *Panel) Delete() {
	in()
	defer out()
	C.del_panel( p.panel )
}


func (p *Panel) ToTop() {
	in()
	defer out()
	C.top_panel( p.panel )
}

func (p *Panel) ToBottom() {
	in()
	defer out()
	C.top_panel( p.panel )
}

func NewPanel( w *Window ) (p *Panel, err os.Error ) {
	in()
	defer out()
	p = new(Panel)
	p.panel = C.new_panel( (*C.WINDOW)(w) )
	if p.panel == nil {
		err = CursesError{ "Could not allocate PANEL" }
		p = nil
	} else { 
		panelByPtr[uintptr(unsafe.Pointer(p.panel) )] = p
		p.hidden = false
	}
	return
}

func (p *Panel) Above() *Panel {
	in(); defer out()
	return panelByPtr[ uintptr(unsafe.Pointer(C.panel_above( p.panel ))) ]
}

func (p *Panel) Below() *Panel {
	in(); defer out()
	return panelByPtr[ uintptr(unsafe.Pointer(C.panel_below( p.panel ))) ]
}

func (p *Panel) Move(x,y int) ( err os.Error ) {
	in(); defer out()
	
	if (int)(C.move_panel( p.panel, C.int( y ), C.int( x ) )) == C.ERR {
		err = CursesError{ "Can't move PANEL" }
	}
	return
}

func (p *Panel) Hidden() bool {
	return p.hidden
}
