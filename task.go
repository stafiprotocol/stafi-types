package types

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	stafiMainnetTypes = "https://raw.githubusercontent.com/stafiprotocol/stafi-types/main/mainnet/stafi.json"
	stafiTestnetTypes = "https://raw.githubusercontent.com/stafiprotocol/stafi-types/main/testnet/stafi.json"
	retryLimit        = 6
	SecondsWait       = 3 * time.Second
	mainnetChain      = "Stafi"
	devChain          = "Development"
	seiyaChain        = "Stafi Testnet Seiya"
)

type Types struct {
	tickerSeconds int64
	stop          chan struct{}
	Endpoint      string
	logger        LoggerInterface
	sarpc         SarpcInterface
	stafiTypesUrl string
	stafiJsonBts  []byte
}

func NewTypes(sarpc SarpcInterface, logger LoggerInterface, tickerSeconds int64, endpoint string) (*Types, error) {
	chain, err := sarpc.GetSystemChain()
	if err != nil {
		return nil, err
	}
	var useUrl string
	switch chain {
	case mainnetChain:
		useUrl = stafiMainnetTypes
	case devChain:
		useUrl = stafiTestnetTypes
	case seiyaChain:
		useUrl = stafiTestnetTypes
	default:
		return nil, fmt.Errorf("not support chain: %s", chain)
	}

	typesBts, err := getStafiTypes(useUrl)
	if err != nil {
		return nil, err
	}

	s := &Types{
		tickerSeconds: tickerSeconds,
		stop:          make(chan struct{}),
		Endpoint:      endpoint,
		stafiTypesUrl: useUrl,
		stafiJsonBts:  typesBts,
		sarpc:         sarpc,
		logger:        logger,
	}
	return s, nil
}

func (task *Types) StartMonitor() {
	SafeGoWithRestart(task.Handler)
}

func (task *Types) Stop() {
	close(task.stop)
}

func (task *Types) GetStafiJsonTypes() []byte {
	return task.stafiJsonBts
}

func (task *Types) Handler() {
	ticker := time.NewTicker(time.Duration(task.tickerSeconds) * time.Second)
	defer ticker.Stop()
out:
	for {

		select {
		case <-task.stop:
			break out
		case <-ticker.C:
			var remoteStafiTypesBts []byte
			var err error
			retry := 0
			for {
				if retry > retryLimit {
					task.logger.Error("getStafiTypes reach retry limit", "err", err)
					break
				}
				remoteStafiTypesBts, err = getStafiTypes(task.stafiTypesUrl)
				if err != nil {
					task.logger.Warn("getStafiTypes failed", "err", err)
					time.Sleep(SecondsWait)
					retry++
					continue
				}
				break
			}

			if !bytes.Equal(task.stafiJsonBts, remoteStafiTypesBts) {
				task.stafiJsonBts = remoteStafiTypesBts
				task.sarpc.RegCustomTypes(task.stafiJsonBts)
				task.logger.Info("got new stafi types, already update")
			}
		}
	}
}

func getStafiTypes(url string) ([]byte, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d", rsp.StatusCode)
	}

	defer rsp.Body.Close()
	rspBts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if len(rspBts) == 0 {
		return nil, fmt.Errorf("rsp body empty")
	}
	return rspBts, nil
}
