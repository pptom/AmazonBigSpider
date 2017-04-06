# !/usr/bin/python3.4
# -*-coding:utf-8-*-
# on 2016/11/5.
# 功能:
#   生成IP给Squid配置文件用

# acl ip1 localip 192.168.1.2
# acl ip2 localip 192.168.1.3
# acl ip3 localip 192.168.1.4
# tcp_outgoing_address 192.168.1.2 ip1
# tcp_outgoing_address 192.168.1.3 ip2
# tcp_outgoing_address 192.168.1.4 ip3

if __name__ == "__main__":
    ii = 0
    dudu = []
    try:
        while True:
            ips = input("如：146.148.149.202-254:")
            temp = ips.split("-")
            ipend = int(temp[1])

            temptemp = temp[0].split(".")

            ipprefix = ".".join(temptemp[0:3])
            ipbegin = int(temptemp[3])
            for i in range(ipbegin, ipend + 1):
                dudu.append(ipprefix + "." + str(i))
                ii = ii + 1
    except Exception as e:
        print(e)
        pass
    for i in dudu:
        # acl ip3 localip 192.168.1.4
        # tcp_outgoing_address 192.168.1.2 ip1
        print("acl ip" + i + " localip " + i)
        print("tcp_outgoing_address " + i + " " + "ip" + i)
