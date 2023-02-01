package ofac

import (
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/core/types"
)

type Handler struct {
	Chan chan *types.Transaction
}

func (h *Handler) Start() {
	go h.run()
}

func (h *Handler) run() {

	// Process the transaction
	for {
		select {
		case tx := <-h.Chan:
			h.serveHTTP(tx)
		}
	}
}

func (h *Handler) serveHTTP(tx *types.Transaction) {

	// Serialize JSON
	txJSON, err := tx.MarshalJSON()
	if err != nil {
		txJSON = []byte{}
	}

	// POST HTTP
	postParams := url.Values{}
	postParams.Add("Hash", tx.Hash().String())
	postParams.Add("To", tx.To().String())
	postParams.Add("Data", string(tx.Data()))
	postParams.Add("Value", tx.Value().String())
	postParams.Add("JSONMarshal", string(txJSON))

	resp, err := http.PostForm(params.OFACEndpoint, postParams)

	// Log results
	if err != nil {
		log.Warn("OFAC Reporting Error :", err, ". Response :", resp)
	} else {
		log.Info("Reported Transaction to OFAC. Response :", resp, ". Corresponding txHash:", tx.Hash())
	}

}

func (h *Handler) ReportTransaction(tx *types.Transaction) {
	h.Chan <- tx
}
