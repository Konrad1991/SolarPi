#!/bin/bash

# post file to route /CreateFile
# 	file_name := c.Request.FormValue("FileName")
# 	file, _, err := c.Request.FormFile("file")
curl -X POST -F "file=@TestFile.txt" -F "FileName=DesiredFileName" http://localhost:8080/UploadFile
