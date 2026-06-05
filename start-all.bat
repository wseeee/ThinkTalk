@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

cd /d "%~dp0"

set BIN_DIR=%~dp0bin
set LOG_DIR=%~dp0logs

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

echo.
echo ========================================
echo    ThinkTalk - Local One-Click Startup
echo ========================================
echo.

REM ---------- Prerequisites ----------
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Go SDK not found, please install Go first.
    pause
    exit /b 1
)

if not exist .env (
    echo [WARN] .env does not exist, copying from .env.example...
    copy .env.example .env >nul
    echo [INFO] .env created. Please configure it and rerun.
    pause
    exit /b 0
)

REM ---------- Compilation ----------
echo [1/2] Compiling services...
echo.

REM === API ===
echo   Compiling applet-api...
go build -o "%BIN_DIR%\applet-api.exe" ./application/applet
if %errorlevel% neq 0 goto :build_error

echo   Compiling article-api...
go build -o "%BIN_DIR%\article-api.exe" ./application/article/api
if %errorlevel% neq 0 goto :build_error

echo   Compiling chat-api...
go build -o "%BIN_DIR%\chat-api.exe" ./application/chat/api
if %errorlevel% neq 0 goto :build_error

echo   Compiling qa-api...
go build -o "%BIN_DIR%\qa-api.exe" ./application/qa/api
if %errorlevel% neq 0 goto :build_error

REM === RPC ===
echo   Compiling user-rpc...
go build -o "%BIN_DIR%\user-rpc.exe" ./application/user/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling article-rpc...
go build -o "%BIN_DIR%\article-rpc.exe" ./application/article/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling chat-rpc...
go build -o "%BIN_DIR%\chat-rpc.exe" ./application/chat/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling concerned-rpc...
go build -o "%BIN_DIR%\concerned-rpc.exe" ./application/concerned/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling follow-rpc...
go build -o "%BIN_DIR%\follow-rpc.exe" ./application/follow/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling like-rpc...
go build -o "%BIN_DIR%\like-rpc.exe" ./application/like/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling member-rpc...
go build -o "%BIN_DIR%\member-rpc.exe" ./application/member/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling message-rpc...
go build -o "%BIN_DIR%\message-rpc.exe" ./application/message/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling qa-rpc...
go build -o "%BIN_DIR%\qa-rpc.exe" ./application/qa/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling reply-rpc...
go build -o "%BIN_DIR%\reply-rpc.exe" ./application/reply/rpc
if %errorlevel% neq 0 goto :build_error

echo   Compiling tag-rpc...
go build -o "%BIN_DIR%\tag-rpc.exe" ./application/tag/rpc
if %errorlevel% neq 0 goto :build_error

REM === MQ ===
echo   Compiling article-mq...
go build -o "%BIN_DIR%\article-mq.exe" ./application/article/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling chat-mq...
go build -o "%BIN_DIR%\chat-mq.exe" ./application/chat/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling concerned-mq...
go build -o "%BIN_DIR%\concerned-mq.exe" ./application/concerned/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling like-mq...
go build -o "%BIN_DIR%\like-mq.exe" ./application/like/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling member-mq...
go build -o "%BIN_DIR%\member-mq.exe" ./application/member/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling message-mq...
go build -o "%BIN_DIR%\message-mq.exe" ./application/message/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling qa-mq...
go build -o "%BIN_DIR%\qa-mq.exe" ./application/qa/mq
if %errorlevel% neq 0 goto :build_error

echo   Compiling reply-mq...
go build -o "%BIN_DIR%\reply-mq.exe" ./application/reply/mq
if %errorlevel% neq 0 goto :build_error

echo.
echo [INFO] All services compiled successfully!
echo.

REM ---------- Startup ----------
echo [2/2] Starting services...
echo.

REM Clear PID record file
type nul > "%LOG_DIR%\pids.txt"

