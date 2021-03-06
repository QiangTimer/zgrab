/*
 * ZGrab Copyright 2015 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 */

package zlib

import (
	"encoding/json"
	"errors"
	"net"
	"time"
)

type Grab struct {
	Host   net.IP            `json:"host"`
	Domain string            `json:"domain"`
	Time   time.Time         `json:"timestamp"`
	Log    []ConnectionEvent `json:"log"`
}

type ConnectionEvent struct {
	Data  EventData
	Error error
}

type encodedGrab struct {
	Host   string            `json:"host"`
	Domain *string           `json:"domain"`
	Time   string            `json:"time"`
	Log    []ConnectionEvent `json:"log"`
}

type encodedConnectionEvent struct {
	Type  EventType `json:"type"`
	Data  EventData `json:"data"`
	Error *string   `json:"error"`
}

type partialConnectionEvent struct {
	Data  EventData `json:"data"`
	Error *string   `json:"error"`
}

func (ce *ConnectionEvent) MarshalJSON() ([]byte, error) {
	t := ce.Data.GetType()
	var esp *string
	if ce.Error != nil {
		es := ce.Error.Error()
		esp = &es
	}
	obj := encodedConnectionEvent{
		Type:  t,
		Data:  ce.Data,
		Error: esp,
	}
	return json.Marshal(obj)
}

func (ce *ConnectionEvent) UnmarshalJSON(b []byte) error {
	ece := new(partialConnectionEvent)
	tn := struct {
		TypeName string `json:"type"`
	}{}
	if err := json.Unmarshal(b, &tn); err != nil {
		return err
	}
	t, typeErr := EventTypeFromName(tn.TypeName)
	if typeErr != nil {
		return typeErr
	}
	ece.Data = t.GetEmptyInstance()
	if err := json.Unmarshal(b, &ece); err != nil {
		return err
	}
	ce.Data = ece.Data
	if ece.Error != nil {
		ce.Error = errors.New(*ece.Error)
	}
	return nil
}

func (g *Grab) MarshalJSON() ([]byte, error) {
	var domainPtr *string
	if g.Domain != "" {
		domainPtr = &g.Domain
	}
	time := g.Time.Format(time.RFC3339)
	obj := encodedGrab{
		Host:   g.Host.String(),
		Domain: domainPtr,
		Time:   time,
		Log:    g.Log,
	}
	return json.Marshal(obj)
}

func (g *Grab) UnmarshalJSON(b []byte) error {
	eg := new(encodedGrab)
	err := json.Unmarshal(b, eg)
	if err != nil {
		return err
	}
	g.Host = net.ParseIP(eg.Host)
	if eg.Domain != nil {
		g.Domain = *eg.Domain
	}
	if g.Time, err = time.Parse(time.RFC3339, eg.Time); err != nil {
		return err
	}
	g.Log = eg.Log
	return nil
}

func (g *Grab) status() status {
	if len(g.Log) == 0 {
		return status_failure
	}
	for _, entry := range g.Log {
		if entry.Error != nil {
			return status_failure
		}
	}
	return status_success
}
