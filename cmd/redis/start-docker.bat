@echo off

docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not running. Starting Docker Desktop...
    start /min "" "C:\Program Files\Docker\Docker\Docker Desktop.exe"

    REM Ожидание запуска Docker
    :wait_docker
    docker info >nul 2>&1
    if %errorlevel% neq 0 (
        echo Waiting for Docker to start...
        timeout /t 1 >nul
        goto wait_docker
    )
)

echo Docker is running. Proceeding with the task...