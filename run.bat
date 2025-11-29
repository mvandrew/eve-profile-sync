@echo off
echo Running eve-profile-sync...
go run main.go %*

if %ERRORLEVEL% NEQ 0 (
    echo Run failed!
    exit /b %ERRORLEVEL%
)

