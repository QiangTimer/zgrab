/*
 * ZGrab Copyright 2015 Regents of the University of Michigan
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy
 * of the License at http://www.apache.org/licenses/LICENSE-2.0
 */

package zlib

type FTPBannerEvent struct {
	Banner string `json:"banner",omitempty`
}

var FTPBannerEventType = EventType{
	TypeName:         CONNECTION_EVENT_FTP,
	GetEmptyInstance: func() EventData { return new(FTPBannerEvent) },
}

func (f *FTPBannerEvent) GetType() EventType {
	return FTPBannerEventType
}
