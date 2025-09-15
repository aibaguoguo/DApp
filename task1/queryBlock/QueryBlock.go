package queryBlock

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"task1/abi"
	"task1/utils"
)

// 查询erc20余额 通过abi文件调用
func queryERC20Balance(tokenAddress common.Address, address common.Address) *big.Int {
	client := utils.GetEthClient()
	instance, err := abi.NewMyToken(tokenAddress, client)
	if err != nil {
		log.Fatalf("Error NewMyToken: %v", err)
	}
	balance, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatalf("Error BalanceOf: %v", err)
	}
	return balance
}

// 查询ETH余额
func QueryETHBalance(address string) *big.Int {
	client := utils.GetEthClient()
	addr := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		log.Fatalf("Error queryETHBalance By Address: %v", err)
	}
	return balance
}

/**
 * @description: 查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等
 * @param {int64} blockId 区块号
 * @return {*types.Block} 区块信息
 */
func QueryBlock(blockId int64) *types.Block {
	client := utils.GetEthClient()
	blockNumber := big.NewInt(blockId)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("Error queryBlock By Number: %v", err)
	}

	return block
}

// 获取交易收据信息
func QueryTransactionReceipt(txHash common.Hash) *types.Receipt {
	client := utils.GetEthClient()
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatalf("Error queryTransactionReceipt By Hash: %v", err)
	}
	return receipt
}
