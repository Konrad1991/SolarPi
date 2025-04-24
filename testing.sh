#!/bin/bash

curl -k -X POST https://localhost:8080/CreateUser \
  -F "Name=testuser1" \
  -F "Password=myplaintextpassword" \
  -F "UserRootDirectory=/home/testuser"

curl -k https://localhost:8080/GetAllUsers | jq

# curl -k -X DELETE https://localhost:8080/DeleteUser/testuser

# curl -k -X DELETE https://localhost:8080/DeleteUserByID/1
