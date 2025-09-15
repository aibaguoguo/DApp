package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"os"
	"task1/abi"
	"task1/queryBlock"
	"task1/utils"
)

func main() {

	//doQyeryBlock()
	//trans.Trans()
	//trans.DoTrans()
	doAbiFuncCall()
}

// 调用合约函数
func doAbiFuncCall() {
	client := utils.GetEthClient()
	instance, err := abi.NewMyToken(common.HexToAddress(os.Getenv("MY_TOKEN_ADDRESS")), client)
	if err != nil {
		log.Fatalf("Error NewMyToken: %v", err)
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatalf("Error Name: %v", err)
	}
	fmt.Println("name:", name)
	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(os.Getenv("KEY1_ADDRESS")))
	if err != nil {
		log.Fatalf("Error BalanceOf: %v", err)
	}
	fmt.Println("balance1:", balance)

	balance2, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(os.Getenv("KEY2_ADDRESS")))
	if err != nil {
		log.Fatalf("Error BalanceOf: %v", err)
	}
	fmt.Println("balance2:", balance2)
}

func doQyeryBlock() {
	var blockId int64 = 9079729
	block := queryBlock.QueryBlock(blockId)
	// 打印区块信息
	fmt.Printf("Block Number: %d\n", block.Number())
	fmt.Printf("Block Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Block Time: %d\n", block.Time())
	fmt.Printf("Block Transactions: %d\n", block.Transactions().Len())
	// 打印交易信息
	txHex := common.HexToHash("0x4f6c1ace7d743bb3e49e5b1af3d1551488d0f25f6af81fd7bb39078476efe778")
	tx := block.Transaction(txHex)
	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	// 使用最新的伦敦分叉兼容的签名恢复方法 sepolia chain id 11155111
	signer := types.LatestSignerForChainID(big.NewInt(11155111))
	fromAddress, err := types.Sender(signer, tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("交易发送方 (From): %s\n", fromAddress.Hex())
	fmt.Printf("交易接收方 (To): %s\n", tx.To().Hex())

	fmt.Printf("Transaction Value: %s\n", tx.Value().String())
	fmt.Printf("Transaction Gas: %d\n", tx.Gas())
	fmt.Printf("Transaction Gas Price: %s\n", tx.GasPrice().String())
	fmt.Printf("Transaction Nonce: %d\n", tx.Nonce())

	// 如果需要，可以进一步获取交易收据查看状态等信息
	receipt := queryBlock.QueryTransactionReceipt(tx.Hash())
	fmt.Printf("交易状态: 已确认 (Confirmed), 区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 消耗: %d\n", receipt.GasUsed)
}
