curl -X POST http://localhost:8080/api/v1/users/UserKlKBJbUwPjw/workspaces/WorkspaceKlSkJbUwPjw/workflows -H "Content-Type: application/json" -d '{
  "ID": "Workflow1234",
  "Name": "Example Workflow",
  "Description": "This is an example workflow",
  "Tasks": {},
  "IsActive": true
}' | jq .
