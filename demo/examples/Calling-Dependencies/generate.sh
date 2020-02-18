#!/bin/bash

make input=simple.sysl app=simple down=mydependency

ls 
#gen             implementation  makefile        readme.md
#generate.sh     main.go         ordering.yaml   simple.sysl
# Now we can see that there is a "gen" folder with our generated sysl files

cd gen && ls
#Simple       mydependency
# We can see here that the code generated both "mydependency" and "simple"