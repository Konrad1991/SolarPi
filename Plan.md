# 🧾 JWT Login - Kanban Board

## 🟦 To Do
- [ ] `loginUser`: Accept username & password from client
- [ ] Fetch user from DB by name
- [ ] Compare password using `bcrypt.CompareHashAndPassword(...)`
- [ ] Generate JWT token if login is successful
- [ ] Return token to client
- [ ] Add JWT middleware to protected routes
- [ ] Parse `Authorization: Bearer <token>` header in requests
- [ ] Validate token signature and expiration
- [ ] Extract `name` from token and check user exists

## 🟨 In Progress
- [ ] `loginUser` route stub exists
- [ ] Bcrypt hashing & checking already implemented

## ✅ Done
- [x] User creation with bcrypt hashing
- [x] User retrieval from DB
- [x] TLS support for secure transport
