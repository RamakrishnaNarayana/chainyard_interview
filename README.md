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

Run "npm install"

node app.js

successful execution will show the following output
----------------------------------------------
Server started on 4000
----------------------------------------------

# Postman Collection for api testing

# api uses token
please copy the generated token in collection variables

# Manufacturer
https://www.getpostman.com/collections/340f58694a48d1e3ffdc

# Dealer
https://www.getpostman.com/collections/ed1c707b62e079302727
