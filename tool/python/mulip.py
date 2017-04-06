# !/usr/bin/python3.4
# -*-coding:utf-8-*-
# on 2016/11/5.
# 功能:
#   生成IP给Linux配置文件用

# IPADDR=146.148.149.202
# PREFIX=24
# IPADDR1=146.148.149.203
# PREFIX1=24
# IPADDR2=146.148.149.204
# PREFIX2=24
# IPADDR3=146.148.149.205
# PREFIX3=24
# IPADDR4=146.148.150.194
# PREFIX4=24
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
                if ii == 0:
                    a = "IPADDR="
                    b = "PREFIX="
                else:
                    a = "IPADDR" + str(ii) + "="
                    b = "PREFIX" + str(ii) + "="

                dudu.append(a + ipprefix + "." + str(i))
                dudu.append(b + "24")
                ii = ii + 1
    except Exception as e:
        print(e)
        pass
    for i in dudu:
        print(i)
