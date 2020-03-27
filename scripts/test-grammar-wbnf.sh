#!/bin/bash

pass=0
fail=()

for i in $(find $1 -iname "*.sysl"); do
	if $(wbnf test --input=$i --grammar pkg/parser/sysl.wbnf --start sysl_file 1>/dev/null 2> /dev/null); then
	  echo "success" $i
    pass=$((pass+1))
	else
    fail=("${fail[@]}" $i)
  fi
done

for i in ${fail[@]}; do
  echo "fail   " $i
done

echo Passes:$pass Fails:${#fail[@]}
