# Test Execution Results

All integration tests for RBAC, Employee Creation, and IDOR protection have passed successfully. Below is the test output showing the results of running `go test -v ./routes`.

## Test Summary
- **Tests Executed**: 2 main suites, containing 14 sub-test cases.
- **Pass Rate**: 100% (14/14 passed)
- **Status**: SUCCESS

---

## Detailed Test Logs

```
=== RUN   TestIntegration_RBAC_And_PostEmployeeData
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Positive_-_HRD_can_register_employee
[GIN] 2026/06/30 - 15:42:55 | 201 |  67.62ms |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Negative_-_Employee_cannot_register_employee
[GIN] 2026/06/30 - 15:42:55 | 403 |    2.4ms |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Negative_-_Unauthenticated_user_cannot_register_employee
[GIN] 2026/06/30 - 15:42:55 | 401 |      0s |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Negative_-_Employee_cannot_access_process_payroll
[GIN] 2026/06/30 - 15:42:55 | 403 |       0s |                 | POST     "/api/v1/payroll/process"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Email_domain_must_be_@company.co.id
[GIN] 2026/06/30 - 15:42:55 | 400 |   1.54ms |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Duplicate_email
[GIN] 2026/06/30 - 15:42:55 | 201 |  56.76ms |                 | POST     "/api/v1/auth/register"
[GIN] 2026/06/30 - 15:42:55 | 400 |   2.46ms |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Invalid_fields_(short_password)
[GIN] 2026/06/30 - 15:42:55 | 400 |    1.1ms |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Invalid_fields_(negative_salary)
[GIN] 2026/06/30 - 15:42:55 | 400 |  544.2µs |                 | POST     "/api/v1/auth/register"
=== RUN   TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Invalid_fields_(invalid_role)
[GIN] 2026/06/30 - 15:42:55 | 400 |  543.9µs |                 | POST     "/api/v1/auth/register"
--- PASS: TestIntegration_RBAC_And_PostEmployeeData (0.53s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Positive_-_HRD_can_register_employee (0.07s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Negative_-_Employee_cannot_register_employee (0.00s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Negative_-_Unauthenticated_user_cannot_register_employee (0.00s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/RBAC_-_Negative_-_Employee_cannot_access_process_payroll (0.00s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Email_domain_must_be_@company.co.id (0.00s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Duplicate_email (0.06s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Invalid_fields_(short_password) (0.00s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Invalid_fields_(negative_salary) (0.00s)
    --- PASS: TestIntegration_RBAC_And_PostEmployeeData/Post_Employee_-_Negative_-_Invalid_fields_(invalid_role) (0.00s)
=== RUN   TestIntegration_IDOR_Protection
=== RUN   TestIntegration_IDOR_Protection/IDOR_-_Positive_-_Employee_A_can_view_their_own_leave_request
[GIN] 2026/06/30 - 15:42:56 | 200 |   8.21ms |                 | GET      "/api/v1/leaves/4001"
=== RUN   TestIntegration_IDOR_Protection/IDOR_-_Positive_-_HRD_can_view_Employee_A's_leave_request
[GIN] 2026/06/30 - 15:42:56 | 200 |   4.38ms |                 | GET      "/api/v1/leaves/4001"
=== RUN   TestIntegration_IDOR_Protection/IDOR_-_Negative_-_Employee_B_cannot_view_Employee_A's_leave_request
[GIN] 2026/06/30 - 15:42:56 | 403 |   2.79ms |                 | GET      "/api/v1/leaves/4001"
=== RUN   TestIntegration_IDOR_Protection/IDOR_-_Negative_-_Unauthenticated_user_cannot_view_Employee_A's_leave_request
[GIN] 2026/06/30 - 15:42:56 | 401 |       0s |                 | GET      "/api/v1/leaves/4001"
=== RUN   TestIntegration_IDOR_Protection/IDOR_-_Negative_-_View_non-existent_leave_request
[GIN] 2026/06/30 - 15:42:56 | 403 |    2.7ms |                 | GET      "/api/v1/leaves/999999"
--- PASS: TestIntegration_IDOR_Protection (0.31s)
    --- PASS: TestIntegration_IDOR_Protection/IDOR_-_Positive_-_Employee_A_can_view_their_own_leave_request (0.01s)
    --- PASS: TestIntegration_IDOR_Protection/IDOR_-_Positive_-_HRD_can_view_Employee_A's_leave_request (0.00s)
    --- PASS: TestIntegration_IDOR_Protection/IDOR_-_Negative_-_Employee_B_cannot_view_Employee_A's_leave_request (0.00s)
    --- PASS: TestIntegration_IDOR_Protection/IDOR_-_Negative_-_Unauthenticated_user_cannot_view_Employee_A's_leave_request (0.00s)
    --- PASS: TestIntegration_IDOR_Protection/IDOR_-_Negative_-_View_non-existent_leave_request (0.00s)
PASS
ok  	hris-payroll/routes	5.056s
```
