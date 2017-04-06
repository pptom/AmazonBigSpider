# !/usr/bin/python3.4
# -*-coding:utf-8-*-
from selenium import webdriver
from selenium.common.exceptions import NoSuchElementException
from base64 import b64encode
import time

def getFirefox(url, ip, total=1):
    data = ""
    proxy = {
        "host": ip.split(":")[0],
        "port": ip.split(":")[1],
        "user": "smart",
        "pass": "smart2016"
    }
    profile = webdriver.FirefoxProfile()
    # add new header
    profile.add_extension("../modify_headers-0.7.1.1-fx.xpi")
    profile.set_preference("extensions.modify_headers.currentVersion", "0.7.1.1-fx")
    profile.set_preference("modifyheaders.config.active", True)
    profile.set_preference("modifyheaders.headers.count", 1)
    profile.set_preference("modifyheaders.headers.action0", "Add")
    profile.set_preference("modifyheaders.headers.name0", "Proxy-Switch-Ip")
    profile.set_preference("modifyheaders.headers.value0", "yes")
    profile.set_preference("modifyheaders.headers.enabled0", True)

    # add proxy
    profile.set_preference('network.proxy.type', 1)
    profile.set_preference('network.proxy.http', proxy['host'])
    profile.set_preference('network.proxy.http_port', int(proxy['port']))
    profile.set_preference('network.proxy.no_proxies_on', 'localhost, 127.0.0.1')

    # Proxy auto login
    profile.add_extension('../close_proxy_authentication-1.1-sm+tb+fx.xpi')
    credentials = '{user}:{pass}'.format(**proxy)
    credentials = b64encode(credentials.encode('ascii')).decode('utf-8')
    profile.set_preference('extensions.closeproxyauth.authtoken', credentials)
    browser = webdriver.Firefox(profile)
    browser.get(url)
    try:
        if total == 1:
            data = browser.page_source
        else:
            data = browser.find_element_by_xpath("html").text
    except NoSuchElementException:
        print("神马都没有")
    return browser, data


if __name__ == '__main__':
    # url = "http://ip.42.pl"
    url="https://www.amazon.com/dp/B005HX2YT0"
    while True:
        try:
            ip = input("输入需解救的IP：（如146.148.149.206:808）")
            # ip=""
            # # ip="111.13.65.244:80"
            browers, data = getFirefox(url=url, ip=ip)
            print("解救了" + ip + ",暂停5秒后浏览器关闭")
            time.sleep(5)
            browers.close()
        except Exception as e:
            print(e)