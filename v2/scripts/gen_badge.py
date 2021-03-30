#!/usr/bin/env python

import sys
import requests


def to_num(value: str) -> float:
    return float(value.strip('%'))


def request_badge(coverage: str) -> None:
    try:
        coverage_num = to_num(coverage)
        color = 'red' if coverage_num <= 85 else 'brightgreen'
        url = f"https://img.shields.io/badge/coverage-{coverage_num}%25-{color}"
        response = requests.request('GET', url)
        with open('badge.svg', 'w') as f:
            f.write(response.content.decode('utf-8'))

    except Exception as e:
        print('Error')
        print(e)


if __name__ == '__main__':
    request_badge(sys.argv[1])
