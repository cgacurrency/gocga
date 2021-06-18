GoCGA
======

[![buddy pipeline](https://app.buddy.works/pm-45/gocga/pipelines/pipeline/333452/badge.svg?token=71e07ab3f9492e604ce68855f6f4336691f8e5ed8a7964508643c14be25d2597 "buddy pipeline")](https://app.buddy.works/pm-45/gocga/pipelines/pipeline/333452)

GoCGA is primarily a command-line tool written in Go for creating and managing CGA wallets, with the ability to both send and receive amounts. GoCGA also comprises a `wallet` package which handles key management and an `rpc` package for communicating with CGA nodes.


## Prerequirement
```
sudo apt-get update
sudo apt-get install ocl-icd-opencl-dev
sudo apt-get install build-essential
```
---

Install
-------

    go get github.com/cgacurrency/gocga

Command-line usage
------------------

    gocga add

Add a wallet with a known or random seed. The seed can be of the regular CGA variety, namely a 64 character hexadecimal string. Or a BIP39 mnemonic phrase may be used, in which case the passphrase must be known (enter it as the password). Seeds are stored encrypted with the supplied password.

    gocga add ledger

Add a Ledger hardware wallet.

    gocga list

Lists all the wallets that `gocga` knows about. Each wallet is indexed by a number which must be supplied as the `--wallet` (`-w`) flag when intending to operate on that wallet. For instance,

    gocga list -w0

will list the known accounts within wallet #0. To create a new account within wallet #0, use:

    gocga add -w0

This will prompt for the password to wallet #0.

    gocga send -a <account to send from> <account to send to> <amount>

Sends an amount of CGA from one account to another. The source account (supplied as the `--account` or `-a` flag) must be known to one of the wallets. Proof-of-work generation is built-in and uses as many threads as the number of cores. Alternatively an RPC endpoint may be used for work generation, which defaults to `http://[::1]:7132` (can be specified with the `-s` flag).

    gocga receive -w0

Receives all pending amounts for wallet #0. To receive pending amounts for a single account,

    gocga receive -a <account>

Finally,

    gocga rescan -w0

Rescans the wallet for non-empty accounts and adds them to the wallet, and

    gocga change -a <account> <representative>

Changes the representative for an account.

`wallet` package
----------------

    func NewWallet(seed []byte) (w *Wallet, err error)
    func NewBip39Wallet(mnemonic, password string) (w *Wallet, err error)
    func NewLedgerWallet() (w *Wallet, err error)

Main entrypoint to the package. The first function creates a wallet using a traditional seed, the second uses a BIP39 mnemonic and passphrase.

    func (w *Wallet) ScanForAccounts() (err error)

Scans the wallet for non-empty accounts, including not yet opened accounts with pending amounts.

    func (w *Wallet) NewAccount(index *uint32) (a *Account, err error)

Creates a new account within the wallet. If `index` is non-`nil`, derives the account from the seed with the given index.

    func (w *Wallet) GetAccount(address string) *Account
    func (w *Wallet) GetAccounts() (accounts []*Account)

Gets the account with the given `address`, or all the accounts known to the wallet.

    func (w *Wallet) ReceivePendings() (err error)

Receives all pending amounts to the wallet.

    func (a *Account) Address() string

Get the address of the account.

    func (a *Account) Index() uint32

Get the derivation index of the account.

    func (a *Account) Balance() (balance, pending *big.Int, err error)

Get the owned balance and pending balance of the account. Amounts are in raws.

    func (a *Account) Send(account string, amount *big.Int) (hash rpc.BlockHash, err error)

Send `amount` CGA from this to another `account`. The block hash is returned.

    func (a *Account) ReceivePendings() (err error)

Receives all pending amounts to the account.

    func (a *Account) ChangeRep(representative string) (hash rpc.BlockHash, err error)

Change the representative of the account.

    func CGAAmountFromString(s string) (n CGAAmount, err error)

Convert the string representation of an amount to a `CGAAmount`. A `CGAAmount` has a `Raw` field containing the amount in raws. A `CGAAmount` can also be converted back to a `string`.

`rpc` package
-------------

    rpcClient := rpc.Client{URL: "..."}

Create an RPC client with URL.

Not all RPCs are supported. The following methods are available:

    func (c *Client) AccountBalance(account string) (balance, pending *RawAmount, err error)
    func (c *Client) AccountBlockCount(account string) (blockCount uint64, err error)
    func (c *Client) AccountHistory(account string, count int64, head BlockHash) (history []AccountHistory, previous BlockHash, err error)
    func (c *Client) AccountHistoryRaw(account string, count int64, head BlockHash) (history []AccountHistoryRaw, previous BlockHash, err error)
    func (c *Client) AccountInfo(account string) (info AccountInfo, err error)
    func (c *Client) AccountRepresentative(account string) (representative string, err error)
    func (c *Client) AccountWeight(account string) (weight *RawAmount, err error)
    func (c *Client) AccountsBalances(accounts []string) (balances map[string]*AccountBalance, err error)
    func (c *Client) AccountsFrontiers(accounts []string) (frontiers map[string]BlockHash, err error)
    func (c *Client) AccountsPending(accounts []string, count int64) (pending map[string]HashToPendingMap, err error)
    func (c *Client) AvailableSupply() (available *RawAmount, err error)
    func (c *Client) BlockAccount(hash BlockHash) (account string, err error)
    func (c *Client) BlockConfirm(hash BlockHash) (started bool, err error)
    func (c *Client) BlockCount() (cemented, count, unchecked uint64, err error)
    func (c *Client) BlockInfo(hash BlockHash) (info BlockInfo, err error)
    func (c *Client) Blocks(hashes []BlockHash) (blocks map[string]*Block, err error)
    func (c *Client) BlocksInfo(hashes []BlockHash) (blocks map[string]*BlockInfo, err error)
    func (c *Client) Chain(block BlockHash, count int64) (blocks []BlockHash, err error)
    func (c *Client) Delegators(account string) (delegators map[string]*RawAmount, err error)
    func (c *Client) DelegatorsCount(account string) (count uint64, err error)
    func (c *Client) FrontierCount() (count uint64, err error)
    func (c *Client) Frontiers(account string, count int64) (frontiers map[string]BlockHash, err error)
    func (c *Client) Ledger(account string, count int64) (accounts map[string]AccountInfo, err error)
    func (c *Client) Process(block *Block, subtype string) (hash BlockHash, err error)
    func (c *Client) Representatives(count int64) (representatives map[string]*RawAmount, err error)
    func (c *Client) RepresentativesOnline() (representatives map[string]Representative, err error)
    func (c *Client) Republish(hash BlockHash, count, sources, destinations int64) (blocks []BlockHash, err error)
    func (c *Client) Successors(block BlockHash, count int64) (blocks []BlockHash, err error)
    func (c *Client) WorkCancel(hash BlockHash) (err error)
    func (c *Client) WorkGenerate(hash BlockHash, difficulty HexData) (work, difficulty2 HexData, multiplier float64, err error)
    func (c *Client) WorkValidate(hash BlockHash, work HexData) (validAll, validReceive bool, difficulty HexData, multiplier float64, err error)
