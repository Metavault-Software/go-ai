curl -X PUT http://localhost:8080/api/v1/tasks/1 -H "Content-Type: application/json" -d '{
  "id": "1",
  "labels": [
    "foo",
    "bar"
 ]
}' | jq .
