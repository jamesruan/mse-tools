// Copyright 2012 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/acolwell/mse-tools/ebml"
	"github.com/acolwell/mse-tools/webm"
	"log"
)

type BlockGroupClient struct {
	parsedBlock bool
	id          uint64
	rawTimecode int64
	flags       uint8
	blockData   []byte
	writer      *ebml.Writer
}

func (c *BlockGroupClient) OnListStart(offset int64, id int) bool {
	//log.Printf("OnListStart(%d, %s)\n", offset, webm.IdToName(id))
	log.Printf("OnListStart() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *BlockGroupClient) OnListEnd(offset int64, id int) bool {
	//log.Printf("OnListEnd(%d, %s)\n", offset, webm.IdToName(id))
	log.Printf("OnListEnd() : Unexpected element %s\n", webm.IdToName(id))
	return false
}

func (c *BlockGroupClient) OnBinary(id int, value []byte) bool {
	if id == webm.IdBlock {
		blockInfo := webm.ParseSimpleBlock(value)
		c.id = blockInfo.Id
		c.rawTimecode = int64(blockInfo.Timecode)
		c.flags = blockInfo.Flags & 0x0f
		c.blockData = value[blockInfo.HeaderSize:]
		c.parsedBlock = true
		return true
	} else if id == webm.IdBlockAdditions {
		c.writer.Write(id, value)
		return true
	}
	log.Printf("OnBinary() : Unexpected element %s size %d\n", webm.IdToName(id), len(value))
	return false
}

func (c *BlockGroupClient) OnInt(id int, value int64) bool {
	if id == webm.IdDiscardPadding || id == webm.IdReferenceBlock {
		c.writer.Write(id, value)
		return true
	}
	log.Printf("OnInt() : Unexpected element %s %d\n", webm.IdToName(id), value)
	return false
}

func (c *BlockGroupClient) OnUint(id int, value uint64) bool {
	if id == webm.IdBlockDuration {
		c.writer.Write(id, value)
		return true
	}
	log.Printf("OnUint() : Unexpected element %s %u\n", webm.IdToName(id), value)
	return false
}

func (c *BlockGroupClient) OnFloat(id int, value float64) bool {
	log.Printf("OnFloat() : Unexpected element %s %f\n", webm.IdToName(id), value)
	return false
}

func (c *BlockGroupClient) OnString(id int, value string) bool {
	log.Printf("OnString() : Unexpected element %s %s\n", webm.IdToName(id), value)
	return false
}
