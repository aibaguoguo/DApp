package trans

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	"os"
	"task1/utils"
	"time"
)

func DoTrans() {

	// 1. 配置参数
	amount := 1.0                                          // 转账金额 (ETH)
	recipient := os.Getenv("KEY1_ADDRESS")                 // 接收方地址
	senderPrivateKey := os.Getenv("SEPOLIA_ACCOUNTS_KEY2") // 发送方私钥

	// 2. 连接到 Sepolia 节点
	client := utils.GetEthClient()
	defer client.Close()

	// 3. 执行转账
	txHash, err := transferETH(client, senderPrivateKey, recipient, amount)
	if err != nil {
		log.Fatalf("Transfer failed: %v", err)
	}

	fmt.Printf("Transfer initiated! Transaction hash: %s\n", txHash.Hex())

	// 4. 等待交易确认
	receipt, err := waitForTransactionReceipt(client, txHash)
	if err != nil {
		log.Fatalf("Error waiting for transaction: %v", err)
	}

	fmt.Printf("Transaction confirmed in block #%d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)
}

// transferETH 执行 ETH 转账
func transferETH(client *ethclient.Client, privateKeyStr, toAddressStr string, amount float64) (common.Hash, error) {
	// 1. 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid private key: %v", err)
	}

	// 2. 获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Hash{}, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 3. 解析接收方地址
	toAddress := common.HexToAddress(toAddressStr)

	// 4. 获取下一个 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	// 5. 获取建议的 Gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %v", err)
	}

	// 6. 估算 Gas Limit
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Value:    big.NewInt(0), // 仅估算，不发送实际金额
		GasPrice: gasPrice,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Printf("Warning: failed to estimate gas, using default 21000: %v", err)
		gasLimit = 21000 // 标准转账的默认 Gas Limit
	} else {
		// 增加 10% 作为缓冲
		gasLimit = gasLimit * 110 / 100
	}

	// 7. 计算转账金额 (ETH → Wei)
	amountWei := new(big.Int)
	amountFloat := new(big.Float).SetFloat64(amount)
	amountFloat.Mul(amountFloat, big.NewFloat(math.Pow10(18)))
	amountFloat.Int(amountWei)

	// 8. 创建交易
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amountWei,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})

	// 9. 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get network ID: %v", err)
	}

	// 10. 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 11. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash(), nil
}

// waitForTransactionReceipt 等待交易确认
func waitForTransactionReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	// 设置超时时间（3分钟）
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 轮询交易状态
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		if err != ethereum.NotFound {
			return nil, fmt.Errorf("error checking transaction receipt: %v", err)
		}

		// 等待几秒后重试
		select {
		case <-time.After(5 * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for transaction confirmation")
		}
	}
}
