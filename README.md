# dl

[![Build Status](https://travis-ci.org/ejunjsh/dl.svg?branch=master)](https://travis-ci.org/ejunjsh/dl)

a concurrent http file downloader,support rate limit

# install

    go get github.com/ejunjsh/dl

# usage

    # dl
    usage: dl [[rate limit]:url...]
    rate limit: limit the speed,unit is KB
    url...: urls you want to download

# example


    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    172.06MB/607.13MB(28.34%)[=========>                         ] 4m 51s (1.49MB/s)
    143.66MB/607.13MB(23.66%)[=======>                         ] 9m 37s (821.30KB/s)
    139.26MB/607.13MB(22.94%)[=======>                         ] 9m 27s (844.47KB/s)
