// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Julien Vehent jvehent@mozilla.com [:ulfr]
package worker

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net"
	"time"
)

func InitMQ(uri, certPath, keyPath, caPath string) (amqpConn *amqp.Connection, amqpChan *amqp.Channel, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("worker.initMQ() -> %v", e)
		}
	}()
	// create an AMQP configuration with a 10min heartbeat and timeout
	timeout, _ := time.ParseDuration("600s")
	var dialConfig amqp.Config
	dialConfig.Heartbeat = timeout
	dialConfig.Dial = func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, timeout)
	}
	// if amqps, create the TLS configuration
	if len(uri) > 5 && uri[0:5] == "amqps" {
		// import the client certificates
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			panic(err)
		}
		// import the ca cert
		data, err := ioutil.ReadFile(caPath)
		ca := x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM(data); !ok {
			panic("failed to import CA Certificate")
		}
		TLSconfig := tls.Config{Certificates: []tls.Certificate{cert},
			RootCAs:            ca,
			InsecureSkipVerify: false,
			Rand:               rand.Reader}
		dialConfig.TLSClientConfig = &TLSconfig
	}
	// Setup the AMQP broker connection
	amqpConn, err = amqp.DialConfig(uri, dialConfig)
	if err != nil {
		panic(err)
	}
	amqpChan, err = amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	return
}
