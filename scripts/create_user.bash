curl -X POST http://localhost:8080/api/v1/users -H "Content-Type: application/json" -d '{
  "id": "UserKlKBJbUwPjw",
  "email": "john.doe@example.com",
  "password": "secure_password",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": true
}' | jq .
