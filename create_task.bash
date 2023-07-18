curl -X POST http://localhost:8080/api/v1/tasks -H "Content-Type: application/json" -d '{
  "id": "1",
  "name": "OpenAI Chat Task",
  "executor": "OpenAIAgent",
  "args": {
    "messages": [
      {
        "role": "user",
        "content": "Write a Go program to illustrate Go powerful concurrency model"
      }
    ]
  }
}' | jq .
