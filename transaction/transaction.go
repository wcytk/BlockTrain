package transaction

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"io/ioutil"
	"log"
	"math/big"
	"time"
)

const (
	KEYJSON_FILEDIR = "C:\\Users\\wcytk\\AppData\\Roaming\\Ethereum\\keystore\\UTC--2019-05-17T14-14-33.054977900Z--72956040eae0ba41e6ee2a009ae1ea8a504a1008"
	KEYSTORE_DIR    = "C:\\Users\\wcytk\\AppData\\Roaming\\Ethereum\\keystore"
	CHAIN_ID         = 10
)

func StartTransaction(fromAddr string, toAddr string, passphrase string) {
	client, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}

	ethClient, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("ethClient.Dial err", err)
		return
	}

	var result string
	//var result hexutil.Big
	err = client.Call(&result, "eth_getBalance", fromAddr, "latest")
	//err = ec.c.CallContext(ctx, &result, "eth_getBalance", account, "latest")

	if err != nil {
		fmt.Println("client.Call err", err)
		return
	}

	fmt.Printf("fromAccount: %s\nbalance: %s\n", fromAddr, result)

	err = client.Call(&result, "eth_getBalance", toAddr, "latest")
	//err = ec.c.CallContext(ctx, &result, "eth_getBalance", account, "latest")

	if err != nil {
		fmt.Println("client.Call err", err)
		return
	}

	fmt.Printf("toAccount: %s\nbalance: %s\n", toAddr, result)

	//fmt.Print("Input passphrase: ")
	//var passphrase string
	//_, _ = fmt.Scanln(&passphrase)

	ks := keystore.NewKeyStore(
		KEYSTORE_DIR,
		keystore.LightScryptN,
		keystore.LightScryptP)

	go transaction(ks, ethClient, fromAddr, toAddr, passphrase)
}

func transaction(ks *keystore.KeyStore, client *ethclient.Client, from string, to string, passphrase string) {
	for {
		fromAddress := common.HexToAddress(from)
		// Find the signing account
		signAcc, err := ks.Find(accounts.Account{Address: fromAddress})
		if err != nil {
			fmt.Println("account keystore find error:")
			panic(err)
		}
		fmt.Printf("account found: signAcc.addr=%s; signAcc.url=%s\n", signAcc.Address.String(), signAcc.URL)
		fmt.Println()

		// Unlock the signing account
		errUnlock := ks.Unlock(signAcc, passphrase)
		if errUnlock != nil {
			fmt.Println("account unlock error:")
			panic(err)
		}
		fmt.Printf("account unlocked: signAcc.addr=%s; signAcc.url=%s\n", signAcc.Address.String(), signAcc.URL)
		fmt.Println()

		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		value := big.NewInt(100000000000000000) // in wei (0.1 eth)
		gasLimit := uint64(21000)                // in units
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		toAddress := common.HexToAddress(to)
		var data []byte
		tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)


		keyJson, readErr := ioutil.ReadFile(KEYJSON_FILEDIR)
		if readErr != nil {
			fmt.Println("key json read error:")
			panic(readErr)
		}
		keyWrapper, keyErr := keystore.DecryptKey(keyJson, passphrase)
		if keyErr != nil {
			fmt.Println("key decrypt error:")
			panic(keyErr)
		}
		fmt.Printf("key extracted: addr = %s \n", keyWrapper.Address.String())

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(CHAIN_ID)), keyWrapper.PrivateKey)
		if err != nil {
			log.Fatal(err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
		time.Sleep(1 * time.Minute)
	}
}

// 0x72956040eae0ba41e6ee2a009ae1ea8a504a1008
// 0x599066346fe0facf2bdec95f11ee7954ff8641d6
