package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	warehouseLocation = "warehouse"
	pitLocation       = "pit"
	// mineShatfLocation = "mineshaft"

	pocketSize  = 1
	storageSize = 5

	miningSpeed       = 2 * time.Second
	movementSpeed     = 3 * time.Second
	unloadPocketSpeed = 1 * time.Second
)

type warehouse struct {
	storage []*coal
}

func (w *warehouse) isFull() bool {
	return len(w.storage) == storageSize
}

func (w *warehouse) getSize() int {
	return len(w.storage)
}

type coal struct{}

type pit struct {
	minerCounter uint
}

// type mineshaft struct {
// 	direction    string
// 	minerCounter uint
// }

type miner struct {
	pocket   *coal
	location string
}

func (m *miner) gotoPit() {
	time.Sleep(movementSpeed)
	m.location = pitLocation
}

func (m *miner) gotoWarehouse() {
	time.Sleep(movementSpeed)
	m.location = warehouseLocation
}

func (m *miner) whereami() string {
	return m.location
}

func (m *miner) mineCoal() *coal {
	time.Sleep(miningSpeed)
	return &coal{}
}

func (m *miner) grab(coal *coal) bool {
	m.pocket = coal
	return m.pocket != nil
}

func (m *miner) emptyPocket() {
	m.pocket = nil
}

func (m *miner) isPocketFull() bool {
	return m.pocket != nil
}

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	w := &warehouse{}
	m := &miner{}

	for {
		logger.Info().Int("storage_size", w.getSize()).Bool("pocket_full", m.isPocketFull()).Msg("storage size")
		if w.isFull() {
			logger.Fatal().Msg("storage is full")
		}

		logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("going to the pit")
		m.gotoPit()
		logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("i'm in pit")

		var coal *coal
		if m.location == pitLocation {
			logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("mining coal")
			coal = m.mineCoal()
			logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("mining is finshed")
		}

		if coal != nil && m.grab(coal) {
			logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("going to warehouse")
			m.gotoWarehouse()
			logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("i'm in warehouse")
		}

		if m.location == warehouseLocation {
			logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("put coal to the storage")
			w.storage = append(w.storage, m.pocket)
			m.emptyPocket()
			logger.Info().Str("location", m.whereami()).Bool("pocket_full", m.isPocketFull()).Msg("done, my pocket is empty")
		}
	}
}
