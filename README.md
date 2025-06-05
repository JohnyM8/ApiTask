# Quick Start

1. If you do not have go installed, folow the link https://go.dev/

2. After cloning repo, to start a server type in terminal
```
go run server.go
```

# Connecting to local database

To connect to local database you need to change string in server.go file in 57 line.

Template goes like this:
"postgres://User:Password@localhost:5432/DBName?sslmode=disable"

# Creating table

In order to make server cooperate with databe you need to create the table 
```
Create table  wallets (
   ID serial PRIMARY KEY,
   address VARCHAR(255),
   balance INT
);
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
