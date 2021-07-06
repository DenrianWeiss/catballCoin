package service

import (
	"fmt"
	"github.com/DenrianWeiss/catballCoin/model"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/wallet"
	"math/big"
	"sync"

	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

const (
	GasModeByPercent = "percent"
	GasModeByPrice   = "price"
	TransferGas      = 21000
)

var (
	tasks = make(map[int64]chan struct{})
	tasksLock sync.Mutex
)

func NewCoinHarvestTask(config *model.CoinTask) {
	ticker := time.NewTicker(time.Duration(GlobalConfig.Interval) * time.Second)
	stop := make(chan struct{})
	tasksLock.Lock()
	tasks[config.TaskID] = stop
	tasksLock.Unlock()
	go coinHarvestWorker(ticker, stop, config)
}

func coinHarvestWorker(t *time.Ticker, stop chan struct{}, config *model.CoinTask) {
	client, err := jsonrpc.NewClient(config.EthRpc)
	if err != nil {
		tasksLock.Lock()
		delete(tasks, config.TaskID)
		tasksLock.Unlock()
		return
	}
	privateKey := big.NewInt(0)
	privateKey.SetString(config.SourceKey, 16)
	wal, err := wallet.NewWalletFromPrivKey(privateKey.Bytes())
	prevBalance := big.NewInt(0)
	if err != nil {
		tasksLock.Lock()
		delete(tasks, config.TaskID)
		tasksLock.Unlock()
		return
	}
	select {
	case <- t.C: {
		go coinHarvestAction(client, config, prevBalance, wal)
	}
	case <- stop: {
		return
	}
	}
}

func coinHarvestAction(rpc *jsonrpc.Client, config *model.CoinTask, prevBalance *big.Int, wal *wallet.Key) {
	// Get Chain id.
	chain, err := rpc.Eth().ChainID()
	if err != nil {
		fmt.Printf("task %v failed to get chain id: %v\n", config.TaskID, err)
		return
	}
	// First we query currently gas.
	gas, err := rpc.Eth().GasPrice()
	if err != nil {
		fmt.Printf("task %v failed to get gas price: %v\n", config.TaskID, err)
		return
	}
	block, err := rpc.Eth().BlockNumber()
	if err != nil {
		fmt.Printf("task %v failed to get block number: %v\n", config.TaskID, err)
		return
	}
	// Then we get the account's balance
	val, err := rpc.Eth().GetBalance(wal.Address(), web3.BlockNumber(block))
	if err != nil {
		fmt.Printf("task %v failed to get balance: %v\n", config.TaskID, err)
		return
	}
	valCpy := new(big.Int)
	valCpy.Set(val)
	if val.Cmp(big.NewInt(0)) == 0 || val.Cmp(prevBalance) < 0 || val.Cmp(prevBalance) == 0 {
		return
		// Noting to do
	}
	// Calculate gas
	var gasFee uint64
	switch config.GasMode {
	case GasModeByPercent: {
		totalGas := val.Mul(val, big.NewInt(config.Percents))
		gasFee = totalGas.Div(totalGas, big.NewInt(100 * TransferGas)).Uint64()
	}
	case GasModeByPrice: {
		gasFee = gas * uint64(config.Percents) / 100
	}
	default:{
		return
	}
	}
	totalGas := big.NewInt(0).SetUint64(gasFee)
	totalGas = totalGas.Mul(totalGas, big.NewInt(21000))
	targetAddress := web3.Address{0x1}
	err = targetAddress.UnmarshalText([]byte(config.TargetWallet))
	if err != nil {
		return
	}
	// Do transaction
	tx := &web3.Transaction{
		GasPrice: gasFee,
		Gas: 21000,
		Value: valCpy.Sub(valCpy, totalGas),
		To: &targetAddress,
	}

	signer := wallet.NewEIP155Signer(chain.Uint64())
	txn, _ := signer.SignTx(tx, wal)
	hash, err := rpc.Eth().SendRawTransaction(txn.MarshalRLP())
	SendTelegramMessage(fmt.Sprintf("Send %v wei to %v, Etherscan at https://etherscan.io/tx/%x",
		val, config.TargetWallet, hash), config.ChatID)
}
