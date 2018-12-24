#! /usr/bin/python2.7
# coding=utf-8
# vim:expandtab:ts=4:sw=4:
# 
import time
import re
import os
import sys
import shutil
import platform

curpath=os.getcwd();#当前文件目录


def build():
    ret = os.system( "go build" )
    if ret != 0:
        print "build   Error!"
        return False
    print "build  success!"
    return True
   
def kill(grepstr):
    cmd = "ps -aux|grep \"%s\" |grep -v grep|awk '{print $2}'" % (grepstr)
    print cmd
    ret = os.popen( cmd)
    pid=ret.readline().strip() 
    print "get pid:"+pid
    ret.close()
    if pid!="":
        ret = os.system( "kill "+pid )
        if ret != 0:
            print "kill   Error!"
            return False
        print "kill   success"
        return True
    else:
        print "pid not exit"
        return True

def run(cmd):
    ret = os.system(cmd)
    if ret != 0:
        print "run   Error!"
        return
    print "run  success!"


if __name__=="__main__":   
    conf=sys.argv[1].strip()
    if len(sys.argv)>=3:
        param=sys.argv[2].strip()
        if param=="restart":
            print("restart "+conf)
            if kill("shop_server -conf "+conf+".conf")==True:
                run("nohup ./shop_server -conf "+conf+".conf &")
        elif param=="stop":
            print("stop "+conf)
            if kill("shop_server -conf "+conf+".conf")==True:
                print("stop ok")
    else:
        print "build and run "+conf
        if build()==True:
            if kill("shop_server -conf "+conf+".conf")==True:
                run("nohup ./shop_server -conf "+conf+".conf &")
   
