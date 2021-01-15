# Makefile for BBC Micro B Bruce Lee NuLA Version
# 15 January 2021
#

BEEBEM       := C:/Dave/b2/b2_Debug
PNG2BBC      := ../png2bbc/Release/png2bbc.exe
EMACS        := c:/Dave/Emacs-27.1/bin/emacs.exe
SNAP         := ../snap/Release/snap.exe
BEEBEM_FLAGS := -b -0
BEEBASM      := ../beebasm/beebasm.exe
GAME_SSD     := res/blank.ssd
OUTPUT_SSD   := bruce-nula.ssd
MAIN_ASM     := main.asm
RM           := del
CP           := copy

#
# Generated graphics
GFX_OBJECTS := $(shell $(PNG2BBC) -l gfxscript)

#
# Phony targets
.PHONY: all clean run gfx

all: $(OUTPUT_SSD)

$(OUTPUT_SSD): $(MAIN_ASM) Makefile
	$(BEEBASM) -i $(MAIN_ASM) -di $(GAME_SSD) -do $(OUTPUT_SSD)

#$(GFX_OBJECTS): gfxscript
#	$(PNG2BBC) gfxscript

#gfx:
#	$(PNG2BBC) gfxscript
#	$(EMACS) -batch -Q --eval="(package-initialize)" -l repack.el --eval="(reverse-graphic \"bin/fuel.bin\" \"bin/fuel.bbc\")"
#	$(CP) bin\game.pal.new bin\game.pal
#	$(RM) bin\game.pal.new
#	$(SNAP) org/jet-pac bin/platform.bin 7680 bin/jet-pac-nula

clean:
	$(RM) $(OUTPUT_SSD)
	$(RM) /Q bin\*.*

run:
	$(BEEBEM) $(BEEBEM_FLAGS) $(OUTPUT_SSD)
