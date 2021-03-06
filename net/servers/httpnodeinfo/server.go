package httpnodeinfo

import (
	"Elastos.ELA/common/config"
	"Elastos.ELA/core/ledger"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"Elastos.ELA/net/servers"
)

type Info struct {
	NodeVersion   string
	BlockHeight   uint32
	NeighborCnt   int
	Neighbors     []NgbNodeInfo
	HttpRestPort  int
	HttpWsPort    int
	HttpJsonPort  int
	HttpLocalPort int
	NodePort      int
	NodeId        string
}

type NgbNodeInfo struct {
	NgbId   string
	NbrAddr string
}

var templates = template.Must(template.New("info").Parse(page))

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var ngbrNodersInfo []NgbNodeInfo
	var node = servers.NodeForServers

	ngbrNoders := node.GetNeighborNoder()

	for i := 0; i < len(ngbrNoders); i++ {
		ngbrNodersInfo = append(ngbrNodersInfo, NgbNodeInfo{
			NgbId:   fmt.Sprintf("0x%x", ngbrNoders[i].GetID()),
			NbrAddr: ngbrNoders[i].GetAddr(),
		})
	}

	pageInfo := &Info{
		BlockHeight:   ledger.DefaultLedger.Blockchain.BlockHeight,
		NeighborCnt:   len(ngbrNoders),
		Neighbors:     ngbrNodersInfo,
		HttpRestPort:  config.Parameters.HttpRestPort,
		HttpWsPort:    config.Parameters.HttpWsPort,
		HttpJsonPort:  config.Parameters.HttpJsonPort,
		NodePort:      config.Parameters.NodePort,
		NodeId:        fmt.Sprintf("0x%x", node.GetID()),
	}

	err := templates.ExecuteTemplate(w, "info", pageInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartServer() {
	http.HandleFunc("/info", viewHandler)
	http.ListenAndServe(":"+strconv.Itoa(int(config.Parameters.HttpInfoPort)), nil)
}
