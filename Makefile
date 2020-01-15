include .env
.PHONY: all test clean pre

all: clean pre setup dev-env 
dev: dev-env

#### Download binaries	
pre:
    ifeq ($(wildcard fixtures/bin),)
	    @echo "Preparing configtxgen cryptogen tool ..."
	    @cd fixtures && curl -ssL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s -- ${VERSION} ${VERSION} -s -d && cd ../
    endif
    ifeq ($(wildcard fixtures/crypto-config),)
	    @cd fixtures && \
        ../scripts/generateArtifacts.sh ${CHANNEL_NAME}  && cd ../ 

    endif
	@echo "Done pre"
    
setup:
	@echo "Done setup"
	
clean:
	@echo "Start clean..."
	@rm -Rf fixtures/crypto-config
	@rm -Rf fixtures/channel-artifacts/*
	@rm -Rf mount
orderer:
    ifndef NODE_NAME
	    @echo "NODE_NAME is required !"
    endif	
    ifndef PORT
	    @echo "PORT is required !"
    endif
	@sed "s/%NODE_NAME%/${NODE_NAME}/g;s/%PORT%/${PORT}/g" template/orderer.yaml > ${NODE_NAME}.yaml
	@IMAGE_TAG=${VERSION} docker-compose -f ${NODE_NAME}.yaml up -d
peer:
    ifndef NODE_NAME
	    @echo "NODE_NAME is required !"
    endif	
    ifndef MSP_ID
	    @echo "MSP_ID is required !"
    endif
    ifndef LISTEN_PORT
	    @echo "LISTEN_PORT is required !"
    endif
    ifndef EVENT_PORT
	    @echo "EVENT_PORT is required !"
    endif
    ifndef DOMAIN_NAME
	    @echo "DOMAIN_NAME is required !"
    endif
    ifndef COUCH_DB_NAME
	    @echo "COUCH_DB_NAME is required !"
    endif
    ifndef COUCH_DB_PORT
	    @echo "COUCH_DB_PORT is required !"
    endif
	@sed "s/%NODE_NAME%/${NODE_NAME}/g;s/%MSP_ID%/${MSP_ID}/g;s/%LISTEN_PORT%/${LISTEN_PORT}/g;s/%EVENT_PORT%/${EVENT_PORT}/g;s/%DOMAIN_NAME%/${DOMAIN_NAME}/g;s/%COUCH_DB_NAME%/${COUCH_DB_NAME}/g;s/%COUCH_DB_PORT%/${COUCH_DB_PORT}/g;" template/peer.yaml > ${NODE_NAME}.yaml
	@IMAGE_TAG=${VERSION} docker-compose -f ${NODE_NAME}.yaml up -d
anchor-peer:
    ifndef NODE_NAME
	    @echo "NODE_NAME is required !"
    endif	
    ifndef MSP_ID
	    @echo "MSP_ID is required !"
    endif
    ifndef LISTEN_PORT
	    @echo "LISTEN_PORT is required !"
    endif
    ifndef EVENT_PORT
	    @echo "EVENT_PORT is required !"
    endif
    ifndef DOMAIN_NAME
	    @echo "DOMAIN_NAME is required !"
    endif
    ifndef COUCH_DB_NAME
	    @echo "COUCH_DB_NAME is required !"
    endif
    ifndef COUCH_DB_PORT
	    @echo "COUCH_DB_PORT is required !"
    endif
	@sed "s/%NODE_NAME%/${NODE_NAME}/g;s/%MSP_ID%/${MSP_ID}/g;s/%LISTEN_PORT%/${LISTEN_PORT}/g;s/%EVENT_PORT%/${EVENT_PORT}/g;s/%DOMAIN_NAME%/${DOMAIN_NAME}/g;s/%COUCH_DB_NAME%/${COUCH_DB_NAME}/g;s/%COUCH_DB_PORT%/${COUCH_DB_PORT}/g;" template/anchor-peer.yaml > ${NODE_NAME}.yaml
	@IMAGE_TAG=${VERSION} docker-compose -f ${NODE_NAME}.yaml up -d
dev-env:
	@sed "s/%NODE_NAME%/orderer1.ejubc.com/g;s/%PORT%/7050/g" template/orderer.yaml > orderer1.ejubc.com.yaml
	@sed "s/%NODE_NAME%/orderer2.ejubc.com/g;s/%PORT%/8050/g" template/orderer.yaml > orderer2.ejubc.com.yaml
	@sed "s/%NODE_NAME%/orderer3.ejubc.com/g;s/%PORT%/9050/g" template/orderer.yaml > orderer3.ejubc.com.yaml
	
	@sed "s/%NODE_NAME%/peer0.master.ejubc.com/g;s/%MSP_ID%/MasterMSP/g;s/%LISTEN_PORT%/7051/g;s/%EVENT_PORT%/7053/g;s/%DOMAIN_NAME%/master.ejubc.com/g;s/%COUCH_DB_NAME%/couchdb1/g;s/%COUCH_DB_PORT%/5984/g;" template/anchor-peer.yaml > peer0.master.ejubc.com.yaml
	@sed "s/%NODE_NAME%/peer1.master.ejubc.com/g;s/%MSP_ID%/MasterMSP/g;s/%LISTEN_PORT%/8051/g;s/%EVENT_PORT%/8053/g;s/%DOMAIN_NAME%/master.ejubc.com/g;s/%COUCH_DB_NAME%/couchdb2/g;s/%COUCH_DB_PORT%/6984/g;" template/peer.yaml > peer1.master.ejubc.com.yaml
	@sed "s/%NODE_NAME%/peer0.merchant.ejubc.com/g;s/%MSP_ID%/MerchantMSP/g;s/%LISTEN_PORT%/9051/g;s/%EVENT_PORT%/9053/g;s/%DOMAIN_NAME%/merchant.ejubc.com/g;s/%COUCH_DB_NAME%/couchdb3/g;s/%COUCH_DB_PORT%/7984/g;" template/anchor-peer.yaml > peer0.merchant.ejubc.com.yaml
	@sed "s/%NODE_NAME%/peer1.merchant.ejubc.com/g;s/%MSP_ID%/MerchantMSP/g;s/%LISTEN_PORT%/10051/g;s/%EVENT_PORT%/10053/g;s/%DOMAIN_NAME%/merchant.ejubc.com/g;s/%COUCH_DB_NAME%/couchdb4/g;s/%COUCH_DB_PORT%/8984/g;" template/peer.yaml > peer1.merchant.ejubc.com.yaml
	@for i in 1 2 3 ; do \
        echo orderer$${i}"."${DOMAIN}; \
    done
	@IMAGE_TAG=${VERSION} docker-compose -f orderer1.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f orderer2.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f orderer3.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f peer0.master.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f peer1.master.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f peer0.merchant.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f peer1.merchant.ejubc.com.yaml up -d
	@IMAGE_TAG=${VERSION} docker-compose -f cli.yaml up -d
	@docker exec cli scripts/script.sh

