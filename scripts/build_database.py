#!/usr/local/bin/python
# -*- coding: utf-8 -*-

import os
import json
import re
import sys
from os import listdir
from os.path import isfile, join

'''
Parses based on the SPDX uploaded licenses in github https://github.com/spdx/license-list-data

Running this takes a long time so be prepared to wait while it churns away
'''

def clean_text(text):
    text = text.lower()
    return text


def find_ngrams(input_list, n):
    return zip(*[input_list[i:] for i in range(n)])


def load_database():
    with open('database.json', 'r') as file:
        database = file.read()

    licenses = json.loads(database)

    for license in licenses:
        license['clean'] = clean_text(license['text'])
        ngrams = []

        ngramrange = [3, 7, 8]

        if license['shortname'] in ['Artistic-1.0', 'BSD-3-Clause']:
            ngramrange = range(2, 35)

        for x in ngramrange:
            ngrams = ngrams + find_ngrams(license['clean'].split(), x)
        license['ngrams'] = ngrams

    return licenses


def build_database():
    license_dir = './license-list-data/json/details/'

    onlyfiles = [f for f in listdir(license_dir) if isfile(join(license_dir, f))]

    licenses = []

    for license in onlyfiles:
        with open(join(license_dir, license), 'r') as file:
            temp = file.read()
            license_json = json.loads(temp)

            ngrams = []
            ngramrange = [3, 7, 8]

            if license_json['licenseId'] in ['Artistic-1.0', 'BSD-3-Clause']:
                ngramrange = range(2, 35)

            for x in ngramrange:
                ngrams = ngrams + find_ngrams(license_json['licenseText'].split(), x)
            license_json['ngrams'] = ngrams

            licenses.append(license_json)

    fair_source = {
        'name': 'Fair Source License v0.9',
        'licenseId': 'Fair-Source-0.9',
        'licenseText': 'Fair Source License, version 0.9 Copyright (C) [year] [copyright owner] Licensor: [legal name of licensor] Software: [name software and version if applicable] Use Limitation: [number] users License Grant. Licensor hereby grants to each recipient of the Software (\"you\") a non-exclusive, non-transferable, royalty-free and fully-paid-up license, under all of the Licensors copyright and patent rights, to use, copy, distribute, prepare derivative works of, publicly perform and display the Software, subject to the Use Limitation and the conditions set forth below. Use Limitation. The license granted above allows use by up to the number of users per entity set forth above (the \"Use Limitation\"). For determining the number of users, \"you\" includes all affiliates, meaning legal entities controlling, controlled by, or under common control with you. If you exceed the Use Limitation, your use is subject to payment of Licensors then-current list price for licenses. Conditions. Redistribution in source code or other forms must include a copy of this license document to be provided in a reasonable manner. Any redistribution of the Software is only allowed subject to this license. Trademarks. This license does not grant you any right in the trademarks, service marks, brand names or logos of Licensor. DISCLAIMER. THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OR CONDITION, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. LICENSORS HEREBY DISCLAIM ALL LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE. Termination. If you violate the terms of this license, your rights will terminate automatically and will not be reinstated without the prior written consent of Licensor. Any such termination will not affect the right of others who may have received copies of the Software from you.',
    }
    for x in [3, 7, 8]:
        ngrams = ngrams + find_ngrams(fair_source['licenseText'].split(), x)
        fair_source['ngrams'] = ngrams

    licenses.append(fair_source)
    return licenses


if __name__ == '__main__':

    licenses = build_database()
    
    for license in licenses:
        matches = []

        for ngram in license['ngrams']:
            find = ' '.join(ngram)
            ismatch = True

            filtered = [x for x in licenses if x['licenseId'] != license['licenseId']]
            for lic in filtered:
                if find in lic['licenseText']:
                    ismatch = False
                    break

            if ismatch:
                matches.append(find)

        if len(matches) == 0:
            print '>>>>', license['licenseId'], len(matches)
        else:
            print license['licenseId'], len(matches)

        license['keywords'] = matches

    licenses = [{
        'licenseText': x['licenseText'],
        'name': x['name'],
        'licenseId': x['licenseId'],
        'keywords': x['keywords'][:50]
    } for x in licenses]

    with open('database_keywords.json', 'w') as myfile:
        myfile.write(json.dumps(licenses))
