package info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastos/Elastos.ELA/cmd/common"
	"github.com/elastos/Elastos.ELA/utils/http"
	"github.com/elastos/Elastos.ELA/utils/http/jsonrpc"

	"github.com/urfave/cli"
)

func printFormat(data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	json.Indent(buf, dataBytes, "", "    ")
	fmt.Println(string(buf.Bytes()))
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "info",
		Usage:       "show node information",
		Description: "With ela-cli info, you could look up node status, query blocks, transactions, etc.",
		ArgsUsage:   "[args]",
		Subcommands: []cli.Command{
			{
				Name:  "getconnectioncount",
				Usage: "Show how many peers are connected",
				Action: func(c *cli.Context) error {
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getconnectioncount", util.Params{})
					if err != nil {
						fmt.Println("error: get node connections failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getneighbors",
				Usage: "Show neighbor nodes information",
				Action: func(c *cli.Context) error {
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getneighbors", util.Params{})
					if err != nil {
						fmt.Println("error: get node neighbors info failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getnodestate",
				Usage: "Show current node status",
				Action: func(c *cli.Context) error {
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getnodestate", util.Params{})
					if err != nil {
						fmt.Println("error: get node state info failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getcurrentheight",
				Usage: "Get best block height",
				Action: func(c *cli.Context) error {
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getcurrentheight", util.Params{})
					if err != nil {
						fmt.Println("error: get block count failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getbestblockhash",
				Usage: "Get the best block hash",
				Action: func(c *cli.Context) error {
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getbestblockhash", util.Params{})
					if err != nil {
						fmt.Println("error: get best block hash failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getblockhash",
				Usage: "Get a block hash by height",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						cmdcom.PrintErrorMsg("Missing argument. Block height expected.")
						cli.ShowCommandHelpAndExit(c, "getblockhash", 1)
					}

					height := c.Args().First()
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getblockhash", util.Params{"height": height})
					if err != nil {
						fmt.Println("error: get block hash failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getblock",
				Usage: "Get a block details by height or block hash",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						cmdcom.PrintErrorMsg("Missing argument. Block height or hash expected.")
						cli.ShowCommandHelpAndExit(c, "getblock", 1)
					}
					param := c.Args().First()

					height, err := strconv.ParseInt(param, 10, 64)
					if err == nil {
						result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getblockhash", util.Params{"height": height})
						if err != nil {
							fmt.Println("error: get block failed,", err)
							return err
						}
						param = result.(string)
					}
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getblock", util.Params{"blockhash": param, "verbosity": 2})
					if err != nil {
						fmt.Println("error: get block failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getrawtransaction",
				Usage: "Get raw transaction by transaction hash",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						cmdcom.PrintErrorMsg("Missing argument. Transaction hash expected.")
						cli.ShowCommandHelpAndExit(c, "getrawtransaction", 1)
					}
					param := c.Args().First()

					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getrawtransaction", util.Params{"txid": param})
					if err != nil {
						fmt.Println("error: get transaction failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
			{
				Name:  "getrawmempool",
				Usage: "Get transaction details in node mempool",
				Action: func(c *cli.Context) error {
					result, err := jsonrpc.CallParams(cmdcom.LocalServer(), "getrawmempool", util.Params{})
					if err != nil {
						fmt.Println("error: get transaction pool failed,", err)
						return err
					}
					printFormat(result)
					return nil
				},
			},
		},
	}
}
