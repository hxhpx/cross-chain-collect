require('dotenv').config()
const fs = require('fs');

const { ALCHEMY_ETH, PRIVATE_KEY, NOMAD_ETH_CONTRACT } = process.env

const Web3 = require('web3');
//set wss provider to connet network
const web3 = new Web3(new Web3.providers.WebsocketProvider(ALCHEMY_ETH));

//use private key to set  defaultAccount
const myAct = web3.eth.accounts.privateKeyToAccount(PRIVATE_KEY);
web3.eth.defaultAccount = myAct.address;
console.log("default account is ", web3.eth.defaultAccount)

// get the current network name to display in the log
web3.eth.net.getId().then(console.log)

//subscribe contract events
//here is the USEDT on mainnet of Multichain Bridge

// Loading the compiled contract Json
const contractABI = require('../abi/nomad-ETH-bridge0x88.json');
const contract = new web3.eth.Contract(
  contractABI,
  NOMAD_ETH_CONTRACT
);

var step = 400;
var from = 13983800;
var to = from + step;

contract.getPasetEvents('Process', {
	fromBlock: from,
	toBlock: to
}, function(error, result){
	if(error) {
		step -= 5000;
		to = from + step;
	}
	else {
		console.log(result);
	}
});

/*
function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

async function main() {
	while(to <= 15356580) {
		contract.getPastEvents({
			fromBlock: from,
			toBlock: to
		},function(error, result) {
			if(error) {
				step -= 5000;
				to = from + step;
			}
			else {
				console.log(result);
			}
	});
}
}

main();
*/