REM === API ===
if exist "%BIN_DIR%\applet-api.exe" (
    start "applet-api" /B "%BIN_DIR%\applet-api.exe" -f application\applet\etc\applet-api.yaml > "%LOG_DIR%\applet-api.log" 2>&1
    call :record_pid applet-api
) else (
    echo [WARN] applet-api.exe does not exist, skipping.
)

if exist "%BIN_DIR%\article-api.exe" (
    start "article-api" /B "%BIN_DIR%\article-api.exe" -f application\article\api\etc\article-api.yaml > "%LOG_DIR%\article-api.log" 2>&1
    call :record_pid article-api
) else (
    echo [WARN] article-api.exe does not exist, skipping.
)

if exist "%BIN_DIR%\chat-api.exe" (
    start "chat-api" /B "%BIN_DIR%\chat-api.exe" -f application\chat\api\etc\chat-api.yaml > "%LOG_DIR%\chat-api.log" 2>&1
    call :record_pid chat-api
) else (
    echo [WARN] chat-api.exe does not exist, skipping.
)

if exist "%BIN_DIR%\qa-api.exe" (
    start "qa-api" /B "%BIN_DIR%\qa-api.exe" -f application\qa\api\etc\qa-api.yaml > "%LOG_DIR%\qa-api.log" 2>&1
    call :record_pid qa-api
) else (
    echo [WARN] qa-api.exe does not exist, skipping.
)

REM === RPC ===
if exist "%BIN_DIR%\user-rpc.exe" (
    start "user-rpc" /B "%BIN_DIR%\user-rpc.exe" -f application\user\rpc\etc\user.yaml > "%LOG_DIR%\user-rpc.log" 2>&1
    call :record_pid user-rpc
) else (
    echo [WARN] user-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\article-rpc.exe" (
    start "article-rpc" /B "%BIN_DIR%\article-rpc.exe" -f application\article\rpc\etc\article.yaml > "%LOG_DIR%\article-rpc.log" 2>&1
    call :record_pid article-rpc
) else (
    echo [WARN] article-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\chat-rpc.exe" (
    start "chat-rpc" /B "%BIN_DIR%\chat-rpc.exe" -f application\chat\rpc\etc\chat.yaml > "%LOG_DIR%\chat-rpc.log" 2>&1
    call :record_pid chat-rpc
) else (
    echo [WARN] chat-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\concerned-rpc.exe" (
    start "concerned-rpc" /B "%BIN_DIR%\concerned-rpc.exe" -f application\concerned\rpc\etc\concerned.yaml > "%LOG_DIR%\concerned-rpc.log" 2>&1
    call :record_pid concerned-rpc
) else (
    echo [WARN] concerned-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\follow-rpc.exe" (
    start "follow-rpc" /B "%BIN_DIR%\follow-rpc.exe" -f application\follow\rpc\etc\follow.yaml > "%LOG_DIR%\follow-rpc.log" 2>&1
    call :record_pid follow-rpc
) else (
    echo [WARN] follow-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\like-rpc.exe" (
    start "like-rpc" /B "%BIN_DIR%\like-rpc.exe" -f application\like\rpc\etc\like.yaml > "%LOG_DIR%\like-rpc.log" 2>&1
    call :record_pid like-rpc
) else (
    echo [WARN] like-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\member-rpc.exe" (
    start "member-rpc" /B "%BIN_DIR%\member-rpc.exe" -f application\member\rpc\etc\member.yaml > "%LOG_DIR%\member-rpc.log" 2>&1
    call :record_pid member-rpc
) else (
    echo [WARN] member-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\message-rpc.exe" (
    start "message-rpc" /B "%BIN_DIR%\message-rpc.exe" -f application\message\rpc\etc\message.yaml > "%LOG_DIR%\message-rpc.log" 2>&1
    call :record_pid message-rpc
) else (
    echo [WARN] message-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\qa-rpc.exe" (
    start "qa-rpc" /B "%BIN_DIR%\qa-rpc.exe" -f application\qa\rpc\etc\qa.yaml > "%LOG_DIR%\qa-rpc.log" 2>&1
    call :record_pid qa-rpc
) else (
    echo [WARN] qa-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\reply-rpc.exe" (
    start "reply-rpc" /B "%BIN_DIR%\reply-rpc.exe" -f application\reply\rpc\etc\reply.yaml > "%LOG_DIR%\reply-rpc.log" 2>&1
    call :record_pid reply-rpc
) else (
    echo [WARN] reply-rpc.exe does not exist, skipping.
)

