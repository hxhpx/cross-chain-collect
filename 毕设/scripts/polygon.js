require('dotenv').config()
const fs = require('fs');

const { POLYGON_WS, PRIVATE_KEY, POLYGON_CONTRACT, PUBLIC_ADDRESS } = process.env

const Web3 = require('web3');
//set wss provider to connet network
const web3 = new Web3(new Web3.providers.WebsocketProvider(POLYGON_WS));

//use private key to set  defaultAccount
const myAct = web3.eth.accounts.privateKeyToAccount(PRIVATE_KEY);
web3.eth.defaultAccount = myAct.address;
console.log("default account is ", web3.eth.defaultAccount)

// get the current network name to display in the log
web3.eth.net.getId().then(console.log)

//subscribe contract events
//here is the USEDT on mainnet of Multichain Bridge

// Loading the compiled contract Json
const contractABI = require('../abi/polygon_abi.json');
const contract = new web3.eth.Contract(
  contractABI,
  POLYGON_CONTRACT
);

//getPastEvent
/*
contract.getPastEvents('LogAnySwapOut', {
	fromBlock:15808104,
	toBlock:'latest'
},function(error, result) {
	if(error)
		console.log(error);
	else
		console.log(result);
});
*/


//keep listening
contract.events.LogAnySwapIn({
	fromBlock: 34691612
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
