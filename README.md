# aDi

## make
set GOOS=linux
set GOARCH=amd64
go build -o lingl
nohup ./lingl > output 2>&1 &