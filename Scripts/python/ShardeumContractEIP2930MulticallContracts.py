from web3 import Web3
import json
import os
import time
from hexbytes import HexBytes
import math

ShardeumConnectionHTTPS = "https://liberty20.shardeum.org/";
web3 = Web3(Web3.HTTPProvider(ShardeumConnectionHTTPS))

devTestnetPrivateKey = str(os.environ['devTestnetPrivateKey']);

userWallet = (web3.eth.account.from_key(devTestnetPrivateKey)).address

chainIdConnected = web3.eth.chain_id
print("chainIdConnected: " + str(chainIdConnected))

multicallContractAddress= web3.toChecksumAddress("0x41Ae7549023a7F0b6Cb7FE4d1807487b18cbAe10");
multicallContractABI = json.loads('[{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"multiCall","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"setCallOne","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"callContractToCall","outputs":[{"internalType":"contractcontractToCall","name":"","type":"address"}],"stateMutability":"view","type":"function"}]')
multicallContractDeployed = web3.eth.contract(address = multicallContractAddress , abi = multicallContractABI);

contractOneAddress = multicallContractDeployed.functions.callContractToCall().call()
print("contractOneAddress: "+contractOneAddress)
contractOneABI = [{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"slot0","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]
contractOneDeployed = web3.eth.contract(address = contractOneAddress , abi = contractOneABI);

slot0 = contractOneDeployed.functions.slot0().call()
print("slot0: "+ str(slot0) )

getCodeHashContractAddress= web3.toChecksumAddress("0x7d0a188781AdEf198f7056268d2FDa3F1B534287");
getCodeHashContractABI = json.loads('[{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"getCodeHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"isContractBasedOnHash","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"isContractBasedOnSize","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]')
getCodeHashContractDeployed = web3.eth.contract(address = getCodeHashContractAddress , abi = getCodeHashContractABI);

codeHashBytes32 = getCodeHashContractDeployed.functions.getCodeHash(contractOneAddress).call()
codeHashString = "0x" + codeHashBytes32.hex()
print("contractOneAddress codeHash: " + codeHashString )

unixTime = int(math.floor( time.time()*(10**3)) )
print("UNIX TIME: " + str(unixTime) )

EIP_2930_tx = {
    'chainId' : chainIdConnected,
    'to': multicallContractAddress, #WORKS WITH REGULAR WALLETS BUT CANNOT SEND TO CONTRACT FOR SOME REASON?
    'nonce':  web3.eth.getTransactionCount(userWallet)       ,
    'gas': 2100000, #WORKS WITH 1000000. If not try : Remix > deploy and run transactions
    'gasPrice': web3.toWei('30', 'gwei'), # https://etherscan.io/gastracker
    'data' : multicallContractDeployed.encodeABI(fn_name='multiCall', args=[unixTime]) ,
    'type' : 1,
    'accessList' :
                [
                    {
                        "address" : contractOneAddress,
                        "storageKeys": [
                            "0x0000000000000000000000000000000000000000000000000000000000000000",
                            codeHashString  ##Code hash from EXTCODEHASH https://blog.finxter.com/how-to-find-out-if-an-ethereum-address-is-a-contract/
                        ]
                    }
                ]
}

signed_tx = web3.eth.account.signTransaction(EIP_2930_tx, devTestnetPrivateKey)
tx_hash = web3.toHex(web3.eth.sendRawTransaction(signed_tx.rawTransaction))
print("TX HASH: " + tx_hash)

time.sleep(15)

receipt = web3.eth.getTransactionReceipt(tx_hash)
print("TX RECEIPT: " + str(receipt) )
