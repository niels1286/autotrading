# autotrading
打包命令： （代码根目录）
go build

windows版本打包命令  
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

linux版本打包命令  
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


使用方法：  
Windows:  
```
进入命令行工具
D://xxx/xxx/autotrading.exe <prikey> <toAddress>
```
Linux/macOS:
```
./autotrading-linux <prikey> <toAddress>
```