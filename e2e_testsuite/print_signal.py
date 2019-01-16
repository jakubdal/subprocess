#!/usr/bin/env python
import signal, os, sys

def signal_handler(sig, frame):
    sys.stdout.write("SIGINT called")

signal.signal(signal.SIGINT, signal_handler)
signal.pause()
