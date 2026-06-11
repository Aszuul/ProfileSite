@echo off
REM Build Go WASM binary for Breakout game
cd /d "%~dp0go"
echo Building Go WASM binary...
set GOOS=js
set GOARCH=wasm
go build -o ..\static\wasm\main.wasm
if %ERRORLEVEL% equ 0 (
    echo.
    echo Build successful! WASM binary saved to static/wasm/main.wasm
    echo.
    echo Copying wasm_exec.js...
    for /f "tokens=*" %%i in ('go env GOROOT') do set GOROOT=%%i
    if exist "%GOROOT%\lib\wasm\wasm_exec.js" (
        copy "%GOROOT%\lib\wasm\wasm_exec.js" "..\static\wasm\wasm_exec.js" >nul
        echo wasm_exec.js copied successfully.
    ) else (
        echo Warning: wasm_exec.js not found at !GOROOT!\lib\wasm\
    )
) else (
    echo.
    echo Build failed! Check error messages above.
    exit /b 1
)
