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

## [v0.4.0]
- Updating Go version to 1.23.5
- Fixed a garbage collection issue with the ffmpeg package, causing all binaries to be read into memory spiking memory usage.
- Fixed an issue where the GC was not releasing resources properly back to the system.
- Cut memory consumption by about 50% from 190MiB to >100MiB during ffmpeg processing. (still some work to be done here)
  - as an additional note to this change, I should be able to address the issue of memory not releasing back to the OS at some point
- Fixed multiple bugs related to the `AudioPlayer` interface.
- Re-wrote the function to encode audio to Opus, separating it into two branches
  - One branch will process the static audio file with ffmpeg, sending PCM frames to a channel.
  - The other branch will listen for these frames on the channel, and process them to Opus encoded audio frames, sending them to another channel.
- Added XChaCha-Poly1305 encryption
- Removed ffprobe binaries, ignoring metadata for now since it wasn't being used anyways.

## [v0.5.0]
- Fixing the tagged version.

## [v0.5.1]
- Fixing an issue where the program would exit when attempting to set the `version` of the bot using the latest git tags, this obviously doesn't work if the project has no git tags.
- Updating README to reflect some changes and better document how to init the bot.

## [v0.5.2]
- Fixing documentation.

## [v0.6.0]
- Added mp4 to mp3 conversion with ffmpeg

## [v0.6.1]
- Fixed mp4 to mp3 conversion, still some work to be done here I think as it's kind of inneficient to convert from mp3 to mp4

## [v0.6.2]
- Added a global and local rate limiter for all requests that use `HttpRequest` from `general_requests` in the request_util package

## [v0.6.3]
- Fixed the `Interaction.Data` for `InteractionCreateEvent` structs.
- Removed a TODO from `request_util`.

## [v0.6.4]
- Fixed the `ClientSession.Play` call not properly exiting the audio player if ffmpeg fails.