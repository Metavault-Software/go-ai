curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d '{
  "email": "john.doe@example.com",
  "password": "secure_password"
}' | jq .
