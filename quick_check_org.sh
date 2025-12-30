#!/bin/bash
# 快速查看组织验证状态

# API 端点配置
HOST="${1:-localhost:3001}"
AUTH_KEY="${2:-your-auth-key}"

echo "🔍 查询 API Keys 组织验证状态"
echo "================================"
echo ""
echo "如果没有提供参数，请手动修改此脚本中的 AUTH_KEY"
echo ""

# 获取所有 groups
echo "📋 获取所有 Groups..."
GROUPS=$(curl -s "http://$HOST/api/groups" -H "Authorization: Bearer $AUTH_KEY")

# 解析第一个 group_id
GROUP_ID=$(echo "$GROUPS" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$GROUP_ID" ]; then
    echo "❌ 无法获取 Group ID，请检查："
    echo "  1. 服务是否运行：http://$HOST"
    echo "  2. AUTH_KEY 是否正确"
    exit 1
fi

echo "✅ 找到 Group ID: $GROUP_ID"
echo ""

# 获取该 group 的所有 keys
echo "📊 查询 Keys..."
KEYS=$(curl -s "http://$HOST/api/keys?group_id=$GROUP_ID" -H "Authorization: Bearer $AUTH_KEY")

# 统计和显示
echo "$KEYS" | python3 -c "
import json
import sys

try:
    data = json.load(sys.stdin)
    items = data.get('items', [])

    if not items:
        print('❌ 没有找到任何 Key')
        sys.exit(1)

    org_keys = [k for k in items if k.get('is_organization_key')]
    non_org_keys = [k for k in items if not k.get('is_organization_key')]

    print(f'总计: {len(items)} 个 Keys')
    print(f'✅ 组织验证: {len(org_keys)} 个')
    print(f'❌ 非组织: {len(non_org_keys)} 个')
    print()

    if org_keys:
        print('✅ 通过组织验证的 Keys:')
        print('-' * 80)
        print(f'{'ID':<6} {'Org ID':<25} {'Status':<10} {'Requests':<10} {'Failures':<10}')
        print('-' * 80)
        for key in org_keys:
            print(f'{key['id']:<6} {key.get('organization_id', ''):<25} {key['status']:<10} {key['request_count']:<10} {key['failure_count']:<10}')
        print()

    if non_org_keys:
        print('❌ 未通过组织验证的 Keys:')
        print('-' * 60)
        print(f'{'ID':<6} {'Status':<10} {'Requests':<10} {'Failures':<10}')
        print('-' * 60)
        for key in non_org_keys:
            print(f'{key['id']:<6} {key['status']:<10} {key['request_count']:<10} {key['failure_count']:<10}')
        print()

except Exception as e:
    print(f'❌ 解析失败: {e}')
    print('原始数据:')
    print(sys.stdin.read())
"

echo ""
echo "💡 提示："
echo "  - 组织验证的 Key 可以使用所有模型"
echo "  - 非组织 Key 使用高级模型时会记录警告日志"
echo "  - 查看日志：docker logs -f open-load | grep organization"
