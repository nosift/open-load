#!/bin/bash
# 查看 API Keys 的组织验证状态

DB_PATH="./data/open-load.db"

# 如果 open-load.db 不存在，尝试 gpt-load.db
if [ ! -f "$DB_PATH" ]; then
    DB_PATH="./data/gpt-load.db"
fi

if [ ! -f "$DB_PATH" ]; then
    echo "❌ 数据库文件不存在"
    exit 1
fi

echo "📊 API Keys 组织验证状态统计"
echo "================================"
echo ""

# 统计总数
TOTAL=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM api_keys;")
ORG_KEYS=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM api_keys WHERE is_organization_key = 1;")
NON_ORG_KEYS=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM api_keys WHERE is_organization_key = 0;")

echo "总 Key 数量: $TOTAL"
echo "✅ 组织验证 Key: $ORG_KEYS"
echo "❌ 非组织 Key: $NON_ORG_KEYS"
echo ""

# 显示组织验证的 Key
if [ $ORG_KEYS -gt 0 ]; then
    echo "✅ 通过组织验证的 Keys:"
    echo "----------------------------------------"
    sqlite3 -header -column "$DB_PATH" "
        SELECT
            id as ID,
            substr(organization_id, 1, 20) as 'Organization ID',
            status as Status,
            failure_count as Failures,
            request_count as Requests
        FROM api_keys
        WHERE is_organization_key = 1
        ORDER BY id;
    "
    echo ""
fi

# 显示非组织验证的 Key
if [ $NON_ORG_KEYS -gt 0 ]; then
    echo "❌ 未通过组织验证的 Keys:"
    echo "----------------------------------------"
    sqlite3 -header -column "$DB_PATH" "
        SELECT
            id as ID,
            status as Status,
            failure_count as Failures,
            request_count as Requests
        FROM api_keys
        WHERE is_organization_key = 0
        ORDER BY id;
    "
    echo ""
fi

# 显示详细信息（可选）
if [ "$1" == "-v" ] || [ "$1" == "--verbose" ]; then
    echo "📋 详细信息:"
    echo "========================================"
    sqlite3 -header -column "$DB_PATH" "
        SELECT
            id as ID,
            CASE WHEN is_organization_key = 1 THEN '✅ Yes' ELSE '❌ No' END as 'Org Key',
            substr(organization_id, 1, 15) as 'Org ID',
            status as Status,
            failure_count as Fails,
            request_count as Reqs
        FROM api_keys
        ORDER BY is_organization_key DESC, id;
    "
fi

echo ""
echo "💡 提示："
echo "  - 使用 -v 或 --verbose 查看详细信息"
echo "  - 组织验证的 Key 可以使用配置的高级模型"
echo "  - 非组织 Key 使用高级模型时会在日志中记录警告"
