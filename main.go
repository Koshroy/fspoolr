package main

import (
    "log"
    //"time"
    "github.com/Koshroy/fspoolr/spoolr"
    "github.com/Koshroy/fspoolr/statemanage"
    "github.com/Koshroy/fspoolr/fmonitor"
    "github.com/Koshroy/httplistener"
)

const SETTINGS_FNAME = "settings.json"

func main() {
    fm, err := fmonitor.NewFilemonitor()
    if err != nil {
        log.Fatalln("could not open file monitor")
    }

    gs := statemanage.NewGlobalState()

    mySettings, err := spoolr.NewSettings("settings.json")
    if err != nil {
        log.Fatalln("could not open settings.json")
    }

    log.Println("dirs", mySettings.Dirs)
    
    for _, dir := range mySettings.Dirs {
        a, err := spoolr.NewJsonArtifact(dir)
        if err != nil {
            log.Println("could not open json artifact on dir", dir)
            log.Println(err, "\n")
            continue
        }
        fm.AddArtifact(a)
        gs.AddArtifact(a)
    }

    fm.SetEventChan(gs.EventChan())
    fm.Start()
    gs.Start()

    // for {
    //     time.Sleep(300 * time.Millisecond)
    // }

    h := httplistener.NewHTTPListener(mySettings.Port, mySettings.AppRoot)

    h.SetResChannel(gs.ResourceChan())
    h.SetReqChannel(gs.RequestChan())
    h.Start()

    // err = watcher.Watch("/Users/koushik/tmp")
    // if err != nil {
    //     log.Fatal(err)
    // }

    // <-done

    // /* ... do stuff ... */
    // watcher.Close()
}
