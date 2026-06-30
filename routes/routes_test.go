package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"hris-payroll/config"
	"hris-payroll/domain"
	authDelivery "hris-payroll/internal/auth/delivery"
	authUsecase "hris-payroll/internal/auth/usecase"
	leaveDelivery "hris-payroll/internal/leave/delivery"
	leaveRepository "hris-payroll/internal/leave/repository"
	leaveUsecase "hris-payroll/internal/leave/usecase"
	payrollDelivery "hris-payroll/internal/payroll/delivery"
	payrollRepository "hris-payroll/internal/payroll/repository"
	payrollUsecase "hris-payroll/internal/payroll/usecase"
	"hris-payroll/routes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper to generate a test JWT token
func generateTestToken(userID uint, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "juara-coding-super-secret-key-go-hris-payroll-2026-batch-1"
		os.Setenv("JWT_SECRET", secretKey)
	}
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_role": role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// Setup a fully configured testing Gin engine within an active DB transaction
func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)

	// Load environment variables from parent directory's .env file
	_ = godotenv.Load("../.env")

	// Establish connection to database
	db, err := config.ConnectDB()
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Begin a transaction so test data is never permanently persisted
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin database transaction: %v", tx.Error)
	}

	// Run migration on transaction
	err = tx.AutoMigrate(
		&domain.DepartmentBudget{},
		&domain.Employee{},
		&domain.Payroll{},
		&domain.LeaveRequest{},
		&domain.TokenBlacklist{},
	)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to migrate tables in transaction: %v", err)
	}

	// Instantiate layers with the transaction object
	aUsecase := authUsecase.NewAuthUsecase(tx)
	aHandler := authDelivery.NewAuthHandler(aUsecase)

	pRepo := payrollRepository.NewPayrollRepository(tx)
	pUsecase := payrollUsecase.NewPayrollUsecase(pRepo)
	pHandler := payrollDelivery.NewPayrollHandler(pUsecase)

	lRepo := leaveRepository.NewLeaveRepository(tx)
	lUsecase := leaveUsecase.NewLeaveUsecase(lRepo)
	lHandler := leaveDelivery.NewLeaveHandler(lUsecase)

	r := gin.Default()
	routes.SetupRoutes(r, tx, aHandler, pHandler, lHandler)

	// Teardown function to roll back transaction
	teardown := func() {
		tx.Rollback()
	}

	return r, tx, teardown
}

// Seed helper for test suite
func seedTestData(t *testing.T, db *gorm.DB) (uint, uint, string) {
	// Seed a test department
	dept := domain.DepartmentBudget{
		ID:         1000,
		Name:       "Test Engineering",
		BudgetLeft: 500000000,
	}
	if err := db.Create(&dept).Error; err != nil {
		t.Fatalf("failed to seed test department: %v", err)
	}

	// Hash password using Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Secret123!"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Create test HRD
	hrd := domain.Employee{
		ID:           1001,
		Name:         "Test HRD Admin",
		Email:        "hrd-test@company.co.id",
		Password:     string(hashedPassword),
		Role:         "HRD",
		BaseSalary:   12000000,
		DepartmentID: dept.ID,
	}
	if err := db.Create(&hrd).Error; err != nil {
		t.Fatalf("failed to seed test HRD: %v", err)
	}

	// Generate JWT Token for HRD
	token, err := generateTestToken(hrd.ID, hrd.Role)
	if err != nil {
		t.Fatalf("failed to generate test JWT token: %v", err)
	}

	return dept.ID, hrd.ID, token
}

