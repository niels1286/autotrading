// @Title
// @Description
// @Author  Niels  2021/5/11
package main

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/niels1286/nuls-go-sdk"
	"github.com/niels1286/nuls-go-sdk/account"
	txprotocal "github.com/niels1286/nuls-go-sdk/tx/protocal"
)

const chainId = 1

func main() {
	if len(os.Args) < 3 {
		fmt.Println("autotrading <prikey> <toAddress>")
		return
	}
	prikeyHex := os.Args[1]
	toAddress := os.Args[2]
	lockTime := int64(1642359110)
	sdk := GetOfficalSdk()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			doit(prikeyHex, toAddress, lockTime, nil)
		}
	}()
	doit(prikeyHex, toAddress, lockTime, sdk)
}

func doit(prikeyHex string, address string, lockTime int64, sdk *nuls.NulsSdk) {
	fmt.Println("start......")
	count := 0
	for true {
		time.Sleep(1 * time.Second)
		if (time.Now().Unix() + 20) < lockTime {
			count++
			if count%30 == 0 {
				fmt.Println("wait unlock : " + fmt.Sprintf("%d", lockTime-time.Now().Unix()) + "s.")
			}
			continue
		}
		result := transfer(prikeyHex, address, sdk)
		if result {
			break
		}
	}
}

func transfer(prikeyHex string, address string, sdk *nuls.NulsSdk) bool {
	tx := txprotocal.Transaction{
		TxType:   txprotocal.TX_TYPE_TRANSFER,
		Time:     uint32(time.Now().Unix()),
		Remark:   []byte(fmt.Sprintf("%d", time.Now().Unix())),
		Extend:   nil,
		CoinData: nil,
		SignData: nil,
	}
	acc := CreateAccount(prikeyHex)
	nonce, balance := GetNonce(acc.Address, sdk)
	if nil == balance || balance.Cmp(big.NewInt(0)) == 0 {
		return false
	}
	from1 := txprotocal.CoinFrom{
		Coin: txprotocal.Coin{
			Address:       acc.AddressBytes,
			AssetsChainId: chainId,
			AssetsId:      1,
			Amount:        balance,
		},
		Nonce:  nonce,
		Locked: 0,
	}
	val := big.NewInt(0).Add(big.NewInt(0), balance)
	val = val.Sub(val, big.NewInt(1000000))
	to1 := txprotocal.CoinTo{
		Coin: txprotocal.Coin{
			Address:       account.AddressStrToBytes(address),
			AssetsChainId: chainId,
			AssetsId:      1,
			Amount:        val,
		},
		LockValue: 0,
	}
	coinData := txprotocal.CoinData{
		Froms: []txprotocal.CoinFrom{from1},
		Tos:   []txprotocal.CoinTo{to1},
	}
	var err error
	tx.CoinData, err = coinData.Serialize()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	hash, _ := tx.GetHash().Serialize()
	signValue, _ := acc.Sign(hash)
	txSign := txprotocal.CommonSignData{
		Signatures: []txprotocal.P2PHKSignature{{
			SignValue: signValue,
			PublicKey: acc.GetPubKeyBytes(true),
		}},
	}
	tx.SignData, _ = txSign.Serialize()
	resultBytes, _ := tx.Serialize()
	result, err := sdk.BroadcastTx(resultBytes)
	if nil != err {
		return false
	} else {

		fmt.Print(time.Now().String() + " , ")
		fmt.Println("tx hash: " + result)
		return true
	}
}

func GetNonce(address string, sdk *nuls.NulsSdk) ([]byte, *big.Int) {

	status, err := sdk.GetBalance(address, chainId, 1)
	if err != nil {
		return nil, nil
	}
	if status == nil {
		return []byte{0, 0, 0, 0, 0, 0, 0, 0}, nil
	}
	return status.Nonce, status.Balance
}

func CreateAccount(prikeyHex string) *account.Account {
	account, _ := account.GetAccountFromPrkey(prikeyHex, uint16(chainId), "NULS")
	return account
}

func GetOfficalSdk() *nuls.NulsSdk {
	return nuls.NewNulsSdk("https://api.nuls.io/jsonrpc/", "https://public1.nuls.io", 1)
	// return nuls.NewNulsSdk("http://beta.api.nuls.io/jsonrpc/", "http://beta.public1.nuls.io", 2)
}
