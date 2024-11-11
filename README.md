# Go-Assessment

Invoke-WebRequest -Uri "http://localhost:8080/signup" `
    -Method POST `
    -Headers @{"Content-Type"="application/json"} `
    -Body '{"email": "test@example.com", "password": "password123"}'


Invoke-WebRequest -Uri "http://localhost:8080/signin" `
    -Method POST `
    -Headers @{"Content-Type"="application/json"} `
    -Body '{"email": "test@example.com", "password": "password123"}'


Invoke-WebRequest -Uri "http://localhost:8080/refresh" `
    -Method POST `
    -Headers @{"Authorization"="Bearer <your token here>
"}`


Invoke-WebRequest -Uri "http://localhost:8080/protected" `
    -Method GET `
    -Headers @{"Authorization"="Bearer <your token here>
"}`
