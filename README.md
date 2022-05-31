# pigeon
Sign &amp; transmit transactions.

### Start the application

> Application is build for macOS, Linux and Windows systems with x86-64 architecture.

1. Download the latest [Release](https://github.com/core-coin/pigeon/releases) for your system.
    - `pigeon-darwin-x86_64` = macOS
    - `pigeon-linux-x86_64` = Linux
1. [Open the Terminal](https://support.apple.com/guide/terminal/open-or-quit-terminal-apd5265185d-f365-44cb-8b09-71a064a42125/mac) and navigate to file location. Mostly it is: `~/Downloads` You can change directory with the command: `cd ~/Downloads`
1. Grant permissions:
    - Via the terminal: `chmod +x pigeon-…` or
    - Via the properties of the file: Right click on File -> Properties -> Permissions -> Execute.
1. Start the application:
    - Via the terminal: `./pigeon-…`

### Flags

Flags:
- -d, --dry-run                   Test the schema (do not stream, do not sign)
- -f, --file `string`               Input file with transactions
- -g, --gocore `string`             Gocore RPC API endpoint (default "http://127.0.0.1:8545")
- -h, --help                      help for pigeon
- -n, --network `int`               Network to stream on (default 1)
- -o, --output `string`             Output file with signed transactions
- -p, --password-file `string`      File with password to for file
- -k, --private-key-file `string`   File with private key to sign transactions
- -s, --stream-file `string`        File for streaming transactions into blockchain
- -t, --titles                      Skip 1 line (for CSV)
- -i, --tx-ids-file `string`        File where to store streamed tx IDs
- -u, --utc-file `string`           UTC file with encoded private key
- -v, --verbosity  `int`            Verbosity (from 1 to 7) (default 2)

### Example runs

- To sign transactions offline: `pigeon -f {path to file with transactions} -u {path to UTC file} -o {path to file where to save signed transactions}`
- To sign and stream transactions: `pigeon -f {path to file with transactions} -u {path to UTC file} -p {path to file with password}`
- To sign and stream transactions(+ save streamed transaction IDs to file): `pigeon -f {path to file with transactions} -u {path to UTC file} -i {path to file where to save transactions hashes}`
- To stream signed transactions: `pigeon -s {path to file with signed transactions}`
- To stream signed transactions(+ save streamed transaction IDs to file): `pigeon -s {path to file with signed transactions} -i {path to file where to save transactions hashes}`

### Liability

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESSED OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NON-INFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE TO ANY CLAIM, DAMAGES OR
OTHER LIABILITIES, WHETHER IN AN ACTION OF A CONTRACT, TORT, OR OTHERWISE,
ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE, OR
OTHER DEALINGS IN THE SOFTWARE.

### File with transactions scheme :
- JSON: `[
  {
  "from": "from",
  "to": "to",
  "amount": 1.2,
  "energy_limit": "",
  "energy_price": ""
  },
  {
  "from": "from",
  "to": "to",
  "amount": 22,
  "energy_limit": "22000",
  "energy_price": ""
  },
  {
  "from": "from",
  "to": "to",
  "amount": 2.2,
  "energy_limit": "",
  "energy_price": "2000000000"
  },
  {
  "from": "from",
  "to": "to",
  "amount":22,
  "energy_limit": "23000",
  "energy_price": "3000000000"
  }
  ]`
- CSV: `from,to,amount,energy_limit,energy_price
  cb...,cb...,1.123,,
  cb...,cb...,1.123,22000,
  cb...,cb...,1.123,22000,2000000000
  ` or w/o titles - `
  cb...,cb...,1.123,,
  cb...,cb...,1.123,22000,
  cb...,cb...,1.123,22000,2000000000
  `
### License

Released under the [CORE License](LICENSE).

