#!/bin/bash

# 🧪 Script de Pruebas Automatizadas - API Andrei Mes Manur
# Este script ejecuta todas las pruebas de endpoints y verifica permisos

API_URL="http://localhost:8080/api/v1"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Contadores
PASSED=0
FAILED=0
TOTAL=0

# Función para imprimir resultados
print_result() {
    TOTAL=$((TOTAL + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED:${NC} $2"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}❌ FAILED:${NC} $2"
        FAILED=$((FAILED + 1))
    fi
}

# Función para hacer requests y verificar status code
test_endpoint() {
    local method=$1
    local url=$2
    local expected_status=$3
    local description=$4
    local headers=$5
    local data=$6
    
    if [ -n "$headers" ] && [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X $method "$API_URL$url" -H "Content-Type: application/json" -H "$headers" -d "$data")
    elif [ -n "$headers" ]; then
        response=$(curl -s -w "%{http_code}" -X $method "$API_URL$url" -H "$headers")
    elif [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X $method "$API_URL$url" -H "Content-Type: application/json" -d "$data")
    else
        response=$(curl -s -w "%{http_code}" -X $method "$API_URL$url")
    fi
    
    status_code="${response: -3}"
    
    if [ "$status_code" = "$expected_status" ]; then
        print_result 0 "$description (Status: $status_code)"
        return 0
    else
        print_result 1 "$description (Expected: $expected_status, Got: $status_code)"
        return 1
    fi
}

echo -e "${BLUE}🚀 Iniciando pruebas de la API Andrei Mes Manur${NC}\n"

# FASE 1: AUTENTICACIÓN
echo -e "${YELLOW}📋 FASE 1: Pruebas de Autenticación${NC}"

# Login Andrei
echo -e "\n🔐 Autenticando usuarios..."
ANDREI_RESPONSE=$(curl -s -X POST "$API_URL/login" -H "Content-Type: application/json" -d '{"email":"andrei@evil.com","password":"password123"}')
ANDREI_TOKEN=$(echo $ANDREI_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$ANDREI_TOKEN" ]; then
    print_result 0 "Login Andrei exitoso"
    ANDREI_AUTH="Authorization: Bearer $ANDREI_TOKEN"
else
    print_result 1 "Login Andrei falló"
    ANDREI_AUTH=""
fi

# Login Demonio
DEMON_RESPONSE=$(curl -s -X POST "$API_URL/login" -H "Content-Type: application/json" -d '{"email":"shadow@evil.com","password":"demon123"}')
DEMON_TOKEN=$(echo $DEMON_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$DEMON_TOKEN" ]; then
    print_result 0 "Login Demonio exitoso"
    DEMON_AUTH="Authorization: Bearer $DEMON_TOKEN"
else
    print_result 1 "Login Demonio falló"
    DEMON_AUTH=""
fi

# Login Network Admin
ADMIN_RESPONSE=$(curl -s -X POST "$API_URL/login" -H "Content-Type: application/json" -d '{"email":"john.admin@company.com","password":"admin123"}')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$ADMIN_TOKEN" ]; then
    print_result 0 "Login Network Admin exitoso"
    ADMIN_AUTH="Authorization: Bearer $ADMIN_TOKEN"
else
    print_result 1 "Login Network Admin falló"
    ADMIN_AUTH=""
fi

# Pruebas de registro
test_endpoint "POST" "/register" "201" "Registro de nuevo demonio" "" '{"username":"TestDemon2","email":"testdemon2@evil.com","password":"password123","role":"demon"}'
test_endpoint "POST" "/register" "201" "Registro de nuevo network admin" "" '{"username":"TestAdmin2","email":"testadmin2@company.com","password":"password123","role":"network_admin"}'
test_endpoint "POST" "/register" "400" "Registro como Andrei debe fallar" "" '{"username":"FakeAndrei","email":"fake@evil.com","password":"password123","role":"andrei"}'

# FASE 2: ENDPOINTS PÚBLICOS
echo -e "\n${YELLOW}📋 FASE 2: Endpoints Públicos${NC}"
test_endpoint "GET" "/resistance" "200" "Acceso público a página de resistencia"

# FASE 3: FUNCIONALIDADES DE ANDREI
echo -e "\n${YELLOW}📋 FASE 3: Funcionalidades de Andrei${NC}"
if [ -n "$ANDREI_AUTH" ]; then
    test_endpoint "GET" "/admin/users" "200" "Ver todos los usuarios" "$ANDREI_AUTH"
    test_endpoint "GET" "/admin/users/2" "200" "Ver usuario específico" "$ANDREI_AUTH"
    test_endpoint "GET" "/admin/stats" "200" "Ver estadísticas de plataforma" "$ANDREI_AUTH"
    test_endpoint "GET" "/admin/demons/ranking" "200" "Ver ranking de demonios" "$ANDREI_AUTH"
    test_endpoint "GET" "/admin/posts" "200" "Ver todos los posts" "$ANDREI_AUTH"
    test_endpoint "POST" "/admin/posts" "201" "Crear post como Andrei" "$ANDREI_AUTH" '{"title":"Test Post Andrei","body":"Test content","media":""}'
    test_endpoint "POST" "/admin/rewards" "201" "Crear recompensa" "$ANDREI_AUTH" '{"demon_id":2,"type":"reward","title":"Test Reward","description":"Test reward","points":100}'
    test_endpoint "POST" "/admin/rewards" "201" "Crear castigo" "$ANDREI_AUTH" '{"demon_id":2,"type":"punishment","title":"Test Punishment","description":"Test punishment","points":-50}'
fi

# FASE 4: FUNCIONALIDADES DE DEMONIO
echo -e "\n${YELLOW}📋 FASE 4: Funcionalidades de Demonio${NC}"
if [ -n "$DEMON_AUTH" ]; then
    test_endpoint "POST" "/demons/victims" "201" "Registrar nueva víctima" "$DEMON_AUTH" '{"username":"TestVictim","email":"testvictim@company.com","password":"victim123","role":"network_admin"}'
    test_endpoint "GET" "/demons/victims" "200" "Ver mis víctimas" "$DEMON_AUTH"
    test_endpoint "POST" "/demons/reports" "201" "Crear reporte" "$DEMON_AUTH" '{"victim_id":9,"title":"Test Report","description":"Test report description"}'
    test_endpoint "GET" "/demons/reports" "200" "Ver mis reportes" "$DEMON_AUTH"
    test_endpoint "GET" "/demons/stats" "200" "Ver mis estadísticas" "$DEMON_AUTH"
    test_endpoint "POST" "/demons/posts" "201" "Crear post como demonio" "$DEMON_AUTH" '{"title":"Test Post Demon","body":"Test content from demon","media":""}'
fi

# FASE 5: FUNCIONALIDADES DE NETWORK ADMIN
echo -e "\n${YELLOW}📋 FASE 5: Funcionalidades de Network Admin${NC}"
if [ -n "$ADMIN_AUTH" ]; then
    test_endpoint "POST" "/network-admins/posts/anonymous" "201" "Crear post anónimo" "$ADMIN_AUTH" '{"title":"Test Anonymous Post","body":"Test anonymous content","media":""}'
fi

# FASE 6: PRUEBAS DE SEGURIDAD Y PERMISOS
echo -e "\n${YELLOW}📋 FASE 6: Pruebas de Seguridad${NC}"

# Acceso sin token
test_endpoint "GET" "/admin/users" "401" "Acceso sin token debe fallar"
test_endpoint "GET" "/demons/stats" "401" "Acceso sin token debe fallar"

# Token inválido
test_endpoint "GET" "/admin/users" "401" "Token inválido debe fallar" "Authorization: Bearer token_invalido"

# Permisos cruzados (deben fallar)
if [ -n "$DEMON_AUTH" ]; then
    test_endpoint "GET" "/admin/users" "403" "Demonio accediendo a funciones de Andrei debe fallar" "$DEMON_AUTH"
    test_endpoint "DELETE" "/admin/users/5" "403" "Demonio eliminando usuario debe fallar" "$DEMON_AUTH"
    test_endpoint "GET" "/admin/stats" "403" "Demonio viendo estadísticas generales debe fallar" "$DEMON_AUTH"
fi

if [ -n "$ADMIN_AUTH" ]; then
    test_endpoint "POST" "/demons/reports" "403" "Network Admin creando reporte debe fallar" "$ADMIN_AUTH" '{"victim_id":1,"title":"Test","description":"Test"}'
    test_endpoint "GET" "/demons/stats" "403" "Network Admin viendo estadísticas de demonio debe fallar" "$ADMIN_AUTH"
    test_endpoint "GET" "/admin/users" "403" "Network Admin accediendo a funciones de Andrei debe fallar" "$ADMIN_AUTH"
fi

# Verificación final
echo -e "\n${YELLOW}📋 VERIFICACIÓN FINAL${NC}"
test_endpoint "GET" "/resistance" "200" "Página de resistencia actualizada"

# Resumen de resultados
echo -e "\n${BLUE}📊 RESUMEN DE PRUEBAS${NC}"
echo -e "Total de pruebas: $TOTAL"
echo -e "${GREEN}Exitosas: $PASSED${NC}"
echo -e "${RED}Fallidas: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ¡Todas las pruebas pasaron! La API está funcionando correctamente.${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️ Algunas pruebas fallaron. Revisar los errores arriba.${NC}"
    exit 1
fi