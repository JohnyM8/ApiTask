# Quick Start

1. If you do no have go installed, folow the link https://go.dev/

2. After cloning repo, to start a server type in terminal
```
go run server.go
```

# Example queries

To make this it queries you can http://localhost:8080/

## Create new wallet called "Wallet1"
```
mutation createWallet {
  createWallet(input: { address: "Wallet1"}) {
    address
    balance
  }
}
```

## List all wallets
```
query findWallet{
  wallets {
    address
    balance
  }
}
```

## Make a transfer of 2000 BTP tokens from wallet address "Wallet1" to "Wallet2"
```
mutation Transfer{
  transfer(fromAddress: "Wallet1", toAddress: "Wallet2", amount: 2000) {
    fromAddress
    toAddress
    amount
  }
}
```
