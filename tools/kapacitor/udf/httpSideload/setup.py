from setuptools import setup
from os import path

import sys
setup_dir = path.dirname(__file__)
sys.path.insert(0, setup_dir)

from kapacitor.udf import VERSION

setup(name='httpSideload',
    version=VERSION,
    packages=[
        'kapacitor',
        'kapacitor.udf',
    ],
    install_requires=[
        "protobuf==3.4.0",
        "cachetools",
        "requests"
    ],
    maintainer_email="pawel@kontakt.io",
    license="MIT",
    url="github.com/influxdata/kapacitor",
    description="Kapacitor UDF Agent library",
)
