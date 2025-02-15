//
// Copyright IBM Corporation 2020,2021
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
//

package pubsub

import (
	"context"
	"strings"

	"github.com/IBM/kar/core/internal/config"
	"github.com/IBM/kar/core/internal/store"
)

func placementKeyPrefix(t string) string {
	return "pubsub" + config.Separator + "placement" + config.Separator + t
}

func placementKey(t, id string) string {
	return "pubsub" + config.Separator + "placement" + config.Separator + t + config.Separator + id
}

// GetSidecar returns the current sidecar for the given actor type and id or "" if none.
func GetSidecar(ctx context.Context, t, id string) (string, error) {
	s, err := store.Get(ctx, placementKey(t, id))
	if err == store.ErrNil {
		return "", nil
	}
	return s, err
}

// CompareAndSetSidecar atomically updates the sidecar for the given actor type and id.
// Use old = "" to atomically set the initial placement.
// Use new = "" to atomically delete the current placement.
// Returns 0 if unsuccessful, 1 if successful.
func CompareAndSetSidecar(ctx context.Context, t, id, old, new string) (int, error) {
	o := &old
	if old == "" {
		o = nil
	}
	n := &new
	if new == "" {
		n = nil
	}
	return store.CompareAndSet(ctx, placementKey(t, id), o, n)
}

// GetAllActorInstances returns a mapping from actor types to instanceIDs
func GetAllActorInstances(ctx context.Context, actorTypePrefix string) (map[string][]string, error) {
	m := map[string][]string{}
	reply, err := store.Keys(ctx, placementKeyPrefix(actorTypePrefix)+"*")
	if err != nil {
		return nil, err
	}
	for _, key := range reply {
		splitKeys := strings.Split(key, config.Separator)
		actorType := splitKeys[2]
		instanceID := splitKeys[3]
		if m[actorType] == nil {
			m[actorType] = make([]string, 0)
		}
		m[actorType] = append(m[actorType], instanceID)
	}
	return m, nil
}
