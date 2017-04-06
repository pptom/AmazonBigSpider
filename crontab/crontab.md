0 */2 * * * /sbin/ntpdate time.windows.com
10 0 * * * killall go
10 0 * * * ps -ef|grep /tmp/go-build |awk '{print $2}'|xargs -i kill {}
20 0 * * * /data/app/redis-3.2.1/src/redis-cli -a smart2016 flushall
30 0 * * * mysqldump -uroot -psmart2016 -q smartdb $(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/usa$(date -d yesterday +\%Y\%m\%d).sql.gz
50 0 * * * mysqldump -uroot -psmart2016 -q smartdb Asin$(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/usaasin$(date -d yesterday +\%Y\%m\%d).sql.gz
10 1 * * * mysqldump -uroot -psmart2016 -q jp_smartdb $(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/jp$(date -d yesterday +\%Y\%m\%d).sql.gz
30 1 * * * mysqldump -uroot -psmart2016 -q jp_smartdb Asin$(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/jpasin$(date -d yesterday +\%Y\%m\%d).sql.gz
10 1 * * * mysqldump -uroot -psmart2016 -q de_smartdb $(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/de$(date -d yesterday +\%Y\%m\%d).sql.gz
30 1 * * * mysqldump -uroot -psmart2016 -q de_smartdb Asin$(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/deasin$(date -d yesterday +\%Y\%m\%d).sql.gz
10 1 * * * mysqldump -uroot -psmart2016 -q uk_smartdb $(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/uk$(date -d yesterday +\%Y\%m\%d).sql.gz
30 1 * * * mysqldump -uroot -psmart2016 -q uk_smartdb Asin$(date -d yesterday +\%Y\%m\%d) > /data/www/web/go/src/github.com/hunterhug/AmazonBigSpiderWeb/file/sql/ukasin$(date -d yesterday +\%Y\%m\%d).sql.gz
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/usa/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/jp/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/de/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
0 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/ip.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
5 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/urlpool.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
10 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/helpspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
20 2 * * * /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/crontab/uk/asinspider.sh  >> /data/www/web/go/src/github.com/hunterhug/AmazonBigSpider/haha.log 2>&1 &
