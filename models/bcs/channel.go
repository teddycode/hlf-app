package bcs

import (
	"bytes"
	"encoding/json"
	"github.com/fabric-app/pkg/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	cb "github.com/hyperledger/fabric/protos/common"
	"strconv"
)

func (c *Client) GetBlockHeight() (string, error) {
	chainInfo, err := c.lc.QueryInfo()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(int64(chainInfo.BCI.Height), 10), nil
}

func (c *Client) QueryTxByID(txid string, endpoint string) ([]byte, error) {
	transactions, err := c.lc.QueryTransaction(fab.TransactionID(txid), ledger.WithTargetEndpoints(endpoint))
	if err != nil {
		logging.Error("QueryTransaction error: " + err.Error())
		return nil, err
	}

	//var reads, writes []string
	var writes []string
	bf := bytes.Buffer{}
	rws := getReadWriteSet((*cb.Envelope)(transactions.GetTransactionEnvelope()))
	for _, u := range rws {
		for _, w := range u.KVRWSet.Writes {
			if w.IsDelete {
				continue
			}
			bw, _ := json.Marshal(w)
			writes = append(writes, string(bw))
		}
	}

	bf.WriteString("{\"writes\":[")
	for i, x := range writes {
		bf.WriteString(x)
		if i < len(writes)-1 {
			bf.WriteString(",")
		}
	}
	bf.WriteString("]}")

	return bf.Bytes(), nil
}
