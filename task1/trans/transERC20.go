package trans

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
)

// transferERC20 执行 ERC-20 代币转账
func transferERC20(client *ethclient.Client, privateKeyStr, tokenAddressStr, toAddressStr string, amount int) (common.Hash, error) {
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

	// 3. 解析代币合约地址
	tokenAddress := common.HexToAddress(tokenAddressStr)

	// 4. 解析接收方地址
	toAddress := common.HexToAddress(toAddressStr)

	// 5. 获取下一个 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	// 6. 获取建议的 Gas 价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %v", err)
	}

	// 7. 构建 transfer 函数调用数据
	data, err := buildTransferData(toAddress, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to build transfer data: %v", err)
	}

	// 8. 估算 Gas Limit
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &tokenAddress,
		GasPrice: gasPrice,
		Value:    big.NewInt(0), // ERC-20 转账不需要发送 ETH
		Data:     data,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Printf("Warning: failed to estimate gas, using default 100000: %v", err)
		gasLimit = 100000 // 默认 Gas Limit
	} else {
		// 增加 20% 作为缓冲
		gasLimit = gasLimit * 120 / 100
	}

	// 9. 创建交易
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tokenAddress,
		Value:    big.NewInt(0), // ERC-20 转账不需要发送 ETH
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	// 10. 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get network ID: %v", err)
	}

	// 11. 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 12. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash(), nil
}

// buildTransferData 构建 transfer 函数调用数据
func buildTransferData(toAddress common.Address, amount int) ([]byte, error) {
	// ERC-20 的 transfer 函数 ABI
	transferABI := `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	// 解析 ABI
	parsedABI, err := abi.JSON(strings.NewReader(transferABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// 打包调用数据
	amountBig := big.NewInt(int64(amount))
	data, err := parsedABI.Pack("transfer", toAddress, amountBig)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transfer data: %v", err)
	}

	return data, nil
}
