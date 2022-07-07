package modules

import (
	"encoding/json"
	"fmt"
	"strings"
	_ "strings"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/oxodao/metaprint/utils"
)

type NetUsage struct {
    Prefix   string
    Suffix   string
    Format   string
    Device   string
    Rounding int
    Unit     string
}

func (r NetUsage) Print(args []string) string {
    devs, err := linuxproc.ReadNetworkStat("/proc/net/dev")
    if err != nil {
        panic(fmt.Errorf("proc/net/dev read error %s", err))
    }
    if r.Device == "" {
        panic(fmt.Errorf("No network device specified"))
    }

    devindex := 0
    for i, v := range devs {
        if v.Iface != r.Device {
            continue
        }

        devindex = i
        break
    }

    elapsed := 0.0
    onceavg := linuxproc.NetworkStat{}
    twoceavg := linuxproc.NetworkStat{}
    timer := time.Now()
    rxSpeed, txSpeed := utils.NewHSpeed(0), utils.NewHSpeed(0)
    for {
        defer func ()  {
            if r:=recover(); r!= nil {
                fmt.Println(fmt.Errorf("err %s", r))
                panic(r)
            }
        }()
        
        if onceavg.Iface!="" {
            twoceavg = r.readDev(devindex) 
            elapsed = (time.Now().Sub(timer)).Seconds()
            deltaRx := twoceavg.RxBytes - onceavg.RxBytes
            deltaTx := twoceavg.TxBytes - onceavg.TxBytes
            rxps := float64(deltaRx) / elapsed
            txps := float64(deltaTx) / elapsed
            // dump.Println(twoceavg, elapsed, deltaRx, deltaTx, rxps, txps)

            rxSpeed = utils.NewHSpeed(rxps)
            txSpeed = utils.NewHSpeed(txps)
            // dump.V(rxSpeed)
            // dump.V(txSpeed)

            fmt.Println(utils.ReplaceVariables(r.Format, map[string]interface{}{
                "rx": rxSpeed.String(),
                "tx": txSpeed.String(),
            }))
        }

        timer = time.Now()
        onceavg = r.readDev(devindex)
        time.Sleep(1 * time.Second)
    }

}

func (r NetUsage) readDev(ind int) linuxproc.NetworkStat {
    devs, err := linuxproc.ReadNetworkStat("/proc/net/dev")
    if err != nil {
        panic(fmt.Errorf("proc/net/dev read error %s", err))
    }
    if ind < 0 || ind >= len(devs) {
        panic(fmt.Errorf("cannot find device %s at index %d", r.Device, ind))
    }

    return devs[ind]
}

type NetStats struct {
    iface        string
    prevStats    linuxproc.NetworkStat
    currentStats linuxproc.NetworkStat
    deltaTime    time.Duration
}

func (r NetUsage) getAvg(v linuxproc.NetworkStat) (float64, float64) {
    stats := make(map[string]interface{})
    netstats, ok := json.Marshal(v)
    if ok != nil {
        fmt.Println(fmt.Errorf("json marshall err %s", ok))
        return 0, 0
    }
    var x map[string]interface{}
    _ = json.Unmarshal(netstats, &x)

    // if x["iface"] == "lo" {
    //    continue
    // }
    iface := strings.ReplaceAll(x["iface"].(string), "-", "_")
    for s, val := range x {
        key := fmt.Sprintf("%s_%s", iface, s)
        stats[key] = val
    }

    idletime := float64(0)
    cputime := float64(0)
    for n, proc := range x {
        if n == "id" {
            continue
        }
        if n == "idle" {
            idletime = proc.(float64)
        }
        cputime = cputime + proc.(float64)
    }
    return cputime, idletime
}

func (r NetUsage) GetPrefix() string {
    return r.Prefix
}

func (r NetUsage) GetSuffix() string {
    return r.Suffix
}
