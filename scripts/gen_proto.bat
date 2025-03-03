@echo off
setlocal enabledelayedexpansion

echo Checking if protoc is installed...
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo protoc not found, please install Protocol Buffers compiler first
    echo You can download it from: https://github.com/protocolbuffers/protobuf/releases
    exit /b 1
)

echo Checking if protoc-gen-go is installed...
where protoc-gen-go >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing protoc-gen-go...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    if %errorlevel% neq 0 (
        echo Failed to install protoc-gen-go
        exit /b 1
    )
)

echo Starting Protocol Buffers code generation...

set PROJECT_ROOT=%~dp0..
set PROTO_PATH=%PROJECT_ROOT%
set GO_OUT=%PROJECT_ROOT%

cd %PROJECT_ROOT%
protoc --proto_path=. --go_out=. --go_opt=paths=source_relative internal/transport/api/v1/api.proto
if %errorlevel% neq 0 (
    echo Generation failed
    exit /b 1
)

echo Protocol Buffers code generation completed!
exit /b 0