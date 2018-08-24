# dl

[![Build Status](https://travis-ci.org/ejunjsh/dl.svg?branch=master)](https://travis-ci.org/ejunjsh/dl)

a concurrent http file downloader,support rate limit, resume from break-point.

# install

    go get github.com/ejunjsh/dl

# usage

    # dl
    usage: dl [[rate limit]:url...]
    rate limit: limit the speed,unit is KB
    url...: urls you want to download

# example

## concurrent download

    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    172.06MB/607.13MB(28.34%)[=========>                         ] 4m 51s (1.49MB/s)
    143.66MB/607.13MB(23.66%)[=======>                         ] 9m 37s (821.30KB/s)
    139.26MB/607.13MB(22.94%)[=======>                         ] 9m 27s (844.47KB/s)

## rate limit

below example shows the download speed that is limited in 200KB

    ➜ dl 200:https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    172.06MB/607.13MB(28.34%)[=========>                         ] 4m 51s (198KB/s)

## resume from break-point

below shows two commands,the second command resume from the first command

    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    314.74MB/607.13MB(51.84%)[=================>               ] 6m 15s (798.37KB/s)

    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    319.59MB/607.13MB(52.64%)[===================>                 ] 0s (318.32MB/s)

