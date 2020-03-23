#!/bin/bash

pass=0
fail=0

for i in $(find $1 -iname "*.sysl"); do
	if $(wbnf test --input=$i --grammar pkg/parser/sysl.wbnf --start sysl_file 1>/dev/null 2> /dev/null); then
	  echo $i   "success"
    pass=$((pass+1))
	else
	  echo $i   "fail"
    fail=$((fail+1))
  fi
done

echo Passes:$pass Fails:$fail
