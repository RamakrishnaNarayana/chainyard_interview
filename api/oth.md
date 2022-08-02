
// pm2 start myServer.js --node-args="--production --port=1337"

// node --max-old-space-size=8192 --nouse-idle-notification --expose-gc

// loadtest -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDgxNjQ3MDMsInVzZXJuYW1lIjoiS2lzaG9yaVJhZGhhIiwib3JnTmFtZSI6Ik9yZzEiLCJpYXQiOjE2NDgxMjg3MDN9.Stnow3evLPZoKlP57e5l5R8q3n-ioJIaiDTaCPx9UGk" -n 100 -c 1 -k "http://localhost:4000/channels/mychannel/chaincodes/fot?args=["0001_006"]&fcn=GetProductByUniqueID"

