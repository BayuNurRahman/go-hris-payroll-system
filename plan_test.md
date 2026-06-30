# Testing Plan: RBAC, Employee Creation, and IDOR Protection

This document contains the specifications and test scenarios designed to verify Role-Based Access Control (RBAC), Employee Registration, and IDOR Protection for the HRIS & Payroll System.

---

## 1. Role (Authorization RBAC) Testing Plan

Ensure that route-level access control successfully restricts actions based on the authenticated user's role.

### Test Scenarios
- **Positive Test Case**:
  - Request: `POST /api/v1/auth/register` with `HRD` token.
  - Expected Outcome: `201 Created` (Access Granted).
- **Negative Test Cases**:
  - Request: `POST /api/v1/auth/register` with `EMPLOYEE` token.
  - Expected Outcome: `403 Forbidden` (Access Denied).
  - Request: `POST /api/v1/auth/register` with no authentication header.
  - Expected Outcome: `401 Unauthorized` (Access Denied).
  - Request: `POST /api/v1/payroll/process` with `EMPLOYEE` token.
  - Expected Outcome: `403 Forbidden` (Access Denied).

---

## 2. Post Employee Data Testing Plan

Ensure input registration fields are properly validated inside the usecase and structural layers before database insertion.

### Test Scenarios
- **Positive Test Case**:
  - Request: `POST /api/v1/auth/register` with valid input fields:
    - Name: `Valid Employee`
    - Email: `valid-emp@company.co.id`
    - Password: `Password123!` (>= 6 chars)
    - Role: `EMPLOYEE`
    - Base Salary: `6000000` (> 0)
    - Department ID: valid referenced ID
  - Expected Outcome: `201 Created`.
- **Negative Test Cases**:
  - **Email Domain Restriction**: Register with email suffix other than `@company.co.id` (e.g. `gmail-user@gmail.com`).
    - Expected Outcome: `400 Bad Request` containing "email harus menggunakan domain resmi @company.co.id".
  - **Duplicate Email Prevention**: Register an employee using an already registered email.
    - Expected Outcome: `400 Bad Request` containing "email sudah terdaftar di sistem".
  - **Short Password Check**: Register with password less than 6 characters.
    - Expected Outcome: `400 Bad Request` (Validation error).
  - **Negative Salary Check**: Register with base salary `<= 0`.
    - Expected Outcome: `400 Bad Request` (Validation error).
  - **Invalid Role Check**: Register with a role other than `HRD` or `EMPLOYEE`.
    - Expected Outcome: `400 Bad Request` (Validation error).

---

## 3. IDOR (Insecure Direct Object Reference) Protection Testing Plan

Ensure defensive database queries prevent users from accessing resource records belonging to other users.

### Test Scenarios
- **Positive Test Cases**:
  - **Employee Viewing Own Record**: Employee A requests details of Employee A's leave request (GET `/api/v1/leaves/4001`).
    - Expected Outcome: `200 OK` (Access Granted).
  - **HRD Viewing Any Record**: HRD requests details of Employee A's leave request (GET `/api/v1/leaves/4001`).
    - Expected Outcome: `200 OK` (Access Granted).
- **Negative Test Cases**:
  - **Cross-User Data Peek**: Employee B attempts to request Employee A's leave request details (GET `/api/v1/leaves/4001`).
    - Expected Outcome: `403 Forbidden` containing "IDOR Blocked" (Access Denied).
  - **Unauthenticated Peek**: Unauthenticated client attempts to retrieve Employee A's leave details.
    - Expected Outcome: `401 Unauthorized` (Access Denied).
  - **Non-Existent ID**: View a non-existent leave ID.
    - Expected Outcome: `403 Forbidden` (Access Denied).
