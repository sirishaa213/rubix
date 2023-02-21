package core

import (
	"fmt"
	"net/http"

	"github.com/EnsurityTechnologies/config"
	"github.com/EnsurityTechnologies/ensweb"
	"github.com/EnsurityTechnologies/helper/jsonutil"
	"github.com/EnsurityTechnologies/logger"
)

const (
	ExplorerBasePath       string = "/api/v2/services/app/Rubix/"
	ExplorerCreateDIDAPI   string = "CreateOrUpdateRubixUser"
	ExplorerTransactionAPI string = "CreateOrUpdateRubixTransaction"
	ExplorerMapDIDAPI      string = "map-did"
)

type ExplorerClient struct {
	ensweb.Client
	log logger.Logger
}

type ExplorerDID struct {
	PeerID    string `json:"peerid"`
	DID       string `json:"user_did"`
	IPAddress string `json:"ipaddress"`
	Balance   int    `json:"balance"`
}

type ExplorerMapDID struct {
	OldDID string `json:"old_did"`
	NewDID string `json:"new_did"`
	PeerID string `json:"peer_id"`
}

type ExplorerTrans struct {
	TID         string   `json:"transaction_id"`
	SenderDID   string   `json:"sender_did"`
	ReceiverDID string   `json:"receiver_did"`
	TokenTime   float64  `json:"token_time"`
	TokenIDs    []string `json:"token_id"`
	Amount      float64  `json:"amount"`
	TrasnType   int      `json:"transaction_type"`
	QuorumList  []string `json:"quorum_list"`
}

type ExplorerResponse struct {
	Message string `json:"Message"`
	Status  bool   `json:"Status"`
}

func (c *Core) InitRubixExplorer() error {
	url := "deamon-explorer.azurewebsites.net"
	if c.testNet {
		url = "rubix-deamon-api.ensurity.com"
	}
	cl, err := ensweb.NewClient(&config.Config{ServerAddress: url, ServerPort: "443", Production: "true"}, c.log)
	if err != nil {
		return err
	}
	c.ec = &ExplorerClient{
		Client: cl,
		log:    c.log.Named("explorerclient"),
	}
	return nil
}

func (ec *ExplorerClient) SendExploerJSONRequest(method string, path string, input interface{}, output interface{}) error {
	req, err := ec.JSONRequest(method, ExplorerBasePath+path, input)
	if err != nil {
		return err
	}
	resp, err := ec.Do(req)
	if err != nil {
		ec.log.Error("Failed r get response from explorer", "err", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		str := fmt.Sprintf("Http Request failed with status %d", resp.StatusCode)
		ec.log.Error(str)
		return fmt.Errorf(str)
	}
	if output == nil {
		return nil
	}
	err = jsonutil.DecodeJSONFromReader(resp.Body, output)
	if err != nil {
		ec.log.Error("Invalid response from the node", "err", err)
		return err
	}
	return nil
}

func (ec *ExplorerClient) ExplorerCreateDID(peerID string, did string) error {
	ed := ExplorerDID{
		PeerID: peerID,
		DID:    did,
	}
	var er ExplorerResponse
	err := ec.SendExploerJSONRequest("POST", ExplorerCreateDIDAPI, &ed, &er)
	if err != nil {
		return err
	}
	if !er.Status {
		ec.log.Error("Failed to update explorer", "msg", er.Message)
		return fmt.Errorf("failed to update explorer")
	}
	return nil
}

func (ec *ExplorerClient) ExplorerMapDID(oldDid string, newDID string, peerID string) error {
	ed := ExplorerMapDID{
		OldDID: oldDid,
		NewDID: newDID,
		PeerID: peerID,
	}
	var er ExplorerResponse
	err := ec.SendExploerJSONRequest("POST", ExplorerMapDIDAPI, &ed, &er)
	if err != nil {
		return err
	}
	if !er.Status {
		ec.log.Error("Failed to update explorer", "msg", er.Message)
		return fmt.Errorf("failed to update explorer")
	}
	return nil
}

func (ec *ExplorerClient) ExplorerTransaction(et *ExplorerTrans) error {
	var er ExplorerResponse
	err := ec.SendExploerJSONRequest("POST", ExplorerTransactionAPI, et, &er)
	if err != nil {
		return err
	}
	if !er.Status {
		ec.log.Error("Failed to update explorer", "msg", er.Message)
		return fmt.Errorf("failed to update explorer")
	}
	return nil
}