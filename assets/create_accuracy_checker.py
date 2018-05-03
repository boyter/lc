#!/usr/local/bin/python
# -*- coding: utf-8 -*-
# SPDX-License-Identifier: MIT

import os
import json
import re
import sys
from os import listdir
from os.path import isfile, join
import codecs


def make_sure_path_exists(path):
    try:
        os.makedirs(path)
    except OSError as exception:
        pass


def read_database():
    license_json = {}
    with open('../database_keywords.json', 'r') as file:
        temp = file.read()
        license_json = json.loads(temp)
    return license_json


def create_folders(licenses={}):
    for license in licenses:
        lid = license['licenseId']
        text = license['licenseText']

        directory = './accuracy/' + lid + '/'

        make_sure_path_exists(directory)
        file = codecs.open(directory + 'LICENSE.txt', 'w', 'utf-8')
        file.write(text)
        file.close()


if __name__ == '__main__':
    licenses = read_database()
    create_folders(licenses)
