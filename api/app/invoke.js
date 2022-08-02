const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const helper = require('./helper');


const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name) => {
    try {
        const ccp = await helper.getCCP(org_name);

        const walletPath = await helper.getWalletPath(org_name);
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }


        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true },
            // eventHandlerOptions: EventStrategies.NONE
	        eventHandlerOptions: {
                commitTimeout: 100,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ANYFORTX
            }
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);

       
        const network = await gateway.getNetwork(channelName);
        const contract = await network.getContract(chaincodeName);


        // Multiple smartcontract in one chaincode
        let result=null;

        switch (fcn) {
            case "ManufactureCar":
		        console.log("====== Executing Manufacture Car ======")
                result = await contract.submitTransaction(fcn, args[0]);
                console.log(`Output: Manufacture Car= ${result.toString()}`)
                result = result.toString()
                break;
            case "UpdateDealer":
                console.log("====== Executing Update Dealer ======")
                result = await contract.submitTransaction(fcn, args[0], args[1]);
                console.log(`Output: Update Dealer Car= ${result.toString()}`)
                result = result.toString()
                break;
            case "SellCar":
                console.log("====== Executing Sell Car ======")
                result = await contract.submitTransaction(fcn, args[0], args[1]);
                console.log(`Output: Update Dealer Car= ${result.toString()}`)
                result = result.toString()
                break;
            case "ChangeCarOwner":
                console.log("====== Executing Change Car Owner ======")
                result = await contract.submitTransaction(fcn, args[0], args[1], args[2]);
                console.log(`Output: Change Car owner = ${result.toString()}`)
                result = result.toString()
                break;
            default:	
                break;
        }
       
        await gateway.disconnect();

        const response_payload = {
            result: result,
            error: false,
            errorData: null
        }

        console.log("Invoke Transaction Try Block Response Payload: ", response_payload)
        return response_payload;


    } catch (error) {

        console.log(`Invoke Transaction Getting error: ${error}`)
        const response_payload = {
            result: null,
            error: true,
            errorData: error.message
        }
        console.log("Invoke Transaction Catch Block Response Payload: ", response_payload)
        return response_payload

    }
}

exports.invokeTransaction = invokeTransaction;
