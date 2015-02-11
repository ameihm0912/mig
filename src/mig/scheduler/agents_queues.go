// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Julien Vehent jvehent@mozilla.com [:ulfr]
package main

import (
	"encoding/json"
	"fmt"
	"mig"
	"time"

	"github.com/streadway/amqp"
)

// startHeartbeatsListener initializes the routine that receives heartbeats from agents
func startHeartbeatsListener(ctx Context) (heartbeatChan <-chan amqp.Delivery, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("startHeartbeatsListener() -> %v", e)
		}
		ctx.Channels.Log <- mig.Log{OpID: ctx.OpID, Desc: "leaving startHeartbeatsListener()"}.Debug()
	}()

	_, err = ctx.MQ.Chan.QueueDeclare("mig.agt.heartbeats", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = ctx.MQ.Chan.QueueBind("mig.agt.heartbeats", "mig.agt.heartbeats", "mig", false, nil)
	if err != nil {
		panic(err)
	}

	err = ctx.MQ.Chan.Qos(0, 0, false)
	if err != nil {
		panic(err)
	}

	heartbeatChan, err = ctx.MQ.Chan.Consume("mig.agt.heartbeats", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	ctx.Channels.Log <- mig.Log{OpID: ctx.OpID, Desc: "agents heartbeats listener initialized"}

	return
}

// getHeartbeats processes the heartbeat messages sent by agents
func getHeartbeats(msg amqp.Delivery, ctx Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("getHeartbeats() -> %v", e)
		}
		if ctx.Debug.Heartbeats {
			ctx.Channels.Log <- mig.Log{OpID: ctx.OpID, Desc: "leaving getHeartbeats()"}.Debug()
		}
	}()

	var agt mig.Agent
	err = json.Unmarshal(msg.Body, &agt)
	if err != nil {
		panic(err)
	}
	if ctx.Debug.Heartbeats {
		desc := fmt.Sprintf("Received heartbeat for Agent '%s' QueueLoc '%s'", agt.Name, agt.QueueLoc)
		ctx.Channels.Log <- mig.Log{Desc: desc}.Debug()
	}
	// discard expired heartbeats
	agtTimeOut, err := time.ParseDuration(ctx.Agent.TimeOut)
	if err != nil {
		panic(err)
	}
	expirationDate := time.Now().Add(-agtTimeOut)
	if agt.HeartBeatTS.Before(expirationDate) {
		desc := fmt.Sprintf("Expired heartbeat received from Agent '%s'", agt.Name)
		ctx.Channels.Log <- mig.Log{Desc: desc}.Notice()
		return
	}
	// if agent is not authorized, ack the message and skip the registration
	// nothing is returned to the agent. it's simply ignored.
	ok, err := isAgentAuthorized(agt.QueueLoc, ctx)
	if err != nil {
		panic(err)
	}
	if !ok {
		desc := fmt.Sprintf("getHeartbeats(): Agent '%s' is not authorized", agt.QueueLoc)
		ctx.Channels.Log <- mig.Log{Desc: desc}.Warning()
		// send an event to notify workers of the failed agent auth
		err = sendEvent("mig.events.agent.authentication.failed", msg.Body, ctx)
		if err != nil {
			panic(err)
		}
		// agent authorization failed so we drop this heartbeat and return
		return
	}

	// write to database in a goroutine to avoid blocking
	go func() {
		// if an agent already exists in database, we update it, otherwise we insert it
		agent, err := ctx.DB.AgentByQueueAndPID(agt.QueueLoc, agt.PID)
		if err != nil {
			agt.DestructionTime = time.Date(9998, time.January, 11, 11, 11, 11, 11, time.UTC)
			agt.Status = mig.AgtStatusOnline
			// create a new agent
			err = ctx.DB.InsertAgent(agt)
			if err != nil {
				ctx.Channels.Log <- mig.Log{Desc: fmt.Sprintf("Heartbeat DB insertion failed with error '%v' for agent '%s'", err, agt.Name)}.Err()
			}
		} else {
			// the agent exists in database. reuse the existing ID, and keep the status if it was
			// previously set to destroyed or upgraded. otherwise set status to online
			agt.ID = agent.ID
			if agent.Status == mig.AgtStatusDestroyed || agent.Status == mig.AgtStatusUpgraded {
				agt.Status = agent.Status
			} else {
				agt.Status = mig.AgtStatusOnline
			}
			err = ctx.DB.UpdateAgentHeartbeat(agt)
			if err != nil {
				ctx.Channels.Log <- mig.Log{Desc: fmt.Sprintf("Heartbeat DB update failed with error '%v' for agent '%s'", err, agt.Name)}.Err()
			}
			// if the agent that exists in the database has a status of 'destroyed'
			// we should not be received a heartbeat from it. so, if detectmultiagents
			// is set in the scheduler configuration, we pass the agent queue over to the
			// routine than handles the destruction of agents
			if agent.Status == mig.AgtStatusDestroyed && ctx.Agent.DetectMultiAgents {
				ctx.Channels.DetectDupAgents <- agent.QueueLoc
			}
		}
	}()

	return
}

// startResultsListener initializes the routine that receives heartbeats from agents
func startResultsListener(ctx Context) (resultsChan <-chan amqp.Delivery, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("startResultsListener() -> %v", e)
		}
		ctx.Channels.Log <- mig.Log{OpID: ctx.OpID, Desc: "leaving startResultsListener()"}.Debug()
	}()

	_, err = ctx.MQ.Chan.QueueDeclare("mig.agt.results", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = ctx.MQ.Chan.QueueBind("mig.agt.results", "mig.agt.results", "mig", false, nil)
	if err != nil {
		panic(err)
	}

	err = ctx.MQ.Chan.Qos(0, 0, false)
	if err != nil {
		panic(err)
	}

	resultsChan, err = ctx.MQ.Chan.Consume("mig.agt.results", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	ctx.Channels.Log <- mig.Log{OpID: ctx.OpID, Desc: "agents results listener initialized"}

	return
}
