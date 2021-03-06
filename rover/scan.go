package rover

import (
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/patrickmn/go-cache"
    "log"
    "time"
    "math/rand"
)

type Probe struct{
     Address string
     Strength int64
     Timestamp time.Time
     Location string
}

func Scan(detected chan<- Probe, device string, granularity time.Duration, location string, threshold int64, sampleRate float64){
    snapshotLen := int32(1024)
    promiscuous := false
    timeout := time.Duration(30) * time.Second
    var err error
    var handle *pcap.Handle
    c := cache.New(granularity, time.Minute)
 
    // Open device
    handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
    if err != nil {log.Fatal(err) }
    defer handle.Close()

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packetSource.Packets() {
        if probe, ok := getProbe(packet); ok {
	   if _, found := c.Get(probe.Address); !found {
	      if rand.Float64() <= sampleRate && probe.Strength < threshold {
	      	      c.Set(probe.Address, true, cache.DefaultExpiration)
	      	      probe.Timestamp = time.Now()
	      	      probe.Location = location
	      	      detected <- probe
	      }
           }
	}
    }
}

func getProbe(packet gopacket.Packet) (p Probe, valid bool){
    // Get Packet and check for Probe Req
    wlanLayer := packet.Layer(layers.LayerTypeDot11)
    if wlanLayer != nil {
        wlanPacket, _ := wlanLayer.(*layers.Dot11)
	if wlanPacket.Type.String() == "MgmtProbeReq" {
           p.Address = wlanPacket.Address2.String()
	   if tap := packet.Layer(layers.LayerTypeRadioTap); tap != nil {
	      dot11r, _ := tap.(*layers.RadioTap)
	      p.Strength = int64(-1 * dot11r.DBMAntennaSignal)
	      valid = true 
	   }
	}
    }
    return
}