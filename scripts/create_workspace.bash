curl -X POST http://localhost:8080/api/v1/users/UserKlKBJbUwPjw/workspaces -H "Content-Type: application/json" -d '{
  "ID": "WorkspaceKlSkJbUwPjw",
  "Name": "My Workspace",
  "Description": "A description of my workspace",
  "IsActive": true,
  "OwnerID": "UserKlKBJbUwPjw"
}' | jq .
