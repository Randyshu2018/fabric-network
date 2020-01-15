make orderer NODE_NAME=orderer1.ejubc.com PORT=7050
make orderer NODE_NAME=orderer2.ejubc.com PORT=8050
make orderer NODE_NAME=orderer3.ejubc.com PORT=9050

make anchor-peer NODE_NAME=peer0.master.ejubc.com MSP_ID=MasterMSP LISTEN_PORT=7051 EVENT_PORT=7053 DOMAIN_NAME=master.ejubc.com COUCH_DB_NAME=couchdb1 COUCH_DB_PORT=5984
make peer        NODE_NAME=peer1.master.ejubc.com MSP_ID=MasterMSP LISTEN_PORT=8051 EVENT_PORT=8053 DOMAIN_NAME=master.ejubc.com COUCH_DB_NAME=couchdb2 COUCH_DB_PORT=6984
make anchor-peer NODE_NAME=peer0.merchant.ejubc.com MSP_ID=MerchantMSP LISTEN_PORT=9051 EVENT_PORT=9053 DOMAIN_NAME=merchant.ejubc.com COUCH_DB_NAME=couchdb3 COUCH_DB_PORT=7984
make peer        NODE_NAME=peer1.merchant.ejubc.com MSP_ID=MerchantMSP LISTEN_PORT=10051 EVENT_PORT=10053 DOMAIN_NAME=merchant.ejubc.com COUCH_DB_NAME=couchdb4 COUCH_DB_PORT=8984

IMAGE_TAG=1.4.3 docker-compose -f cli.yaml up -d

#docker exec cli scripts/script.sh

#scp -r admin@10.122.144.8:/opt/qiushichain-network/fixtures/crypto-config /opt/shurenwei/development/blockchain/qiushichain-network/fixtures/
