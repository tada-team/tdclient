package tdclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/pion/webrtc/v2"
	"github.com/tada-team/tdproto"
)

func NewPeerConnection(login string, iceServer string) (peerConnection *webrtc.PeerConnection, offer webrtc.SessionDescription, outputTrack *webrtc.Track, err error) {
	peerConnection, err = webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{
			URLs: []string{
				iceServer,
			},
		}},
	})
	if err != nil {
		return nil, offer, nil, fmt.Errorf("%v: NewPeerConnection fail: %v", login, err)
	}

	mediaEngine := webrtc.MediaEngine{}
	mediaEngine.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	// Add codecs
	audioCodecs := mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeAudio)
	if len(audioCodecs) == 0 {
		return nil, offer, nil, fmt.Errorf("%v: offer contained no video codecs", login)
	}
	outputTrack, err = peerConnection.NewTrack(audioCodecs[0].PayloadType, rand.Uint32(), "audio", "pion")
	if err != nil {
		return nil, offer, nil, err
	}
	if _, err = peerConnection.AddTrack(outputTrack); err != nil {
		return nil, offer, nil, err
	}

	offer, err = peerConnection.CreateOffer(nil)
	if err != nil {
		return nil, offer, nil, fmt.Errorf("%v: CreateOffer fail: %v", login, err)
	}

	err = mediaEngine.PopulateFromSDP(offer)
	if err != nil {
		return nil, offer, nil, err
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		return nil, offer, nil, fmt.Errorf("%v: SetLocalDescription fail: %v", login, err)
	}

	//write output if program "hear" something
	peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
		log.Printf("%v: got new track, id: %v \n", login, track.ID())
	})

	return peerConnection, offer, outputTrack, nil
}

func GetIceServer(host string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, host+"/features.json", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new http request: %v", err)
	}

	cli := http.Client{
		Timeout: time.Second * 2,
	}
	res, getErr := cli.Do(req)
	if getErr != nil {
		return "", fmt.Errorf("failed to get features.json: %v", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", fmt.Errorf("failed to read request body: %v", readErr)
	}

	features := tdproto.Features{}
	jsonErr := json.Unmarshal(body, &features)
	if jsonErr != nil {
		return "", fmt.Errorf("failed to unmarshal features json: %v", jsonErr)
	}

	if len(features.ICEServers) > 0 {
		return features.ICEServers[0].Urls, nil
	}

	return "", nil
}

func SendCallOffer(c *WsSession, userName string, callJid *tdproto.JID, sdp string) (res webrtc.SessionDescription, err error) {
	callOffer := new(tdproto.ClientCallOffer)
	callOffer.Name = callOffer.GetName()
	callOffer.Params.Jid = *callJid
	callOffer.Params.Trickle = false
	callOffer.Params.Sdp = sdp
	c.Send(callOffer)
	log.Printf("%v: OFFER: %v\n", userName, callOffer.String())

	callAnswer := new(tdproto.ServerCallAnswer)
	err = c.WaitFor(callAnswer)
	if err != nil {
		return res, fmt.Errorf("%v: server.call.answer fail: %v", userName, err)
	}
	log.Printf("%v: ANSWER: %v\n", userName, callAnswer.String())

	return webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  callAnswer.Params.JSEP.SDP,
	}, nil
}

func SendCallLeave(c *WsSession, userName string, callJid *tdproto.JID) {
	callLeave := new(tdproto.ClientCallLeave)
	callLeave.Name = callLeave.GetName()
	callLeave.Params.Jid = *callJid
	callLeave.Params.Reason = ""
	c.Send(callLeave)

	serverLeaveAnswer := new(tdproto.ServerCallLeave)
	if err := c.WaitFor(serverLeaveAnswer); err != nil {
		log.Println(fmt.Sprintf("%v: server.call.leave fail: %v", userName, err))
	} else {
		log.Printf("%v: LEAVE: %v\n", userName, serverLeaveAnswer.String())
	}
}
