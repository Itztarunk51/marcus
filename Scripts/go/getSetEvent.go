//Step 1: Create abi file by running: solc --abi Store.sol > store.abi
//Step 2: Create bin file by running: solc --bin Store.sol > store.bin
//Step 3: Remove comments above the abi and bin files.
//Step 4: Generate Go contract interaction file by running:  abigen --bin=store.bin --abi=store.abi --pkg=store --out=store.go
//Step 5: Run: getSetEvent.go
package main

import (
    "fmt"
    "log"
    "os"
    "context"
    "crypto/ecdsa"
    "math/big"

    store "storeProject/contracts" //LOOK AT "go.mod" FOR YOUR RELATIVE PROJECT PATH TO FIND CONTRACT INTERFACE!

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
)

func main() {

     // Use this endpoint when you are running your own node on a specific chain (no events)
     // client, chainID := clientSetup("http://localhost:8545")

     // Use this endpoint when you are running your own node on a specific chain (events allowed)
     // client, chainID := clientSetup("ws://localhost:8546")

     client, chainID := clientSetup(os.Getenv("goerliWebSocketSecureEventsInfuraAPIKey"))

     fmt.Println("chainID: ", chainID)

     contractAddress := common.HexToAddress("0xdbaA7dfBd9125B7a43457D979B1f8a1Bd8687f37")
     contract := connectContractAddress(client,contractAddress)
     fmt.Println("contract type object: ")
     fmt.Printf("%T",contract)
     fmt.Println("")

     auth, fromAddress := connectWallet(os.Getenv("devTestnetPrivateKey"),client,chainID)

     storedData := getstoredData(contract)
     fmt.Println("storedData:", storedData)

     setUintValue := big.NewInt(444)
     SetStoredDataTx(setUintValue,client,auth,fromAddress,contract);

     SubscribeToEvents(client, contractAddress, contract)

}

func clientSetup(wssConnectionURL string) (client *ethclient.Client, chainID *big.Int) {

  client, err := ethclient.Dial(wssConnectionURL)
  if err != nil {
      log.Fatal(err)
  }

  chainID, err = client.NetworkID(context.Background())
  if err != nil {
     log.Fatal(err)
  }
  return
}

func connectContractAddress(client *ethclient.Client, contractAddress common.Address) (contract *store.Store) {

  contract, err := store.NewStore(contractAddress, client)
  if err != nil {
      log.Fatal(err)
  }
  return
}

func connectWallet(privateKeyString string, client *ethclient.Client, chainID *big.Int) (auth *bind.TransactOpts, fromAddress common.Address) {

   privateKey, err := crypto.HexToECDSA(privateKeyString)
   if err != nil {
       log.Fatal(err)
   }

   publicKey := privateKey.Public()
   publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
   if !ok {
       log.Fatal("error casting public key to ECDSA")
   }

   fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

   auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
   if err != nil {
       log.Fatal(err)
   }

   return

}

func getstoredData(contract *store.Store) (storedData *big.Int) {

  storedData, err := contract.StoredData(&bind.CallOpts{})
  if err != nil {
        log.Fatal(err)
  }
  return

}

func SetStoredDataTx(setUintValue *big.Int, client *ethclient.Client, auth *bind.TransactOpts, fromAddress common.Address, contract *store.Store) {

  gasPrice, err := client.SuggestGasPrice(context.Background())
  if err != nil {
      log.Fatal(err)
  }

  nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
  if err != nil {
      log.Fatal(err)
  }

  auth.Nonce = big.NewInt(int64(nonce))
  auth.GasLimit = uint64(3000000) // in units
  auth.GasPrice = gasPrice
  auth.Value = big.NewInt(0)     // in wei

  tx, err := contract.Set(auth, setUintValue)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println("Tx hash:", tx.Hash().Hex()) // tx sent

  return
}

func SubscribeToEvents(client *ethclient.Client, contractAddress common.Address, contract *store.Store) {
  //Subscribe to events from smart contract address.
  query := ethereum.FilterQuery{
      Addresses: []common.Address{contractAddress},
  }

  logs := make(chan types.Log)
  sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
  if err != nil {
      log.Fatal(err)
  }

  for {
      select {
      case err := <-sub.Err():
          log.Fatal(err)
      case vLog := <-logs:
          fmt.Println("New Event Log:", vLog) // pointer to event log

          storedData, err := contract.StoredData(&bind.CallOpts{})
            if err != nil {
                log.Fatal(err)
          }
          fmt.Println("storedData:", storedData)
      }
  }

  return
}
