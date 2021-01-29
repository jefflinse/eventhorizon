// Copyright (c) 2021 - The Event Horizon authors.
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

package kafka

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/looplab/eventhorizon/eventbus"
)

func TestEventBusIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Connect to localhost if not running inside docker
	broker := os.Getenv("KAFKA_BROKER_HOST")
	if broker == "" {
		broker = "localhost:9092"
	}

	// Get a random app ID.
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	appID := "app-" + hex.EncodeToString(b)

	bus1, err := NewEventBus(broker, appID)
	if err != nil {
		t.Fatal("there should be no error:", err)
	}

	bus2, err := NewEventBus(broker, appID)
	if err != nil {
		t.Fatal("there should be no error:", err)
	}

	// Need to wait here for the topic to be created.
	time.Sleep(3 * time.Second)

	// NOTE: Long timeout needed to accommodate for topic creation times.
	eventbus.AcceptanceTest(t, bus1, bus2, 10*time.Second)
}

func TestEventBusLoadtest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Connect to localhost if not running inside docker
	broker := os.Getenv("KAFKA_BROKER_HOST")
	if broker == "" {
		broker = "localhost:9092"
	}

	// Get a random app ID.
	bts := make([]byte, 8)
	if _, err := rand.Read(bts); err != nil {
		t.Fatal(err)
	}
	appID := "app-" + hex.EncodeToString(bts)

	bus, err := NewEventBus(broker, appID)
	if err != nil {
		t.Fatal("there should be no error:", err)
	}

	// Need to wait here for the topic to be created.
	time.Sleep(3 * time.Second)

	eventbus.LoadTest(t, bus)
}

func BenchmarkEventBus(b *testing.B) {
	// Connect to localhost if not running inside docker
	broker := os.Getenv("KAFKA_BROKER_HOST")
	if broker == "" {
		broker = "localhost:9092"
	}

	// Get a random app ID.
	bts := make([]byte, 8)
	if _, err := rand.Read(bts); err != nil {
		b.Fatal(err)
	}
	appID := "app-" + hex.EncodeToString(bts)

	bus, err := NewEventBus(broker, appID)
	if err != nil {
		b.Fatal("there should be no error:", err)
	}

	// Need to wait here for the topic to be created.
	time.Sleep(3 * time.Second)

	eventbus.Benchmark(b, bus)
}
