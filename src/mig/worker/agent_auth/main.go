// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Julien Vehent jvehent@mozilla.com [:ulfr]
package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"log"
	"mig/worker"
	"os"
)

type Config struct {
	MQ struct {
		Uri                     string
		TLScert, TLSkey, CAcert string
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s - a worker verifying agents that fail to authenticate\n"+
			"Usage: %s -c /etc/mig/agent_auth_worker.cfg\n",
			os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	var err error
	var configPath = flag.String("c", "/etc/mig/agent_auth_worker.cfg", "Load configuration from file")
	flag.Parse()

	var conf Config
	err = gcfg.ReadFileInto(&conf, *configPath)
	if err != nil {
		panic(err)
	}
	_, amqpChan, err := worker.InitMQ(conf.MQ.Uri, conf.MQ.TLScert, conf.MQ.TLSkey, conf.MQ.CAcert)
	if err != nil {
		log.Fatal(err)
	}
	_, err = amqpChan.QueueDeclare("mig.workers.agent.authentication", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	err = amqpChan.QueueBind("mig.workers.agent.authentication", "mig.events.agent.authentication.*", "mig", false, nil)
	if err != nil {
		panic(err)
	}
	err = amqpChan.Qos(0, 0, false)
	if err != nil {
		panic(err)
	}
	evChan, err := amqpChan.Consume("mig.workers.agent.authentication", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	for event := range evChan {
		fmt.Printf("%s\n", event.Body)
	}
	return
}
