package utils

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	// 加载.env文件
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// 获取以太坊客户端
func GetEthClient() *ethclient.Client {

	sepolia_infura_url := os.Getenv("SEPOLIA_URL")
	//连接区块链网络
	client, err := ethclient.Dial(sepolia_infura_url)
	if err != nil {
		log.Fatalf("Error connecting to Sepolia network: %v", err)
	}
	return client
}
