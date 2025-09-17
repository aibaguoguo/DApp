## 1.安装abigen
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
## 2.生成abi文件 cd ./abi
### 使用solc 编译合约
solcjs --bin MyToken.sol
### 使用solc 生成abi文件
solcjs --abi MyToken.sol
### 生成go文件
abigen --abi MyToken.abi --pkg abi --type MyToken --out MyToken.go
#### 同时指定bin和abi文件生成go文件才会有deploy方法
abigen --bin myToken_sol_MyERC20.bin --abi MyToken_sol_MyERC20.abi --pkg abi --type MyToken --out MyToken2.go
## 3.重新整理依赖
go mod tidy
## 4.使用abi
instance = myToken.new...
instance.func() 执行合约函数
