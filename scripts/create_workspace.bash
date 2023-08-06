curl -X POST http://localhost:8080/api/v1/users/UserKlKBJbUwPjw/workspaces -H "Content-Type: application/json" -d '{
  "id": "WorkspaceKlSkJbUwPjw",
  "name": "My Workspace",
  "description": "A description of my workspace",
  "is_active": true,
  "owner_id": "UserKlKBJbUwPjw"
}' | jq .
