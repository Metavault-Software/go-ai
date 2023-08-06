curl -X POST http://localhost:8080/api/v1/users/UserKlKBJbUwPjw/workspaces/WorkspaceKlSkJbUwPjw/workflows -H "Content-Type: application/json" -d '{
  "id": "Workflow1234",
  "name": "Example Workflow",
  "description": "This is an example workflow",
  "tasks": {},
  "is_active": true
}' | jq .
