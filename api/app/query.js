const { Gateway, Wallets, } = require('fabric-network');
const fs = require('fs');
// const path = require("path")
// const log4js = require('log4js');
// const logger = log4js.getLogger('FOTNetwork');
const util = require('util')


const helper = require('./helper')
const query = async (channelName, chaincodeName, args, fcn, username, org_name) => {

    try {

        // load the network configuration
        // const ccpPath = path.resolve(__dirname, '..', 'config', 'connection-org1.json');
        // const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
        const ccp = await helper.getCCP(org_name) //JSON.parse(ccpJSON);

        // Create a new file system based wallet for managing identities.
        const walletPath = await helper.getWalletPath(org_name) //.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true }
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(channelName);

        // Get the contract from the network.
        const contract = await network.getContract(chaincodeName);
        let result=null;

        switch (fcn) {
            case "QueryCar":
                console.log("======= Executing Query Car =======")
                result = await contract.evaluateTransaction(fcn, args[0]);
                result = JSON.parse(result.toString());
                break;
            case "GetHistoryForCar":
	        	console.log("======= Executing Get History For Car =======")
	        	result = await contract.evaluateTransaction(fcn, args[0]);
                result = JSON.parse(result.toString());
	        	break;
            case "CarExists":
	        	console.log("======= Executing Get Car Exists =======")
	        	result = await contract.evaluateTransaction(fcn, args[0]);
                result = JSON.parse(result.toString());
	        	break;
            default:
                break;
        }

        const response_payload = {
            result: result,
            error: false,
            errorData: null
        }

        console.log("Query Transaction Try Block Response Payload: ", response_payload)
        return response_payload

    } catch (error) {

        console.log(`Query Transaction Getting error: ${error}`)
        const response_payload = {
            result: null,
            error: true,
            errorData: error.message
        }
        console.log("Query Transaction Catch Block Response Payload: ", response_payload)
        return response_payload

    }
}

exports.query = query
