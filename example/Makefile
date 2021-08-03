LSS_FILES  := $(wildcard src/*.ls)

LS_FILES = $(subst src,ls,$(patsubst %.ls,%.ls,$(LSS_FILES)))
TP_FILES = $(subst ls,bin,$(patsubst %.ls,%.tp,$(LS_FILES)))

ls/%.ls: src/%.ls doc/spreadsheet.xlsx
	fexcel compile --sheet Sheet1 --dins P2 --douts S2 --rins J2 --routs M2 --numregs A2 --posregs D2 --constants G2 --ualms V2 --noupdate "doc/spreadsheet.xlsx" $< -o $@

bin/%.tp: ls/%.ls
	#tplint $< -I src
	#maketp $< $@
	touch $@

all: ${LS_FILES} ${TP_FILES}

.PHONY: clean

clean:
	rm ls/*.ls
	rm bin/*.tp
