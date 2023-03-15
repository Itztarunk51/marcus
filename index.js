// import web3 library
import Web3 from 'web3';

const web3 = new Web3("https://endpoints.omniatech.io/v1/eth/goerli/public");
const Web3 = require('web3');

const contractAddress = "0xcc2bDFAA45FD5fe6021117F198491134C3F75aE9";
const contractABI = [
	{
		"inputs": [],
		"name": "computeBruteForceWin",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "currentBoard",
		"outputs": [
			{
				"internalType": "uint8",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "didThePlayerWin",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
];

const web3 = new Web3("https://sphinx.shardeum.org/");
const contract = new web3.eth.Contract(contractABI, contractAddress);

let web3js;
let selectedAddress;

async function connectWallet() {
  if (window.ethereum) {
    try {
      await window.ethereum.request({ method: 'eth_requestAccounts' });
      web3js = new Web3(window.ethereum);
      selectedAddress = await getSelectedAddress();
      console.log("Wallet connected: ", selectedAddress);
    } catch (error) {
      console.error(error);
    }
  } else if (window.web3) {
    web3js = new Web3(web3.currentProvider);
    selectedAddress = await getSelectedAddress();
    console.log("Wallet connected: ", selectedAddress);
  } else {
    alert('Please install a wallet extension to connect to Ethereum');
  }
}

async function getSelectedAddress() {
  const accounts = await web3js.eth.getAccounts();
  return accounts[0];
}

async function computeBruteForceWin() {
  await contract.methods.computeBruteForceWin().send({ from: selectedAddress });
}

async function didThePlayerWin() {
    const selectedAddress = await getSelectedAddress();
    const result = await contract.methods.didThePlayerWin(selectedAddress).call();
  
    if (result) {
      console.log("Congratulations! You won.");
    } else {
      console.log("Sorry, you did not win.");
    }
  }
  
  async function connectWallet() {
    await window.ethereum.enable();
  
    const selectedAddress = await getSelectedAddress();
    document.getElementById("connected-address").textContent = selectedAddress;
  }
  
  document.getElementById("compute-button").addEventListener("click", computeBruteForceWin);
  document.getElementById("win-check-button").addEventListener("click", didThePlayerWin);
  document.getElementById("connect-wallet-button").addEventListener("click", connectWallet);
