// Copyright 2024 Team 3256. All Rights Reserved.
// Author: patrick@warriorb.org (Patrick Woollvin)
//
// Methods for interfacing with the field SMCs (Scoring and Management Computers) over websocket.
//
// Heavy "inspiration" from the cheesy-arena plc package

package smc_ws

import (
	"github.com/Team3256/warrior-arena/websocket"
	"time"
)

type Smc interface {
	SetAddress(address string)
	Connect()
	GetConnectionStatus() bool
	GetIsEnabled() bool
	IsHealthy() bool
	lastHeartbeat() time.Time
	Run()
	GetFieldEStop() bool
	GetTeamEStops() [3]bool
	GetTeamAStops() [3]bool
	GetButtonsConnected() [3]bool
	ResetMatch()
	SetStackLights(red, blue, orange, green bool)
	SetStackBuzzer(state bool)
	SetFieldResetLight(state bool)
	GetCycleState(max, index, duration int) bool
	GetInputNames() []string
	GetInputIds() []string
	GetAmpButtons() (bool, bool)
	GetAmpNoteCounts() int
	GetSpeakerNoteCounts() int
	SetSpeakerMotors(state bool)
	SetSpeakerLights(redState, blueState bool)
	SetSubwooferCountdown(redState, blueState bool)
	SetAmpLights(low, high, coop bool)
	SetPostMatchSubwooferLights(state bool)
}

type SmcWS struct {
	address          string
	RedDSWebsocket   websocket.Websocket
	BlueDSWebsocket  websocket.Websocket
	Aux1Websocket    websocket.Websocket // For 2024: Aux1 is the Red Amp
	Aux2Websocket    websocket.Websocket // For 2024: Aux1 is the Blue Amp
	Aux3Websocket    websocket.Websocket // General use
	Aux4Websocket    websocket.Websocket // General use
	Aux3Enabled      bool
	Aux4Enabled      bool
	isHealthy        bool
	inputs           [inputCount]bool
	registers        [registerCount]uint16
	oldInputs        [inputCount]bool
	oldRegisters     [registerCount]uint16
	oldCoils         [coilCount]bool
	cycleCounter     int
	matchResetCycles int
	scoringEquipment bool
}

// Discrete inputs
//
//go:generate stringer -type=input
type input int

const (
	fieldEStop input = iota
	red1EStop
	red1AStop
	red2EStop
	red2AStop
	red3EStop
	red3AStop
	blue1EStop
	blue1AStop
	blue2EStop
	blue2AStop
	blue3EStop
	blue3AStop
	redConnected1
	redConnected2
	redConnected3
	blueConnected1
	blueConnected2
	blueConnected3
	redAmplify
	redCoop
	blueAmplify
	blueCoop
	inputCount
)

// 16-bit registers
//
//go:generate stringer -type=register
type register int

const (
	fieldIoConnection register = iota
	redSpeaker
	blueSpeaker
	redAmp
	blueAmp
	miscounts
	registerCount
)

// Coils
//
//go:generate stringer -type=coil
type coil int

const (
	heartbeat coil = iota
	matchReset
	stackLightGreen
	stackLightOrange
	stackLightRed
	stackLightBlue
	stackLightBuzzer
	fieldResetLight
	speakerMotors
	redSpeakerLight
	blueSpeakerLight
	redSubwooferCountdown
	blueSubwooferCountdown
	redAmpLightLow
	redAmpLightHigh
	redAmpLightCoop
	blueAmpLightLow
	blueAmpLightHigh
	blueAmpLightCoop
	postMatchSubwooferLights
	coilCount
)

func (ws *SmcWS) Run() {
	for {
		if plc.handler == nil {
			if !plc.IsEnabled() {
				// No PLC is configured; just allow the loop to continue to simulate inputs and outputs.
				plc.isHealthy = false
			} else {
				err := plc.connect()
				if err != nil {
					log.Printf("PLC error: %v", err)
					time.Sleep(time.Second * plcRetryIntevalSec)
					plc.isHealthy = false
					continue
				}
			}
		}

		startTime := time.Now()
		plc.update()
		time.Sleep(time.Until(startTime.Add(time.Millisecond * plcLoopPeriodMs)))
	}
}
