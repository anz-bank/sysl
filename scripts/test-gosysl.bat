@ECHO OFF
SETLOCAL EnableDelayedExpansion

set ROOT=tests
for /R %ROOT% %%f in (*.sysl) do (
	bin\sysl.exe --log debug pb --mode textpb --root %ROOT% -o %ROOT%\%%~nf.win.txt  /%%~nf.sysl || exit /b !errorlevel!
)

del tests\*.win.txt
