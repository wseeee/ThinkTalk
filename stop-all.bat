@echo off
cd /d "%~dp0"

if not exist "logs\pids.txt" (
    echo [ERROR] No running services found - logs\pids.txt does not exist.
    pause
    exit /b 1
)

echo [INFO] Stopping all ThinkTalk services...
echo.

for /f "tokens=1,2" %%a in (logs\pids.txt) do (
    if "%%b" == "" (
        taskkill /F /FI "IMAGENAME eq %%a.exe" 2>nul
        if errorlevel 1 (
            echo   - %%a does not exist
        ) else (
            echo   - Stopped %%a
        )
    ) else (
        taskkill /F /PID %%a 2>nul
        if errorlevel 1 (
            taskkill /F /FI "IMAGENAME eq %%b.exe" 2>nul
            if errorlevel 1 (
                echo   - %%b does not exist
            ) else (
                echo   - Stopped %%b by name
            )
        ) else (
            echo   - Stopped %%b PID %%a
        )
    )
)

del logs\pids.txt 2>nul

echo.
echo [INFO] Stop completed.
pause
