import NonFungibleToken from 0x631e88ae7f1d7c20

pub contract Cyber: NonFungibleToken {
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event Minted(id: UInt64, name: String, artist:String, description:String, arLink:String ,ipfsLink: String,MD5Hash: String,type: String)

    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath
    pub let MinterStoragePath: StoragePath

    pub var totalSupply: UInt64

    pub resource interface NFTPublic {
        pub let id: UInt64
        pub let metadata: Metadata
    }

    pub struct Metadata {
        pub let name: String
        pub let artist: String
        pub let description: String
        pub let arLink: String
        pub let ipfsLink: String
        pub let MD5Hash: String
        pub let type: String

        init(name: String,artist: String,description: String,arLink: String,ipfsLink: String,MD5Hash: String,type: String) {
            self.name=name
            self.artist=artist
            self.description=description
            //Stored in the arweave
            self.arLink=arLink
            //Stored in the ipfs
            self.ipfsLink=ipfsLink
            //MD5 hash of file
            self.MD5Hash=MD5Hash
            self.type=type
        }
    }

   pub resource NFT: NonFungibleToken.INFT, NFTPublic {
        pub let id: UInt64
        pub let metadata: Metadata
        init(initID: UInt64,metadata: Metadata) {
            self.id = initID
            self.metadata=metadata
        }
    }

    pub resource interface CyberCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowArt(id: UInt64): &Cyber.NFT? {
            post {
                (result == nil) || (result?.id == id):
                    "Cannot borrow Cyber reference: The ID of the returned reference is incorrect"
            }
        }
    }

    pub resource Collection: CyberCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic {
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) ?? panic("missing NFT")

            emit Withdraw(id: token.id, from: self.owner?.address)

            return <-token
        }

        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @Cyber.NFT

            let id: UInt64 = token.id

            let oldToken <- self.ownedNFTs[id] <- token

            emit Deposit(id: id, to: self.owner?.address)

            destroy oldToken
        }


        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return &self.ownedNFTs[id] as &NonFungibleToken.NFT
        }

        pub fun borrowArt(id: UInt64): &Cyber.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = &self.ownedNFTs[id] as auth &NonFungibleToken.NFT
                return ref as! &Cyber.NFT
            } else {
                return nil
            }
        }

        destroy() {
            destroy self.ownedNFTs
        }

        init () {
            self.ownedNFTs <- {}
        }
    }

    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create Collection()
    }

    pub struct NftData {
        pub let metadata: Cyber.Metadata
        pub let id: UInt64
        init(metadata: Cyber.Metadata, id: UInt64) {
            self.metadata= metadata
            self.id=id
        }
    }

    pub fun getNft(address:Address) : [NftData] {
        var artData: [NftData] = []
        let account=getAccount(address)

        if let artCollection= account.getCapability(self.CollectionPublicPath).borrow<&{Cyber.CyberCollectionPublic}>()  {
            for id in artCollection.getIDs() {
                var art=artCollection.borrowArt(id: id)
                artData.append(NftData(metadata: art!.metadata,id: id))
            }
        }
        return artData
    }

	pub resource NFTMinter {
		pub fun mintNFT(
		recipient: &{NonFungibleToken.CollectionPublic},
		name: String,
        artist: String,
        description: String,
        arLink: String,
        ipfsLink: String,
        MD5Hash: String,
        type: String) {
            emit Minted(id: Cyber.totalSupply,  name: name,artist:artist,description:description,arLink:arLink,ipfsLink: ipfsLink,MD5Hash: MD5Hash,type:type)

			recipient.deposit(token: <-create Cyber.NFT(
			    initID: Cyber.totalSupply,
			    metadata: Metadata(
                    name: name,
                    artist: artist,
                    description:description,
                    arLink:arLink,
                    ipfsLink:ipfsLink,
                    MD5Hash:MD5Hash,
                    type:type
                )))

            Cyber.totalSupply = Cyber.totalSupply + (1 as UInt64)
		}
	}

    init() {
        self.CollectionStoragePath = /storage/CyberCollection
        self.CollectionPublicPath = /public/CyberCollection
        self.MinterStoragePath = /storage/CyberMinter

        self.totalSupply = 0

        let minter <- create NFTMinter()
        self.account.save(<-minter, to: self.MinterStoragePath)

        emit ContractInitialized()
    }
}