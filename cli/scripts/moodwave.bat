@echo off
:: Moodwave CLI — Windows batch bootstrap
:: Double-click or run from Command Prompt.
:: This script launches moodwave.exe if it exists in PATH or local directory.

setlocal enabledelayedexpansion

set "BINARY=moodwave.exe"
set "INSTALL_DIR=%LOCALAPPDATA%\moodwave\bin"

:: Check if moodwave is in PATH
where %BINARY% >nul 2>&1
if %ERRORLEVEL% == 0 (
    %BINARY% %*
    goto :eof
)

:: Check local directory
if exist "%~dp0%BINARY%" (
    "%~dp0%BINARY%" %*
    goto :eof
)

:: Check install directory
if exist "%INSTALL_DIR%\%BINARY%" (
    "%INSTALL_DIR%\%BINARY%" %*
    goto :eof
)

:: Not found — prompt to install
echo.
echo  Moodwave CLI not found.
echo.
echo  To install, run this command in PowerShell:
echo.
echo    irm https://raw.githubusercontent.com/moodwave/moodwave/main/scripts/install.ps1 ^| iex
echo.
echo  Or build from source:
echo    git clone https://github.com/moodwave/moodwave
echo    cd moodwave
echo    go build -o moodwave.exe ./cmd/moodwave
echo.
pause
exit /b 1

:eof
endlocal
