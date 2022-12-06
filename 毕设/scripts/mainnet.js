require('dotenv').config()
const fs = require('fs');

const { ALCHEMY_ETH, PRIVATE_KEY, CONTRACT_ETH_FOR_USDT } = process.env

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
const contractABI = require('../abi/mainnet_abi.json');
const contract = new web3.eth.Contract(
  contractABI,
  CONTRACT_ETH_FOR_USDT
);

//getPastEvent from 20018693, which is the earliest block

const latestBlock = 15836153;
var step = 40000;
var toBlk = latestBlock;
var fromBlk = toBlk - step;

function sleep(ms) {
	return new Promise(resolve => setTimeout(resolve, ms));
}

async function main() {
	while(fromBlk >= 14000000) {
		contract.getPastEvents('LogAnySwapOut', {
		  fromBlock:fromBlk,
		  toBlock:toBlk
		}, function(error, result) {
		  if(error) {
			   step -= 5000;
			   fromBlk = toBlk - step;
		  }
		  else{
			result.forEach(function(it,idx) {
			  if(it.returnValues.fromChainID == 1 && it.returnValues.from != CONTRACT_ETH_FOR_USDT){
				var txn = it.transactionHash;
				var account = it.returnValues.to;
				//var fromChainID = it.returnValues.fromChainID;
				//var toChainID = it.returnValues.toChainID;
				var amount = it.returnValues.amount;
				var obj = txn+'\t'+account+'\t'+amount+'\n';
				fs.writeFile('./pastEvent/mainnet-to-BSC-USDT.json',obj,{flag:'a'},function(error){
				  if(error)  console.log("!!!write file error \n", error);
				});
			  }
			});
			step = 40000;
			toBlk = fromBlk - 1;
			fromBlk = toBlk - step;
		  }
		});

		await sleep(3000);
	}

	if(fromBlk < 14000000)
  	  console.log("---------------- query has done -----------------")
}

main();

/*
//keep listening
contract.events.LogAnySwapOut({
	fromBlock: 15809103
},function(error, event) { console.log(event); })
.on('connected', function(subscriptionId) {
	console.log('subscriptionId is ', subscriptionId)
})
.on('data', function(event) {
	console.log('new event is ', event)
})
.on('error', function(error, receipt) {
	console.log('error! ', error)
});


/*
var subscription = web3.eth.subscribe('logs', {
	address: CONTRACT_ADDRESS,
},function(error, result) {
    if (error)
	console.log(error);
    if (!error)
        console.log(result);
  })
  .on("data", function(sync){
	console.log(sync);
  })
  .on("changed", function(isSyncing){
    if(isSyncing) {
	console.log("isSyncing: ",isSyncing);
    }
});
*/
