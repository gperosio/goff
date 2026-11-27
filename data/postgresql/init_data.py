#!/usr/bin/env python3

import os
import sys 

def main():
    with open(os.path.join(os.path.dirname(__file__), 'init_data.sql'), 'r') as f:
        sql = f.read()
    print(sql)