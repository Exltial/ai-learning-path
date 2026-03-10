#!/bin/bash

# ============================================
# AI 学习平台 - API 自动化测试脚本
# ============================================
# 使用方法:
#   ./test-api.sh [base_url]
# 示例:
#   ./test-api.sh http://localhost:8080
# ============================================

set -e

# 配置
BASE_URL="${1:-http://localhost:8080}"
API_PREFIX="/api/v1"
TEST_USER="testuser_$(date +%s)"
TEST_EMAIL="test_${TEST_USER}@example.com"
TEST_PASSWORD="Password123"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果数组
declare -a TEST_RESULTS

# 打印分隔线
print_separator() {
    echo -e "${BLUE}============================================${NC}"
}

# 打印测试结果
print_result() {
    local test_name="$1"
    local status="$2"
    local expected="$3"
    local actual="$4"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" == "PASS" ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo -e "${GREEN}✓ PASS${NC}: $test_name"
        TEST_RESULTS+=("{\"name\":\"$test_name\",\"status\":\"passed\",\"expected\":\"$expected\",\"actual\":\"$actual\"}")
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo -e "${RED}✗ FAIL${NC}: $test_name"
        echo -e "  期望：$expected"
        echo -e "  实际：$actual"
        TEST_RESULTS+=("{\"name\":\"$test_name\",\"status\":\"failed\",\"expected\":\"$expected\",\"actual\":\"$actual\"}")
    fi
}

# 检查服务是否可用
check_service() {
    echo -e "${YELLOW}检查服务可用性...${NC}"
    if curl -s -o /dev/null -w "%{http_code}" "$BASE_URL" | grep -q "200\|404"; then
        echo -e "${GREEN}✓ 服务可访问：$BASE_URL${NC}"
        return 0
    else
        echo -e "${RED}✗ 服务不可访问：$BASE_URL${NC}"
        echo "请确保后端服务已启动"
        return 1
    fi
}

# ============================================
# 认证模块测试
# ============================================
test_authentication() {
    print_separator
    echo -e "${BLUE}认证模块测试${NC}"
    print_separator
    
    local token=""
    
    # 测试 1: 用户注册 - 正常流程
    echo -e "\n${YELLOW}测试：用户注册 - 正常流程${NC}"
    local register_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$API_PREFIX/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USER\",
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    local register_status=$(echo "$register_response" | tail -n1)
    local register_body=$(echo "$register_response" | head -n-1)
    
    if [ "$register_status" == "201" ]; then
        print_result "用户注册 - 正常流程" "PASS" "201" "$register_status"
        token=$(echo "$register_body" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    else
        print_result "用户注册 - 正常流程" "FAIL" "201" "$register_status"
    fi
    
    # 测试 2: 用户注册 - 用户名重复
    echo -e "\n${YELLOW}测试：用户注册 - 用户名重复${NC}"
    local duplicate_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$API_PREFIX/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$TEST_USER\",
            \"email\": \"another@example.com\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    local duplicate_status=$(echo "$duplicate_response" | tail -n1)
    
    if [ "$duplicate_status" == "409" ]; then
        print_result "用户注册 - 用户名重复" "PASS" "409" "$duplicate_status"
    else
        print_result "用户注册 - 用户名重复" "FAIL" "409" "$duplicate_status"
    fi
    
    # 测试 3: 用户注册 - 邮箱重复
    echo -e "\n${YELLOW}测试：用户注册 - 邮箱重复${NC}"
    local email_duplicate_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$API_PREFIX/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"another_user\",
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    local email_duplicate_status=$(echo "$email_duplicate_response" | tail -n1)
    
    if [ "$email_duplicate_status" == "409" ]; then
        print_result "用户注册 - 邮箱重复" "PASS" "409" "$email_duplicate_status"
    else
        print_result "用户注册 - 邮箱重复" "FAIL" "409" "$email_duplicate_status"
    fi
    
    # 测试 4: 用户注册 - 密码强度不足
    echo -e "\n${YELLOW}测试：用户注册 - 密码强度不足${NC}"
    local weak_pwd_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$API_PREFIX/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"weakpwd_user\",
            \"email\": \"weakpwd@example.com\",
            \"password\": \"123\"
        }")
    local weak_pwd_status=$(echo "$weak_pwd_response" | tail -n1)
    
    if [ "$weak_pwd_status" == "400" ]; then
        print_result "用户注册 - 密码强度不足" "PASS" "400" "$weak_pwd_status"
    else
        print_result "用户注册 - 密码强度不足" "FAIL" "400" "$weak_pwd_status"
    fi
    
    # 测试 5: 用户登录 - 正常流程
    echo -e "\n${YELLOW}测试：用户登录 - 正常流程${NC}"
    local login_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$API_PREFIX/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"$TEST_PASSWORD\"
        }")
    local login_status=$(echo "$login_response" | tail -n1)
    local login_body=$(echo "$login_response" | head -n-1)
    
    if [ "$login_status" == "200" ]; then
        print_result "用户登录 - 正常流程" "PASS" "200" "$login_status"
        token=$(echo "$login_body" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    else
        print_result "用户登录 - 正常流程" "FAIL" "200" "$login_status"
    fi
    
    # 测试 6: 用户登录 - 错误密码
    echo -e "\n${YELLOW}测试：用户登录 - 错误密码${NC}"
    local wrong_pwd_response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL$API_PREFIX/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$TEST_EMAIL\",
            \"password\": \"WrongPassword\"
        }")
    local wrong_pwd_status=$(echo "$wrong_pwd_response" | tail -n1)
    
    if [ "$wrong_pwd_status" == "401" ]; then
        print_result "用户登录 - 错误密码" "PASS" "401" "$wrong_pwd_status"
    else
        print_result "用户登录 - 错误密码" "FAIL" "401" "$wrong_pwd_status"
    fi
    
    # 测试 7: 获取当前用户信息
    echo -e "\n${YELLOW}测试：获取当前用户信息${NC}"
    if [ -n "$token" ]; then
        local user_response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$API_PREFIX/users/me" \
            -H "Authorization: Bearer $token")
        local user_status=$(echo "$user_response" | tail -n1)
        
        if [ "$user_status" == "200" ]; then
            print_result "获取当前用户信息" "PASS" "200" "$user_status"
        else
            print_result "获取当前用户信息" "FAIL" "200" "$user_status"
        fi
    else
        print_result "获取当前用户信息" "FAIL" "200" "No token"
    fi
    
    echo "$token"
}

