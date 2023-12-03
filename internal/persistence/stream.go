// Copyright © 2020 Chris Camel <camel.christophe@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package persistence

import (
	"github.com/AsynkronIT/protoactor-go/actor"

	persistencev1 "github.com/ccamel/playground-protoactor.go/internal/persistence/v1"
)

type EventPredicate func(event *persistencev1.EventRecord) bool

type SubscriptionID string

type EventStreamStore interface {
	// Subscribe does a subscription for events.
	Subscribe(pid *actor.PID, start *string, predicate EventPredicate) SubscriptionID
	Unsubscribe(SubscriptionID)
}
