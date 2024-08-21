# simple-discord

simple-discord is designed to be a "simple" to use framework for interfacing with the Discord API. It takes advantage of the Gateway that Discord uses for requests back and forth between the simple-discord client and the actual API. Currently, the framework is limited in what it can do and is certainly not close to fully integrating with the Discord API.

## Features

- **Golang Performance**: Utilizes Go's performance by running the sending and receiving of data on different channels.
- **Custom JSON Unmarshalling**: Populates an active session with easy-to-use Golang structs.
- **Active Cache**: Maintains an active cache while the bot is running.

## Installation

To run the project, simply install Go v1.22.3 or above, and run the following command to install the latest package into your project:

```sh
go get github.com/Carmen-Shannon/simple-discord@latest
```

## Version
**v0.4.0**
- Technically the first release, because it's the first time I learned how to properly use tags. This version lays the basic groundwork for the rest of the work to come. The basic read/write logic as well as handling events to/from the Discord API is completed.