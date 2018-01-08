#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=org1')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo
echo "POST request Enroll on Org2 ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Barry&orgName=org2')
echo $ORG2_TOKEN
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG2 token is $ORG2_TOKEN"
echo

echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1", "peer2"],
	"chaincodeName":"itemcc",
	"chaincodePath":"github.com/item",
	"chaincodeVersion":"v0"
}'
echo
echo


echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer1","peer2"],
	"chaincodeName":"itemcc",
	"chaincodePath":"github.com/item",
	"chaincodeVersion":"v0"
}'
echo
echo

echo "POST instantiate chaincode on peer1 of Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/logchannel/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"chaincodeName":"itemcc",
	"chaincodeVersion":"v0",
	"args":[]
}'
echo
echo

sleep 5s

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/logchannel/chaincodes/itemcc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"fcn":"initItem",
	"args":["IPR","Intellectual property right","10"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

sleep 5s
echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/logchannel/chaincodes/itemcc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "fcn":"initItem",
        "args":["zhishi","Intellectual property right","10"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo

sleep 5s
echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/logchannel/chaincodes/itemcc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
        "fcn":"initItem",
        "args":["zhishichanquan","Intellectual property right","10"]
}')
echo "Transacton ID is $TRX_ID"
echo
echo


echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/logchannel/chaincodes/itemcc?peer=peer1&fcn=queryItemsByItemPropertyOwner&args=%5B%22%22%2c%22%22%2c%22Jim%22%5D&startKey=0&endKey=9" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