# ============================================
# 课程模块测试
# ============================================
test_courses() {
    local token="$1"
    
    print_separator
    echo -e "${BLUE}课程模块测试${NC}"
    print_separator
    
    # 测试 1: 获取课程列表
    echo -e "\n${YELLOW}测试：获取课程列表${NC}"
    local courses_response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$API_PREFIX/courses")
    local courses_status=$(echo "$courses_response" | tail -n1)
    
    if [ "$courses_status" == "200" ]; then
        print_result "获取课程列表" "PASS" "200" "$courses_status"
    else
        print_result "获取课程列表" "FAIL" "200" "$courses_status"
    fi
    
    # 测试 2: 获取课程列表 - 分页
    echo -e "\n${YELLOW}测试：获取课程列表 - 分页${NC}"
    local paginated_response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$API_PREFIX/courses?page=1&limit=10")
    local paginated_status=$(echo "$paginated_response" | tail -n1)
    
    if [ "$paginated_status" == "200" ]; then
        print_result "获取课程列表 - 分页" "PASS" "200" "$paginated_status"
    else
        print_result "获取课程列表 - 分页" "FAIL" "200" "$paginated_status"
    fi
    
    # 测试 3: 获取课程列表 - 难度筛选
    echo -e "\n${YELLOW}测试：获取课程列表 - 难度筛选${NC}"
    local filtered_response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$API_PREFIX/courses?difficulty=beginner")
    local filtered_status=$(echo "$filtered_response" | tail -n1)
    
    if [ "$filtered_status" == "200" ]; then
        print_result "获取课程列表 - 难度筛选" "PASS" "200" "$filtered_status"
    else
        print_result "获取课程列表 - 难度筛选" "FAIL" "200" "$filtered_status"
    fi
    
    # 测试 4: 获取课程详情（需要有效的课程 ID）
    echo -e "\n${YELLOW}测试：获取课程详情${NC}"
    # 先从课程列表获取一个课程 ID
    local courses_body=$(curl -s -X GET "$BASE_URL$API_PREFIX/courses")
    local course_id=$(echo "$courses_body" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [ -n "$course_id" ]; then
        local detail_response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$API_PREFIX/courses/$course_id")
        local detail_status=$(echo "$detail_response" | tail -n1)
        
        if [ "$detail_status" == "200" ]; then
            print_result "获取课程详情" "PASS" "200" "$detail_status"
        else
            print_result "获取课程详情" "FAIL" "200" "$detail_status"
        fi
    else
        print_result "获取课程详情" "FAIL" "200" "No course found"
    fi
    
    # 测试 5: 获取不存在的课程详情
    echo -e "\n${YELLOW}测试：获取不存在的课程详情${NC}"
    local notfound_response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$API_PREFIX/courses/00000000-0000-0000-0000-000000000000")
    local notfound_status=$(echo "$notfound_response" | tail -n1)
    
    if [ "$notfound_status" == "404" ]; then
        print_result "获取不存在的课程详情" "PASS" "404" "$notfound_status"
    else
        print_result "获取不存在的课程详情" "FAIL" "404" "$notfound_status"
    fi
}

# ============================================
# 性能测试
# ============================================
test_performance() {
    print_separator
    echo -e "${BLUE}性能测试（并发请求）${NC}"
    print_separator
    
    # 测试 1: 课程列表并发请求
    echo -e "\n${YELLOW}测试：课程列表并发请求 (10 个并发)${NC}"
    local start_time=$(date +%s%N)
    
    for i in {1..10}; do
        curl -s -X GET "$BASE_URL$API_PREFIX/courses" > /dev/null &
    done
    wait
    
    local end_time=$(date +%s%N)
    local duration=$(( (end_time - start_time) / 1000000 ))
    
    if [ $duration -lt 5000 ]; then
        print_result "课程列表并发请求" "PASS" "<5000ms" "${duration}ms"
    else
        print_result "课程列表并发请求" "FAIL" "<5000ms" "${duration}ms"
    fi
}

# ============================================
# 主函数
# ============================================
main() {
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE}   AI 学习平台 - API 自动化测试${NC}"
    echo -e "${BLUE}============================================${NC}"
    echo ""
    echo "基础 URL: $BASE_URL"
    echo "测试用户：$TEST_USER"
    echo ""
    
    # 检查服务
    if ! check_service; then
        echo ""
        echo -e "${RED}测试中止：服务不可用${NC}"
        echo ""
        echo "请启动后端服务:"
        echo "  cd /home/admin/.openclaw/workspace/projects/ai-learning-platform/backend"
        echo "  docker-compose up -d"
        echo "  go run cmd/main.go"
        exit 1
    fi
    
    echo ""
    
    # 执行测试
    token=$(test_authentication)
    test_courses "$token"
    test_performance
    
    # 输出统计
    print_separator
    echo -e "${BLUE}测试统计${NC}"
    print_separator
    echo "总测试数：$TOTAL_TESTS"
    echo -e "${GREEN}通过：$PASSED_TESTS${NC}"
    echo -e "${RED}失败：$FAILED_TESTS${NC}"
    
    if [ $TOTAL_TESTS -gt 0 ]; then
        local pass_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
        echo "通过率：$pass_rate%"
    fi
    
    # 生成 JSON 结果
    echo ""
    echo -e "${YELLOW}生成测试结果 JSON...${NC}"
    generate_json_report
    
    print_separator
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "${GREEN}✓ 所有测试通过！${NC}"
        exit 0
    else
        echo -e "${RED}✗ 部分测试失败，请检查日志${NC}"
        exit 1
    fi
}

# 生成 JSON 报告
generate_json_report() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local json_file="/home/admin/.openclaw/workspace/projects/ai-learning-platform/docs/testing/api-test-results.json"
    
    cat > "$json_file" << EOF
{
  "test_run": {
    "timestamp": "$timestamp",
    "base_url": "$BASE_URL",
    "total_tests": $TOTAL_TESTS,
    "passed": $PASSED_TESTS,
    "failed": $FAILED_TESTS,
    "pass_rate": $((TOTAL_TESTS > 0 ? PASSED_TESTS * 100 / TOTAL_TESTS : 0))
  },
  "results": [
    $(IFS=,; echo "${TEST_RESULTS[*]}")
  ]
}
EOF
    
    echo "JSON 报告已生成：$json_file"
}

# 运行主函数
main
