#!/bin/bash
# Make sure that you use the Makefile that is used in this demonstration, (or go to https://github.com/anz-bank/sysl-template/blob/master/makefile)

# We need to edit our makefile and add the following lines at the top:
all: sysl

input = simple.sysl
app = simple
down =  # this can be a list separated by a space or left empty
out = gen
# Current go import path, this should be what's in your go.mod plus the path you're in
basepath = github.service.anz/sysl/syslbyexample/_examples/Code-Generation

make input=simple.sysl app=simple

ls 
#gen             implementation  makefile        readme.md
#generate.sh     main.go         ordering.yaml   simple.sysl
# Now we can see that there is a "gen" folder with our generated sysl files
