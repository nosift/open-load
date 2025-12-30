@echo off
REM 查看 API Keys 的组织验证状态

SET DB_PATH=.\data\open-load.db

REM 如果 open-load.db 不存在，尝试 gpt-load.db
IF NOT EXIST "%DB_PATH%" SET DB_PATH=.\data\gpt-load.db

IF NOT EXIST "%DB_PATH%" (
    echo ❌ 数据库文件不存在
    exit /b 1
)

echo 📊 API Keys 组织验证状态统计
echo ================================
echo.

REM 使用 curl 或直接查询（需要安装 sqlite3）
REM 这里提供 API 查询方式

echo 请使用以下方式查看 Key 的组织验证状态：
echo.
echo 方式 1: 通过管理界面
echo   - 访问 http://localhost:3001
echo   - 进入"密钥管理"页面
echo   - 查看每个 Key 的详细信息
echo.
echo 方式 2: 通过 API
echo   curl -X GET "http://localhost:3001/api/keys?group_id=YOUR_GROUP_ID" \
echo     -H "Authorization: Bearer YOUR_AUTH_KEY"
echo.
echo 方式 3: 直接查询数据库（需要 sqlite3）
echo   sqlite3 %DB_PATH% "SELECT id, is_organization_key, organization_id, status FROM api_keys;"
echo.

REM 如果系统有 sqlite3，显示统计
where sqlite3 >nul 2>&1
IF %ERRORLEVEL% EQU 0 (
    echo 检测到 sqlite3，显示统计信息：
    echo ----------------------------------------

    FOR /F "tokens=*" %%i IN ('sqlite3 "%DB_PATH%" "SELECT COUNT(*) FROM api_keys;"') DO SET TOTAL=%%i
    FOR /F "tokens=*" %%i IN ('sqlite3 "%DB_PATH%" "SELECT COUNT(*) FROM api_keys WHERE is_organization_key = 1;"') DO SET ORG_KEYS=%%i
    FOR /F "tokens=*" %%i IN ('sqlite3 "%DB_PATH%" "SELECT COUNT(*) FROM api_keys WHERE is_organization_key = 0;"') DO SET NON_ORG_KEYS=%%i

    echo 总 Key 数量: %TOTAL%
    echo ✅ 组织验证 Key: %ORG_KEYS%
    echo ❌ 非组织 Key: %NON_ORG_KEYS%
    echo.

    IF %ORG_KEYS% GTR 0 (
        echo ✅ 通过组织验证的 Keys:
        sqlite3 -header -column "%DB_PATH%" "SELECT id, organization_id, status, failure_count, request_count FROM api_keys WHERE is_organization_key = 1;"
        echo.
    )

    IF %NON_ORG_KEYS% GTR 0 (
        echo ❌ 未通过组织验证的 Keys:
        sqlite3 -header -column "%DB_PATH%" "SELECT id, status, failure_count, request_count FROM api_keys WHERE is_organization_key = 0;"
        echo.
    )
)

echo 💡 提示：
echo   - 组织验证的 Key 可以使用配置的高级模型
echo   - 非组织 Key 使用高级模型时会在日志中记录警告
echo   - 查看日志：docker logs open-load
pause
