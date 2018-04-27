# figaro-client
Figaro Blockchain 

WIP Package Layout
cmd/
    boot/
        main.go
    client/
        main.go
    masternode/
        main.go
    node/
        main.go
consensus/
    <Block, Blockchain>
    ...ConsensusEngine
crypto/
    _
    ...CryptoService
db/
    <ArchiveTrie, StateTrie>
    ...AccountService
    ...BlockService
    ...ReceiptService
    ...TransactionService
encoding/
    rlp/
        ...EncodingService
internal/
    (shame folder)
mock/
    (shared mock folder)
network/
    libp2p/
        _
        ...NetworkAdapter
store/
    badger/
        _
        ...StoreService
// everything below here is either types without dependencies or abstract interfaces
account.go
block.go
blockchain.go
consensus.go
crypto.go
encoding.go
network.go 
receipt.go
store.go
transaction.go
trie.go