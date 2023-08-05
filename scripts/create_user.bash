curl -X POST http://localhost:8080/api/v1/users -H "Content-Type: application/json" -d '{
  "ID": "UserKlKBJbUwPjw",
  "Email": "john.doe@example.com",
  "Password": "secure_password",
  "FirstName": "John",
  "LastName": "Doe",
  "IsActive": true
}' | jq .
