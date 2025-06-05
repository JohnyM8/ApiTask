package graph

import (
	"ApiTask/graph/model"
	"database/sql"
	"sync"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *sql.DB

	walletsNames   []string
	namesListMutex *sync.RWMutex

	walletsMutexes         map[string]*sync.RWMutex
	wallets                map[string]*model.Wallet
	walletsInUse           map[string]bool
	walletsAmountToBalance map[string]int
}

func (r *Resolver) InitWallets() {
	if r.DB == nil {
		return
	}

	r.walletsMutexes = make(map[string]*sync.RWMutex)
	r.wallets = make(map[string]*model.Wallet)
	r.walletsInUse = make(map[string]bool)
	r.walletsAmountToBalance = make(map[string]int)

	r.namesListMutex = &sync.RWMutex{}

	rows, err := r.DB.Query("SELECT address, balance FROM wallets")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		w := &model.Wallet{}
		err := rows.Scan(&w.Address, &w.Balance)

		if err != nil {
			return
		}

		r.AddWallet(w)
	}
}
func (r *Resolver) AddWallet(wallet *model.Wallet) {
	r.namesListMutex.Lock()

	r.walletsNames = append(r.walletsNames, wallet.Address)
	r.walletsMutexes[wallet.Address] = &sync.RWMutex{}
	r.wallets[wallet.Address] = wallet
	r.walletsInUse[wallet.Address] = false
	r.walletsAmountToBalance[wallet.Address] = 0

	r.namesListMutex.Unlock()
}

func (r *Resolver) CheckIfAddressExists(address string) bool {
	r.namesListMutex.Lock()

	flag := false

	for _, walletName := range r.walletsNames {
		if walletName == address {
			flag = true
		}
	}

	if !flag {
		return false
	}

	r.namesListMutex.Unlock()

	return true
}

func (r *Resolver) UpdateWalletBalanceDBandMAP(address string, newBalance int32) bool {
	_, err := r.DB.Query("UPDATE wallets SET balance = $1 WHERE address = $2", newBalance, address)
	if err != nil {
		return false
	}

	r.wallets[address].Balance = int32(newBalance)

	return true
}
