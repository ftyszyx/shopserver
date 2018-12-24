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

def update(path):
    print"\n"
    print "begin update "+path
    os.chdir( path )
    ret = os.system( "git pull" )
    if ret != 0:
        print "build "+path+" Error!"
        return
    
    print "update "+path+" success!"
    print"\n"
    os.chdir( curpath )


if __name__=="__main__":
   update(shopserver)
   update(shopclient)
   update(erp)
   print "all done!"
   raw_input()
   
