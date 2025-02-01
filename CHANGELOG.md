## [v0.1.0]
- Re-wrote the session/gateway logic
- Set up basic gateway management with reconnect/resume functionality
- Set up audio playback and AES256 encoding

## [v0.1.1]
- Removed multiple fmt.Println statements to clean up the console, still a lot scattered around but they should stay for now.
- Fixed the `decodePayload` handler in the base session to properly decode from a map of payload types
- Fixed an issue with the UDP handler not properly decoding packets into the correct payload type
- Fixed an issue with the audio player not connecting the UDP session properly, related to the previous issue
- Moved the 

## [v0.2.0]
- Fixed an error related to socket closure from internal socket calls
- Added VoiceSession resume functionality so now voice sessions should be able to stay connected indefinitely, without breaking the udp connection. (hopefully)

## [v0.3.0]
- Stabilized reconnect logic for client and voice sessions so they should now be able to keep uptime indefinitely.