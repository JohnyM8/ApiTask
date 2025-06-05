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

	walletsNames      []string
	namesListMutex    *sync.RWMutex
	databaseMutex     *sync.RWMutex
	threadsCountMutex *sync.RWMutex

	walletsMutexes               map[string]*sync.RWMutex
	wallets                      map[string]*model.Wallet
	walletsInUse                 map[string]bool
	walletsLockedPositiveThreads map[string]int
}

func (r *Resolver) InitWallets() {
	if r.DB == nil {
		return
	}

	r.walletsMutexes = make(map[string]*sync.RWMutex)
	r.wallets = make(map[string]*model.Wallet)
	r.walletsInUse = make(map[string]bool)
	r.walletsLockedPositiveThreads = make(map[string]int)

	r.namesListMutex = &sync.RWMutex{}
	r.databaseMutex = &sync.RWMutex{}
	r.threadsCountMutex = &sync.RWMutex{}

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

		r.AddWalletMAP(w)
	}
}
func (r *Resolver) AddWalletMAP(wallet *model.Wallet) {

	r.namesListMutex.Lock()

	r.walletsNames = append(r.walletsNames, wallet.Address)

	r.namesListMutex.Unlock()

	r.walletsMutexes[wallet.Address] = &sync.RWMutex{}
	r.wallets[wallet.Address] = wallet
	r.walletsInUse[wallet.Address] = false
	r.walletsLockedPositiveThreads[wallet.Address] = 0
}

func (r *Resolver) CheckIfAddressExists(address string) bool {

	r.namesListMutex.Lock()
	defer r.namesListMutex.Unlock()

	flag := false

	for _, walletName := range r.walletsNames {
		if walletName == address {
			flag = true
		}
	}

	return flag
}

func (r *Resolver) UpdateWalletBalanceDBandMAP(address string, newBalance int32) bool {

	r.databaseMutex.Lock()
	defer r.databaseMutex.Unlock()

	_, err := r.DB.Query("UPDATE wallets SET balance = $1 WHERE address = $2", newBalance, address)
	if err != nil {
		return false
	}

	r.wallets[address].Balance = int32(newBalance)

	return true
}

func (r *Resolver) AddWalletDB(wallet *model.Wallet) bool {

	r.databaseMutex.Lock()
	defer r.databaseMutex.Unlock()

	_, err := r.DB.Exec("INSERT INTO wallets (address, balance) VALUES ($1, $2)", wallet.Address, wallet.Balance)

	return err == nil
}
