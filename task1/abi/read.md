## 1.安装abigen
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
## 2.生成abi文件 cd ./abi
abigen --abi MyToken.abi --pkg abi --type MyToken --out MyToken.go
## 3.重新整理依赖
go mod tidy
## 4.使用abi
instance = myToken.new...
instance.func() 执行合约函数
