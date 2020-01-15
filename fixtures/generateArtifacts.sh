#!/bin/bash
CRYPTOGEN=./bin/cryptogen
CONFIGTXGEN=./bin/configtxgen
CHANNEL_NAME=$1
SUPERVISION_CHANNEL="supervisionchannel"
: ${CHANNEL_NAME:="qiushichannel"}
SYS_CHANNEL_NAME="systemchannel"

## Generates Org certs using cryptogen tool
function generateCerts (){
	echo
	echo "##########################################################"
	echo "##### Generate certificates using cryptogen tool #########"
	echo "##########################################################"
	if [[ -d "crypto-config" ]]; then
      rm -Rf crypto-config
    fi
#	${CRYPTOGEN} generate --config=../fixtures/crypto-config.yaml
	echo
    cp -r /opt/shurenwei/development/blockchain/fabric-ca/fixtures/crypto-config /opt/shurenwei/development/blockchain/qiushichain-network/fixtures/
}

## Generate orderer genesis block , channel configuration transaction and anchor peer update transactions
function generateChannelArtifacts() {
	echo "##########################################################"
	echo "#########  Generating Orderer Genesis block ##############"
	echo "##########################################################"
	set -x
	${CONFIGTXGEN} -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID ${SYS_CHANNEL_NAME}
	res=$?
    set +x
    if [[ ${res} -ne 0 ]]; then
      echo "Failed to generate Orderer Genesis block update for Org1MSP..."
      exit 1
    fi

	echo
	echo "#################################################################"
	echo "### Generating channel configuration transaction '${CHANNEL_NAME}.tx' ###"
	echo "#################################################################"
	set -x
	${CONFIGTXGEN} -profile QiushiChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID ${CHANNEL_NAME} 
	res=$?
    set +x
    if [[ ${res} -ne 0 ]]; then
      echo "Failed to generate channel configuration transaction update for MSP..."
      exit 1
    fi
	
    echo
    echo "#################################################################"
    echo "#######    Generating anchor peer update for MasterMSP   ##########"
    echo "#################################################################"
    set -x
    ${CONFIGTXGEN} -profile QiushiChannel -outputAnchorPeersUpdate ./channel-artifacts/MasterMSPanchors.tx -channelID ${CHANNEL_NAME} -asOrg MasterMSP
    res=$?
    set +x
    if [[ ${res} -ne 0 ]]; then
      echo "Failed to generate anchor peer update for MasterMSP..."
      exit 1
    fi

    echo
    echo "#################################################################"
    echo "#######    Generating anchor peer update for MerchantMSP   ##########"
    echo "#################################################################"
    set -x
    ${CONFIGTXGEN} -profile QiushiChannel -outputAnchorPeersUpdate ./channel-artifacts/MerchantMSPanchors.tx -channelID ${CHANNEL_NAME} -asOrg MerchantMSP
    res=$?
    set +x
    if [[ ${res} -ne 0 ]]; then
      echo "Failed to generate anchor peer update for MerchantMSP..."
      exit 1
    fi
    echo
    
    echo
	echo "#################################################################"
	echo "### Generating channel configuration transaction '${SUPERVISION_CHANNEL}.tx' ###"
	echo "#################################################################"
	set -x
	${CONFIGTXGEN} -profile QiushiChannel -outputCreateChannelTx ./channel-artifacts/${SUPERVISION_CHANNEL}.tx -channelID ${SUPERVISION_CHANNEL} 
	res=$?
    set +x
    if [[ ${res} -ne 0 ]]; then
      echo "Failed to generate channel configuration transaction..."
      exit 1
    fi
}

generateCerts
generateChannelArtifacts