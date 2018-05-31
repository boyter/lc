#!/usr/local/bin/python
# -*- coding: utf-8 -*-
# SPDX-License-Identifier: MIT

import sys
import os
import commands
import re
import time


def check_lc(target):
    data = commands.getstatusoutput('lc -f csv %s' % (target))
    split = data[1].split(',')
    return split[13]


def check_license_detector(targets):
    data = commands.getstatusoutput('license-detector %s' % (target))
    license = target.split('/')[2]

    if 'no license file was found' not in data[1]:
        return data[1].split('\n')[1].split('\t')[2]
    return "UNKNOWN"


def check_identify_license(targets):
    data = commands.getstatusoutput('identify_license %s' % (target))
    license = target.split('/')[2]
    return data[1].split('\n')[2].split(' ')[1]


def check_askalono(targets):
    data = commands.getstatusoutput('askalono id %s' % (target))
    license = target.split('/')[2]
    return data[1].split('\n')[0].split(' ')[1]


if __name__ == '__main__':
    targets = ['./ref1k/' + x + '/' for x in os.listdir('./ref1k')][:10]

    for target in targets:
        license = [x for x in os.listdir(target) if 'license' in x.lower()]

        if len(license) == 0:
            print 'No License File For %s' % (target)
        else:
            target += license[0]
            print target
            guesses = {
                'lc': check_lc(target),
                # 'license-detector': check_license_detector(target),
                'identify_license': check_identify_license(target),
                'askalono': check_askalono(target),
            }

            print guesses


    # print 'count::%s' % (len(targets))
    # print 'checking::lc'
    
    # start = time.time()
    # correct = check_lc(targets)
    # print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)

    # start = time.time()
    # print 'checking::license-detector'
    # correct = check_license_detector(targets)
    # print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)

    # start = time.time()
    # print 'checking::identify_license'
    # correct = check_identify_license(targets)
    # print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)

    # start = time.time()
    # print 'checking::askalono'
    # correct = check_askalono(targets)
    # print 'correct:%s::%s percent::time:%s' % (correct, float(correct) / float(len(targets)) * 100, time.time() - start)
