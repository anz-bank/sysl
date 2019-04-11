@ECHO OFF
SETLOCAL EnableDelayedExpansion

set ROOT=sysl2\sysl\tests
for /R %ROOT% %%f in (*.sysl) do (
	dist\gosysl.exe -mode textpb -root %ROOT% -o %ROOT%\%%~nf.win.txt -log debug /%%~nf.sysl || exit /b !errorlevel!
)

del sysl2\sysl\tests\*.win.txt
