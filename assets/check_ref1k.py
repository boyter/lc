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
    return split[11]


def check_license_detector(target):
    data = commands.getstatusoutput('license-detector %s' % (target))

    if 'no license file was found' not in data[1]:
        return data[1].split('\n')[1].split('\t')[2]
    return "UNKNOWN"


def check_identify_license(target):
    data = commands.getstatusoutput('identify_license %s' % (target))
    return data[1].split('\n')[2].split(' ')[1]


def check_askalono(target):
    data = commands.getstatusoutput('askalono id %s' % (target))
    return data[1].split('\n')[0].split(' ')[1]


if __name__ == '__main__':
    targets = ['./ref1k/' + x + '/' for x in os.listdir('./ref1k')]

    print 'targetfile,lc,license-detector,identify_license,askalono'
    for target in targets:
        license = [x for x in os.listdir(target) if 'license' in x.lower()]

        if len(license) == 0:
            print '%s,NA,NA,NA,NA' % (target)
        else:
            # All except license-detector support single files
            targetFile = target + license[0]

            guesses = {
                'lc': check_lc(targetFile),
                'license-detector': check_license_detector(target),
                'identify_license': check_identify_license(targetFile),
                'askalono': check_askalono(targetFile),
            }

            print '%s,%s,%s,%s,%s' % (
                targetFile, 
                guesses['lc'],
                guesses['license-detector'],
                guesses['identify_license'],
                guesses['askalono'],
            )
