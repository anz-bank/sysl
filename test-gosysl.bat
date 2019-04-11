@ECHO OFF
SETLOCAL EnableDelayedExpansion

set ROOT=sysl2\sysl\tests
for /R %ROOT% %%f in (*.sysl) do (
	gosysl\gosysl.exe -mode textpb -root %ROOT% -o %ROOT%\%%~nf.txt -log debug /%%~nf.sysl || exit /b !errorlevel!
)

