@echo off

docker start redis >nul 2>&1
if %errorlevel% neq 0 (
    echo Redis container is not exists. Starting Redis Container...
    docker run --name redis -d -p 6379:6379 redis

    :wait_run
    docker start redis >nul 2>&1
    if %errorlevel% neq 0 (
        echo Waiting for Redis container to start...
        timeout /t 1 >nul
        goto wait_run
    )
)

echo Redis container is running. Proceeding with the task...