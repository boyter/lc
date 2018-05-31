#!/usr/local/bin/python
# -*- coding: utf-8 -*-
# SPDX-License-Identifier: MIT

import sys
import os
import commands
import re
import time


def check_lc(targets):
    correct = 0

    for target in targets:
        data = commands.getstatusoutput('lc -f csv %s' % (target))
        split = data[1].split(',')

        if split[11] == split[13]:
            correct += 1
    return correct


def check_license_detector(targets):
    correct = 0

    for target in targets:
        data = commands.getstatusoutput('license-detector %s' % (target))
        license = target.split('/')[2]

        if 'no license file was found' not in data[1]:
            if license == data[1].split('\n')[1].split('\t')[2]:
                correct += 1

    return correct


def check_identify_license(targets):
    correct = 0

    for target in targets:
        data = commands.getstatusoutput('identify_license %s/LICENSE.txt' % (target))
        license = target.split('/')[2]

        if license == data[1].split('\n')[2].split(' ')[1]:
            correct += 1

    return correct


def check_askalono(targets):
    correct = 0

    for target in targets:
        data = commands.getstatusoutput('askalono id %s/LICENSE.txt' % (target))
        license = target.split('/')[2]

        if license == data[1].split('\n')[0].split(' ')[1]:
            correct += 1

    return correct


if __name__ == '__main__':
    targets = ['./accuracy/' + x + '/' for x in os.listdir('./accuracy')][:1]

    print 'count::%s' % (len(targets))
    print 'checking::lc'
    
    start = time.time()
    correct = check_lc(targets)
    print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)

    start = time.time()
    print 'checking::license-detector'
    correct = check_license_detector(targets)
    print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)

    start = time.time()
    print 'checking::identify_license'
    correct = check_identify_license(targets)
    print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)

    start = time.time()
    print 'checking::askalono'
    correct = check_askalono(targets)
    print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)