if exist "%BIN_DIR%\tag-rpc.exe" (
    start "tag-rpc" /B "%BIN_DIR%\tag-rpc.exe" -f application\tag\rpc\etc\tag.yaml > "%LOG_DIR%\tag-rpc.log" 2>&1
    call :record_pid tag-rpc
) else (
    echo [WARN] tag-rpc.exe does not exist, skipping.
)

REM === MQ ===
if exist "%BIN_DIR%\article-mq.exe" (
    start "article-mq" /B "%BIN_DIR%\article-mq.exe" -f application\article\mq\etc\article.yaml > "%LOG_DIR%\article-mq.log" 2>&1
    call :record_pid article-mq
) else (
    echo [WARN] article-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\chat-mq.exe" (
    start "chat-mq" /B "%BIN_DIR%\chat-mq.exe" -f application\chat\mq\etc\chat.yaml > "%LOG_DIR%\chat-mq.log" 2>&1
    call :record_pid chat-mq
) else (
    echo [WARN] chat-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\concerned-mq.exe" (
    start "concerned-mq" /B "%BIN_DIR%\concerned-mq.exe" -f application\concerned\mq\etc\concerned.yaml > "%LOG_DIR%\concerned-mq.log" 2>&1
    call :record_pid concerned-mq
) else (
    echo [WARN] concerned-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\like-mq.exe" (
    start "like-mq" /B "%BIN_DIR%\like-mq.exe" -f application\like\mq\etc\like.yaml > "%LOG_DIR%\like-mq.log" 2>&1
    call :record_pid like-mq
) else (
    echo [WARN] like-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\member-mq.exe" (
    start "member-mq" /B "%BIN_DIR%\member-mq.exe" -f application\member\mq\etc\member.yaml > "%LOG_DIR%\member-mq.log" 2>&1
    call :record_pid member-mq
) else (
    echo [WARN] member-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\message-mq.exe" (
    start "message-mq" /B "%BIN_DIR%\message-mq.exe" -f application\message\mq\etc\message.yaml > "%LOG_DIR%\message-mq.log" 2>&1
    call :record_pid message-mq
) else (
    echo [WARN] message-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\qa-mq.exe" (
    start "qa-mq" /B "%BIN_DIR%\qa-mq.exe" -f application\qa\mq\etc\qa.yaml > "%LOG_DIR%\qa-mq.log" 2>&1
    call :record_pid qa-mq
) else (
    echo [WARN] qa-mq.exe does not exist, skipping.
)

if exist "%BIN_DIR%\reply-mq.exe" (
    start "reply-mq" /B "%BIN_DIR%\reply-mq.exe" -f application\reply\mq\etc\reply.yaml > "%LOG_DIR%\reply-mq.log" 2>&1
    call :record_pid reply-mq
) else (
    echo [WARN] reply-mq.exe does not exist, skipping.
)

echo.
echo ========================================
echo    All services started successfully!
echo ========================================
echo.
echo   API Services:
echo     applet-api   -^> http://localhost:8888
echo     article-api  -^> http://localhost:80
echo     chat-api     -^> http://localhost:8087
echo     qa-api       -^> http://localhost:8890
echo.
echo   View logs: type logs\service-name.log
echo   Stop all:  stop-all.bat
echo.

pause
goto :eof

REM ---------- Subroutines ----------

REM Record process name to pids.txt
:record_pid
echo %1>>"%LOG_DIR%\pids.txt"
goto :eof

REM Compilation error handler
:build_error
echo.
echo [ERROR] Compilation failed! Please check error messages above.
pause
exit /b 1
