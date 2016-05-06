//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"encoding/json"
)

// Dummy operator that simply wraps an item channel.
type Channel struct {
	readonly
}

func NewChannel() *Channel {
	return &Channel{}
}

func (this *Channel) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitChannel(this)
}

func (this *Channel) New() Operator {
	return &Channel{}
}

func (this *Channel) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{"#operator": "Channel"}
	return json.Marshal(r)
}

func (this *Channel) UnmarshalJSON([]byte) error {
	// NOP: Channel has no data structure
	return nil
}