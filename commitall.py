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

shopserver="G:\mywork\go\src\github.com\zyx\shop_server"
shopclient="G:\mywork\shop\shop_client"
erp="G:\mywork\erp"
curpath=os.getcwd();#当前文件目录

def commit(path):
    os.chdir( path )
    print"\n"
    print "begin commit "+path
    print "git add!"
    ret = os.system( "git add -A" )
    if ret != 0:
        print "git add Error!"
        return
    print "git commit!"
    cmd = "git commit -m\" %s\"" % ("zyx")
    ret = os.system( cmd)
    if ret != 0:
        print "git commit Error!"
        return
    print "git push!"
    ret = os.system( "git push" )
    if ret != 0:
        print "git push Error!"
        return
    print"\n"
    os.chdir( curpath )


if __name__=="__main__":
   commit(shopserver)
   commit(shopclient)
   commit(erp)
   print "all done!"
   raw_input()
