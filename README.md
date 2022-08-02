# Running the test network
./network.sh up createChannel -ca -c mychannel -s couchdb

# Deploy chaincode
./network.sh deployCC -ccn car -ccp ../chaincode/car -ccl go

# Copy connection profiles in "testnetwork/organizations/peerOrganizations/org1.example.com" to "api/config" folder
testnetwork/organizations/peerOrganizations/org1.example.com/connection-org1.json >>>> ../api/config/connection-org1.json

# Copy connection profiles in "testnetwork/organizations/peerOrganizations/org2.example.com" to "api/config" folder
testnetwork/organizations/peerOrganizations/org2.example.com/connection-org2.json >>>> ../api/config/connection-org2.json

# Start Node SDK Server
cd ../api/
node app.js

successful execution will show the following output
----------------------------------------------
Server started on 4000
----------------------------------------------

