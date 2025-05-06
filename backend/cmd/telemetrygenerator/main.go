package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"time"
	"turion-takehome/internal/config"
	tdp "turion-takehome/internal/turiondatapacket"
)

func main() {
	config, err := config.NewTelemetryGeneratorConfig()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.Dial("udp", config.GroundStationEmulatorAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	packetCount := uint16(0)
	for {
		data := createTelemetryPacket(&packetCount)
		_, err := conn.Write(data)
		if err != nil {
			log.Printf("Error sending telemetry: %v", err)
			continue
		}
		if packetCount%5 == 0 {
			log.Printf("Sent anomalous telemetry packet #%d\n", packetCount)
		} else {
			log.Printf("Sent normal telemetry packet #%d\n", packetCount)
		}
		time.Sleep(1 * time.Second)
		packetCount++
	}
}
func createTelemetryPacket(seqCount *uint16) []byte {
	buf := new(bytes.Buffer)
	// Create primary header
	// PacketID: Version(3) | Type(1) | SecHdrFlag(1) | APID(11)
	packetID := uint16(tdp.PACKET_VERSION)<<13 |
		uint16(tdp.PACKET_TYPE)<<12 |
		uint16(tdp.SEC_HDR_FLAG)<<11 |
		uint16(tdp.APID)

	// PacketSeqCtrl: SeqFlags(2) | SeqCount(14)
	packetSeqCtrl := uint16(tdp.SEQ_FLAGS)<<14 | (*seqCount & 0x3FFF)
	// Generate telemetry data
	payload := generateTelemetryPayload(*seqCount%5 == 0)
	// Calculate total packet length (excluding primary header first 6 bytes)
	packetDataLength := uint16(binary.Size(tdp.CCSDSSecondaryHeader{}) +
		binary.Size(tdp.TelemetryPayload{}) - 1)

	primaryHeader := tdp.CCSDSPrimaryHeader{
		PacketID:      packetID,
		PacketSeqCtrl: packetSeqCtrl,
		PacketLength:  packetDataLength,
	}
	// Create secondary header
	secondaryHeader := tdp.CCSDSSecondaryHeader{
		Timestamp: uint64(time.Now().Unix()),

		SubsystemID: tdp.SUBSYSTEM_ID,
	}
	// Write headers and payload
	binary.Write(buf, binary.BigEndian, primaryHeader) // CCSDS uses big-endian
	binary.Write(buf, binary.BigEndian, secondaryHeader)
	binary.Write(buf, binary.BigEndian, payload)
	return buf.Bytes()
}

func generateTelemetryPayload(generateAnomaly bool) tdp.TelemetryPayload {
	if generateAnomaly {
		// Randomly choose one parameter to be anomalous
		anomalyType := rand.Intn(4)
		switch anomalyType {
		case 0:
			return tdp.TelemetryPayload{
				Temperature: randomFloat(35.0, 40.0), // High temperature anomaly
				Battery:     randomFloat(70.0, 100.0),
				Altitude:    randomFloat(500.0, 550.0),
				Signal:      randomFloat(-60.0, -40.0),
			}
		case 1:
			return tdp.TelemetryPayload{
				Temperature: randomFloat(20.0, 30.0),
				Battery:     randomFloat(20.0, 40.0), // Low battery anomaly
				Altitude:    randomFloat(500.0, 550.0),
				Signal:      randomFloat(-60.0, -40.0),
			}
		case 2:
			return tdp.TelemetryPayload{
				Temperature: randomFloat(20.0, 30.0),
				Battery:     randomFloat(70.0, 100.0),
				Altitude:    randomFloat(300.0, 400.0), // Low altitude anomaly
				Signal:      randomFloat(-60.0, -40.0),
			}
		default:
			return tdp.TelemetryPayload{
				Temperature: randomFloat(20.0, 30.0),
				Battery:     randomFloat(70.0, 100.0),
				Altitude:    randomFloat(500.0, 550.0),
				Signal:      randomFloat(-90.0, -80.0), // Weak signal anomaly
			}
		}
	}
	return tdp.TelemetryPayload{
		Temperature: randomFloat(20.0, 30.0),   // Normal operating range
		Battery:     randomFloat(70.0, 100.0),  // Battery percentage
		Altitude:    randomFloat(500.0, 550.0), // Orbit altitude
		Signal:      randomFloat(-60.0, -40.0), // Signal strength

	}
}
func randomFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
