package main

import (
	"fmt"
	"math"
	"time"

	"github.com/ninjasphere/go-castv2/controllers"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/channels"
	"github.com/ninjasphere/go-ninja/devices"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/kodi_jsonrpc"
)

type MediaPlayer struct {
	player   *devices.MediaPlayerDevice
	client   kodi_jsonrpc.Connection
	receiver *controllers.ReceiverController
	media    *controllers.MediaController
}

type VolumeNotification struct {
	Data struct {
		Muted  bool    `json:"muted"`
		Volume float64 `json:"volume"`
	} `json:"data"`
	Sender string `json:"sender"`
}

func NewMediaPlayer(driver ninja.Driver, conn *ninja.Connection, id string, client kodi_jsonrpc.Connection) (*MediaPlayer, error) {

	device := &MediaPlayer{
		client: client,
	}

	var result map[string]interface{}
	err := device.call("JSONRPC.Version", nil, &result, time.Second*5)

	if err != nil {
		log.Fatalf("Couldn't get version from %s: %s", id, err)
	}

	versionData := result["version"].(map[string]interface{})
	version := fmt.Sprintf("%.f.%.f.%.f", versionData["major"].(float64), versionData["minor"].(float64), versionData["patch"].(float64))

	log.Infof("Connected to %s. Version: %s", id, version)

	player, err := devices.CreateMediaPlayerDevice(driver, &model.Device{
		NaturalID:     id,
		NaturalIDType: "mdns",
		Name:          &id,
		Signatures: &map[string]string{
			"ninja:manufacturer": "Kodi",
			"ninja:productName":  "Kodi",
			"ninja:thingType":    "mediaplayer",
			"kodi:version":       version,
		},
	}, conn)

	if err != nil {
		return nil, err
	}

	device.player = player

	player.ApplyVolume = device.applyVolume
	if err := player.EnableVolumeChannel(true); err != nil {
		player.Log().Fatalf("Failed to enable volume channel: %s", err)
	}

	player.ApplyPlayPause = device.applyPlayPause
	if err := player.EnableControlChannel([]string{"playing", "paused", "stopped", "buffering", "busy", "idle", "inactive"}); err != nil {
		player.Log().Fatalf("Failed to enable control channel: %s", err)
	}

	/*toggle := true
	value := 0.0
	go func() {
		for {
			toggle = !toggle
			value = value + 0.05

			err := device.applyPlayPause(toggle)

			//err := device.applyVolume(&channels.VolumeState{
			//	Level: &value,
			//	Muted: &toggle,
			//})

			if err != nil {
				log.Warningf("Failed to set play/pause: %s", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()*/

	go func() {
		for notification := range client.Notifications {

			switch notification.Method {
			case "Player.OnPlay":
				player.UpdateControlState(channels.MediaControlEventPlaying)
			case "Player.OnPause":
				player.UpdateControlState(channels.MediaControlEventPaused)
			case "Player.OnStop":
				player.UpdateControlState(channels.MediaControlEventStopped)
			case "Application.OnVolumeChanged":
				var volume VolumeNotification
				err := notification.Read(&volume)
				if err != nil {
					player.Log().Warningf("Failed to read volume notification: %s", err)
				} else {

					vol := float64(volume.Data.Volume) / 100

					err = player.UpdateVolumeState(&channels.VolumeState{
						Level: &vol,
						Muted: &volume.Data.Muted,
					})

					if err != nil {
						player.Log().Warningf("Failed to update volume from notification: %s", err)
					}

				}
			}

			//spew.Dump("notification", notification)
		}
	}()

	return device, nil
}

func (d *MediaPlayer) getPlayerId() (int, error) {
	var result []interface{}
	err := d.call("Player.GetActivePlayers", nil, &result, time.Second*5)

	if err != nil {
		return 0, fmt.Errorf("Failed to get active player id: %s", err)
	}

	if len(result) == 0 {
		return 0, fmt.Errorf("The player is not active")
	}

	return int(result[0].(map[string]interface{})["playerid"].(float64)), nil
}

func (d *MediaPlayer) applyPlayPause(play bool) error {

	playerId, err := d.getPlayerId()
	if err != nil {
		return err
	}

	var result map[string]interface{}

	err = d.call("Player.PlayPause", map[string]interface{}{
		"playerid": playerId,
		"play":     play,
	}, &result, time.Second*5)

	return err
}

func (d *MediaPlayer) call(method string, params map[string]interface{}, result interface{}, timeout time.Duration) error {
	var response kodi_jsonrpc.Response
	if params == nil {
		response = d.client.Send(kodi_jsonrpc.Request{Method: method}, true)
	} else {
		response = d.client.Send(kodi_jsonrpc.Request{Method: method, Params: &params}, true)
	}

	return response.Read(result, timeout/time.Second)
}

func (d *MediaPlayer) applyVolume(state *channels.VolumeState) error {
	d.player.Log().Infof("applyVolume called, volume %v", state)

	if state.Level != nil {
		var result float64

		err := d.call("Application.SetVolume", map[string]interface{}{
			"volume": int(math.Min(*state.Level*100, 100)),
		}, &result, time.Second*5)

		if err != nil {
			return err
		}
		d.player.Log().Debugf("SetVolume Response: %v", result)
	}

	if state.Muted != nil {
		var result bool

		err := d.call("Application.SetMute", map[string]interface{}{
			"mute": *state.Muted,
		}, &result, time.Second*5)

		if err != nil {
			return err
		}
		d.player.Log().Debugf("SetMute Response: %v", result)
	}

	return nil
}
