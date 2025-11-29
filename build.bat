@echo off
echo Building eve-profile-sync...
go build -o eve-profile-sync.exe

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    exit /b %ERRORLEVEL%
)

echo Build successful! Binary created: eve-profile-sync.exe

