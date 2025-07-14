#!/bin/bash

# End-to-End Test Script for GOX Framework Task 2
# Tests the complete workflow of creating projects and generating components

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TEST_COUNT=0
PASSED_COUNT=0
FAILED_COUNT=0

# Function to print test results
print_test() {
    local test_name="$1"
    local status="$2"
    TEST_COUNT=$((TEST_COUNT + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} Test ${TEST_COUNT}: ${test_name}"
        PASSED_COUNT=$((PASSED_COUNT + 1))
    else
        echo -e "${RED}✗${NC} Test ${TEST_COUNT}: ${test_name}"
        FAILED_COUNT=$((FAILED_COUNT + 1))
    fi
}

# Function to run command and check if it succeeds
run_and_check() {
    local cmd="$1"
    local test_name="$2"
    local expected_to_fail="${3:-false}"
    
    echo -e "${BLUE}Running:${NC} $cmd"
    
    if [ "$expected_to_fail" = "true" ]; then
        if ! $cmd >/dev/null 2>&1; then
            print_test "$test_name" "PASS"
            return 0
        else
            print_test "$test_name (should have failed)" "FAIL"
            return 1
        fi
    else
        if $cmd >/dev/null 2>&1; then
            print_test "$test_name" "PASS"
            return 0
        else
            print_test "$test_name" "FAIL"
            echo -e "${RED}Command failed:${NC} $cmd"
            return 1
        fi
    fi
}

# Function to check if file exists
check_file_exists() {
    local file="$1"
    local test_name="$2"
    
    if [ -f "$file" ]; then
        print_test "$test_name" "PASS"
        return 0
    else
        print_test "$test_name" "FAIL"
        echo -e "${RED}File not found:${NC} $file"
        return 1
    fi
}

# Function to check if directory exists
check_dir_exists() {
    local dir="$1"
    local test_name="$2"
    
    if [ -d "$dir" ]; then
        print_test "$test_name" "PASS"
        return 0
    else
        print_test "$test_name" "FAIL"
        echo -e "${RED}Directory not found:${NC} $dir"
        return 1
    fi
}

echo -e "${BLUE}===========================================${NC}"
echo -e "${BLUE}  GOX Framework E2E Tests - Task 2${NC}"
echo -e "${BLUE}===========================================${NC}"
echo

# Setup: Create temporary directory for tests
TEST_DIR=$(mktemp -d)
echo -e "${YELLOW}Test directory:${NC} $TEST_DIR"
cd "$TEST_DIR"

# Build gox binary
echo -e "${YELLOW}Building gox binary...${NC}"
cd /home/jairo/gox
go build -o "$TEST_DIR/gox" ./cmd/gox/

# Add gox to PATH for this session
export PATH="$TEST_DIR:$PATH"
cd "$TEST_DIR"

echo
echo -e "${BLUE}=== Context Detection Tests ===${NC}"

# Test 1: Generate command should fail outside project
run_and_check "gox generate page test" "Generate command fails outside project" true

# Test 2: New command should succeed outside project
run_and_check "gox new project test-app --router=gin --auth=jwt --db=postgres" "New project command succeeds outside project"

# Test 3: Check if project structure was created
check_dir_exists "test-app" "Project directory created"
check_file_exists "test-app/gox.config.yaml" "Project config file created"
check_file_exists "test-app/go.work" "Go workspace file created"
check_file_exists "test-app/docker-compose.yml" "Docker compose file created"

echo
echo -e "${BLUE}=== Gateway Structure Tests ===${NC}"

# Test 4: Check gateway structure
check_dir_exists "test-app/gateway" "Gateway directory created"
check_file_exists "test-app/gateway/main.go" "Gateway main.go created"
check_file_exists "test-app/gateway/go.mod" "Gateway go.mod created"
check_dir_exists "test-app/gateway/pages" "Gateway pages directory created"
check_dir_exists "test-app/gateway/components" "Gateway components directory created"
check_dir_exists "test-app/gateway/shared/ui" "Gateway shared/ui directory created"
check_dir_exists "test-app/gateway/shared/layouts" "Gateway shared/layouts directory created"

# Test 5: Check common structure
check_dir_exists "test-app/common" "Common directory created"
check_dir_exists "test-app/common/middleware" "Common middleware directory created"
check_dir_exists "test-app/common/discovery" "Common discovery directory created"
check_file_exists "test-app/common/middleware/service_router.go" "Service router middleware created"
check_file_exists "test-app/common/discovery/interface.go" "Discovery interface created"

# Test 6: Check services structure
check_dir_exists "test-app/services" "Services directory created"

echo
echo -e "${BLUE}=== Inside Project Tests ===${NC}"

# Change to project directory
cd test-app

# Test 7: New command should fail inside project
run_and_check "gox new project another-app" "New project command fails inside project" true

# Test 8: Generate page command
run_and_check "gox generate page dashboard --auth" "Generate page with auth"
check_file_exists "gateway/pages/dashboard.gox" "Dashboard page file created"

# Test 9: Generate regular component
run_and_check "gox generate component user-card --props=\"name:string,email:string,featured:bool\"" "Generate component with props"
check_file_exists "gateway/components/user-card.gox" "User card component created"

# Test 10: Generate shared component
run_and_check "gox generate component button --shared" "Generate shared component"
check_file_exists "gateway/shared/ui/button.gox" "Shared button component created"

# Test 11: Generate middleware
run_and_check "gox generate middleware rate-limit" "Generate middleware"
check_file_exists "middleware/rate-limit.go" "Rate limit middleware created"

# Test 12: Generate service (microservice)
run_and_check "gox generate service users --api" "Generate service with API"
check_dir_exists "services/users" "Users service directory created"
check_file_exists "services/users/internal/service/service.go" "Service implementation created"
check_file_exists "services/users/internal/handlers/handlers.go" "Service handlers created"
check_file_exists "services/users/cmd/server/main.go" "Service main.go created"

echo
echo -e "${BLUE}=== Content Validation Tests ===${NC}"

# Test 13: Check that generated files contain expected content
if grep -q "Dashboard" "gateway/pages/dashboard.gox"; then
    print_test "Dashboard page contains expected content" "PASS"
else
    print_test "Dashboard page contains expected content" "FAIL"
fi

# Test 14: Check component with props
if grep -q "Name string\|Email string\|Featured bool" "gateway/components/user-card.gox"; then
    print_test "Component contains props" "PASS"
else
    print_test "Component contains props" "FAIL"
fi

# Test 15: Check auth in page
if grep -q "auth" "gateway/pages/dashboard.gox"; then
    print_test "Page contains auth code" "PASS"
else
    print_test "Page contains auth code" "FAIL"
fi

# Test 16: Check gateway main.go has Gin imports (since we used --router=gin)
if grep -q "gin" "gateway/main.go"; then
    print_test "Gateway main.go uses Gin router" "PASS"
else
    print_test "Gateway main.go uses Gin router" "FAIL"
fi

# Test 17: Check docker-compose contains expected services
if grep -q "consul:" "$TEST_DIR/test-app/docker-compose.yml"; then
    print_test "Docker compose includes Consul" "PASS"
else
    print_test "Docker compose includes Consul" "FAIL"
fi

if grep -q "postgres:" "$TEST_DIR/test-app/docker-compose.yml"; then
    print_test "Docker compose includes Postgres" "PASS"
else
    print_test "Docker compose includes Postgres" "FAIL"
fi

echo
echo -e "${BLUE}=== Service Creation Test ===${NC}"

# Change back to test directory
cd "$TEST_DIR"

# Test 18: Create standalone service outside project
run_and_check "gox new service payment-service --db=postgres" "Create standalone service"
check_dir_exists "payment-service" "Payment service directory created"
check_file_exists "payment-service/cmd/server/main.go" "Service main.go created"
check_file_exists "payment-service/go.mod" "Service go.mod created"

echo
echo -e "${BLUE}===========================================${NC}"
echo -e "${BLUE}  Test Results Summary${NC}"
echo -e "${BLUE}===========================================${NC}"
echo -e "${GREEN}Passed: ${PASSED_COUNT}${NC}"
echo -e "${RED}Failed: ${FAILED_COUNT}${NC}"
echo -e "${BLUE}Total:  ${TEST_COUNT}${NC}"
echo

if [ $FAILED_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! Task 2 implementation is complete.${NC}"
    exit_code=0
else
    echo -e "${RED}❌ Some tests failed. Please check the implementation.${NC}"
    exit_code=1
fi

# Cleanup
echo -e "${YELLOW}Cleaning up test directory: $TEST_DIR${NC}"
rm -rf "$TEST_DIR"

exit $exit_code