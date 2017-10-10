# Bouygues SMS

[![Build Status](https://travis-ci.org/tomsquest/bouyguessms.svg?branch=master)](https://travis-ci.org/tomsquest/bouyguessms)

Send up to 5 free SMS per day using Bouygues Telecom unofficial API.  
The program is released as a standalone binary and it can also be used as a GO library.  
The code is written in Golang.  

## Usage

Using from the command line:

```bash
$ bouyguessms \
    -login "yourEmail@domain.com" \
    -pass  "yourPassword" \
    -to "0601010101,0602020202" \
    -msg "Hello World!"
SMS sent successfully at 2017-10-02 07:43:09.051154206 +0200 CEST. Quota left: 4.
```

As a library in your GO program:

```go
import "github.com/tomsquest/bouyguessms"

quota, err := SendSms("yourEmail@domain.com", "yourPassword",
    "Hello Gophers, it's "+time.Now().String(),
    "0601010101")
if err != nil {
    log.Fatalln("error sending sms", err)
}
log.Printf("SMS sent. Quota left: %d.\n", quota)
```

## Inspiration

Original code from [Y3nd](https://github.com/y3nd)'s [bouygues-sms](https://github.com/y3nd/bouygues-sms) written in Javascript. Thanks!