# dl

[![Build Status](https://travis-ci.org/ejunjsh/dl.svg?branch=master)](https://travis-ci.org/ejunjsh/dl)

a concurrent http file downloader,support rate limit, resume from break-point.

# install

    go get github.com/ejunjsh/dl

# usage

    # dl
    usage: dl [-h <header> [ -h <header>]] [[rate limit:]url...]
    -h: specify your http header,format is "key:value"
    rate limit: limit the speed,unit is KB
    url...: urls you want to download

# example

## concurrent download

    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg http://mirrors.neusoft.edu.cn/centos/7/isos/x86_64/CentOS-7-x86_64-Minimal-1804.iso
      ideaIU-2018.2.1.dmg |607.13MB[>                               ]26m13s|384.02KB/s
      CentOS-7-x86_64-Mini|906.00MB[===>                            ] 3m22s|  3.96MB/s

## rate limit

below example shows the download speed that is limited in 200KB

    ➜ dl 200:https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    ideaIU-2018.2.1.dmg |607.13MB[===>                            ]46m14s|199.34KB/s

## resume from break-point

below shows two commands,the second command resume from the first command

    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    ideaIU-2018.2.1.dmg |607.13MB[====>                           ] 5m 1s|  1.73MB/s
    ^C

    ➜ dl https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg
    ideaIU-2018.2.1.dmg |607.13MB[=====>                          ] 3m17s|  2.57MB/s

## customize header

    dl --header aaa:bbb --header ccc:ddd  https://download.jetbrains.com/idea/ideaIU-2018.2.1.dmg

above download will use the "aaa:bbb;ccc:ddd" as its header

## proxy

support `HTTP_PROXY` or `HTTPS_PROXY` environment parameter to setup proxy.

