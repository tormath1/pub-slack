### Pub Slack

Pub Slack is a simple tool in order to be notified on Slack, when you have an alert on your Google Cloud Platform budget. You won't have bad suprise at the end of the month. 
I prefer to run this kind of service with cloud functions, but I prefer to code in Go, so while we're waiting to have a Go runtime environment, I made this simple script.

![slack screenshot](https://raw.githubusercontent.com/tormath1/pub-slack/master/img/image1.png)

### Before

You need to create an application in Slack in order to be able to post on a channel and get the slack token: https://api.slack.com/incoming-webhooks
You're also need to create a Google Service account with at least a readonly permission for pub/sub and create a subscription for this script.

### Build from sources

```shell
$ go version
go version go1.10.3 linux/amd64
$ uname -vr
4.18.1-arch1-1-ARCH #1 SMP PREEMPT Wed Aug 15 21:11:55 UTC 2018
$ git clone https://github.com/tormath1/pub-slack
$ cd pub-slack
$ go build -o pub-slack main.go 
```

### Parameters

```shell
$ ./pub-slack --help
Usage of ./pub-slack:
  -credentials string
    	absolute path to the Google Credentials JSON file
  -project string
    	ID of your Google Project where the topic is created
  -slackURL string
    	Incoming webhook URL for Slack
  -subscription string
    	name of the subscription
```