// Helper to create a request with body and auth header
func performRequest(r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestIntegration_RBAC_And_PostEmployeeData(t *testing.T) {
	router, db, teardown := setupTestRouter(t)
	defer teardown()

	deptID, _, hrdToken := seedTestData(t, db)

	// Create a standard Employee token to test access control restrictions
	empID := uint(2001)
	empToken, err := generateTestToken(empID, "EMPLOYEE")
	if err != nil {
		t.Fatalf("failed to generate employee token: %v", err)
	}

	// 1. RBAC Tests (Authorization RBAC)
	t.Run("RBAC - Positive - HRD can register employee", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Valid Employee",
			Email:        "valid-emp@company.co.id",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, hrdToken)
		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("RBAC - Negative - Employee cannot register employee", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Should Fail",
			Email:        "fail-emp@company.co.id",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, empToken)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("RBAC - Negative - Unauthenticated user cannot register employee", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Should Fail",
			Email:        "fail-emp2@company.co.id",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("RBAC - Negative - Employee cannot access process payroll", func(t *testing.T) {
		payload := []map[string]interface{}{
			{
				"employee_id": 1,
				"bonus":       100000,
			},
		}
		w := performRequest(router, "POST", "/api/v1/payroll/process", payload, empToken)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	// 2. Post Data Employees Validation Tests
	t.Run("Post Employee - Negative - Email domain must be @company.co.id", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Gmail User",
			Email:        "gmail-user@gmail.com",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, hrdToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 Bad Request, got %d. Response: %s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "email harus menggunakan domain resmi @company.co.id") {
			t.Errorf("expected email domain error, got: %s", w.Body.String())
		}
	})

	t.Run("Post Employee - Negative - Duplicate email", func(t *testing.T) {
		// First registration
		payload1 := domain.RegisterReq{
			Name:         "Duplicate Emp",
			Email:        "duplicate@company.co.id",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w1 := performRequest(router, "POST", "/api/v1/auth/register", payload1, hrdToken)
		if w1.Code != http.StatusCreated {
			t.Fatalf("first registration should succeed, got %d. Response: %s", w1.Code, w1.Body.String())
		}

		// Second registration with the same email
		payload2 := domain.RegisterReq{
			Name:         "Duplicate Emp 2",
			Email:        "duplicate@company.co.id",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w2 := performRequest(router, "POST", "/api/v1/auth/register", payload2, hrdToken)
		if w2.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 Bad Request on duplicate, got %d. Response: %s", w2.Code, w2.Body.String())
		}
		if !strings.Contains(w2.Body.String(), "email sudah terdaftar di sistem") {
			t.Errorf("expected email duplicate error, got: %s", w2.Body.String())
		}
	})

	t.Run("Post Employee - Negative - Invalid fields (short password)", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Short Pwd User",
			Email:        "short-pwd@company.co.id",
			Password:     "123", // too short, min is 6
			Role:         "EMPLOYEE",
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, hrdToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 Bad Request, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Post Employee - Negative - Invalid fields (negative salary)", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Neg Salary User",
			Email:        "neg-salary@company.co.id",
			Password:     "Password123!",
			Role:         "EMPLOYEE",
			BaseSalary:   -100.0, // invalid, must be gt=0
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, hrdToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 Bad Request, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Post Employee - Negative - Invalid fields (invalid role)", func(t *testing.T) {
		payload := domain.RegisterReq{
			Name:         "Invalid Role User",
			Email:        "inv-role@company.co.id",
			Password:     "Password123!",
			Role:         "MANAGER", // invalid role, must be HRD or EMPLOYEE
			BaseSalary:   6000000,
			DepartmentID: deptID,
		}
		w := performRequest(router, "POST", "/api/v1/auth/register", payload, hrdToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400 Bad Request, got %d. Response: %s", w.Code, w.Body.String())
		}
	})
}

func TestIntegration_IDOR_Protection(t *testing.T) {
	router, db, teardown := setupTestRouter(t)
	defer teardown()

	deptID, _, hrdToken := seedTestData(t, db)

	// Create Employee A
	empA := domain.Employee{
		ID:           3001,
		Name:         "Employee A",
		Email:        "empa@company.co.id",
		Password:     "Password123!",
		Role:         "EMPLOYEE",
		BaseSalary:   5000000,
		DepartmentID: deptID,
	}
	if err := db.Create(&empA).Error; err != nil {
		t.Fatalf("failed to seed employee A: %v", err)
	}
	empAToken, err := generateTestToken(empA.ID, empA.Role)
	if err != nil {
		t.Fatalf("failed to generate token A: %v", err)
	}

	// Create Employee B
	empB := domain.Employee{
		ID:           3002,
		Name:         "Employee B",
		Email:        "empb@company.co.id",
		Password:     "Password123!",
		Role:         "EMPLOYEE",
		BaseSalary:   5000000,
		DepartmentID: deptID,
	}
	if err := db.Create(&empB).Error; err != nil {
		t.Fatalf("failed to seed employee B: %v", err)
	}
	empBToken, err := generateTestToken(empB.ID, empB.Role)
	if err != nil {
		t.Fatalf("failed to generate token B: %v", err)
	}

	// Seed a leave request for Employee A
	leaveA := domain.LeaveRequest{
		ID:         4001,
		EmployeeID: empA.ID,
		Reason:     "Annual Leave",
		StartDate:  time.Now().AddDate(0, 0, 1),
		EndDate:    time.Now().AddDate(0, 0, 3),
		Status:     "PENDING",
	}
	if err := db.Create(&leaveA).Error; err != nil {
		t.Fatalf("failed to seed leave request: %v", err)
	}

	// 3. IDOR Protection Tests
	t.Run("IDOR - Positive - Employee A can view their own leave request", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/leaves/4001", nil, empAToken)
		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("IDOR - Positive - HRD can view Employee A's leave request", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/leaves/4001", nil, hrdToken)
		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("IDOR - Negative - Employee B cannot view Employee A's leave request", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/leaves/4001", nil, empBToken)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got %d. Response: %s", w.Code, w.Body.String())
		}
		if !strings.Contains(w.Body.String(), "IDOR Blocked") {
			t.Errorf("expected response to indicate IDOR Blocked, got: %s", w.Body.String())
		}
	})

	t.Run("IDOR - Negative - Unauthenticated user cannot view Employee A's leave request", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/leaves/4001", nil, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d. Response: %s", w.Code, w.Body.String())
		}
	})

	t.Run("IDOR - Negative - View non-existent leave request", func(t *testing.T) {
		w := performRequest(router, "GET", "/api/v1/leaves/999999", nil, empAToken)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected status 403 Forbidden, got %d. Response: %s", w.Code, w.Body.String())
		}
	})
}
