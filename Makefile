include $(GOROOT)/src/Make.inc

TARG=curses

GOFILES=curses_defs.go
CGOFILES=curses.go \
	panel.go \

CGO_LDFLAGS=-lpanel -lncurses 

CLEANFILES+=

include $(GOROOT)/src/Make.pkg

curses_defs.go: curses.c
	godefs -g curses curses.c > curses_defs.go

# Simple test programs

%: install %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O
