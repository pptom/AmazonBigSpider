#!/bin/sh
ps -ef|grep go
ps -ef|grep /tmp/go-build |awk '{print $2}'|xargs -i kill {}
killall go
ps -ef|grep go