package cmd

import (
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/cgacurrency/gocga/rpc"
	"github.com/cgacurrency/gocga/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List wallets or accounts within a wallet",
	Run: func(cmd *cobra.Command, args []string) {
		if walletIndex < 0 {
			for i, wi := range wallets {
				n := len(wi.Accounts)
				switch n {
				case 1:
					fmt.Printf("%d: %d account\n", i, n)
				default:
					fmt.Printf("%d: %d accounts\n", i, n)
				}
			}
		} else {
			checkWalletIndex()
			var accounts []string
			for address := range wallets[walletIndex].Accounts {
				accounts = append(accounts, address)
			}
			sort.Strings(accounts)
			rpcClient := rpc.Client{URL: rpcURL}
			var balanceSum, pendingSum big.Int
			for _, address := range accounts {
				balance, pending, err := rpcClient.AccountBalance(address)
				fatalIf(err)
				balanceSum.Add(&balanceSum, &balance.Int)
				pendingSum.Add(&pendingSum, &pending.Int)
				fmt.Print(address)
				printAmounts(&balance.Int, &pending.Int)
			}
			if len(accounts) > 1 {
				fmt.Print(strings.Repeat(" ", 61), "Sum:")
				printAmounts(&balanceSum, &pendingSum)
			}
		}
	},
}

func printAmounts(balance, pending *big.Int) {
	if balance.Sign() > 0 {
		fmt.Printf(" %s", util.CGAAmount{Raw: balance})
	}
	if pending.Sign() > 0 {
		fmt.Printf(" (+ %s pending)", util.CGAAmount{Raw: pending})
	}
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
