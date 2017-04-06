0 */2 * * * /sbin/ntpdate time.windows.com
10 0 * * * killall go
10 0 * * * ps -ef|grep /tmp/go-build |awk '{print $2}'|xargs -i kill {}
20 0 * * * /data/app/redis-3.2.1/src/redis-cli -a smart2016 flushall
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/helpspider.sh >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
