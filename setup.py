#!/usr/bin/env python
from setuptools import find_packages, setup

setup(
    name="RSD-Updater",
    version="2.0.0",
    description="Updates the Red Eclipse server docker images.",
    author="Iceflower S",
    author_email="iceflower@gmx.de",
    license='GPLv3',
    url="https://github.com/IceflowRE/redeclipse-server-docker",
    classifiers=[
        'Programming Language :: Python :: 3.7',
        'License :: OSI Approved :: GNU General Public License v3 (GPLv3)',
        'Development Status :: 5 - Production/Stable',
        'Operating System :: OS Independent',
        'Intended Audience :: Developers',
        'Natural Language :: English',
        'Environment :: Console',
    ],
    packages=find_packages(include=['rsdupdater', 'rsdupdater.*']),
    python_requires='>=3.7',
    install_requires=[

    ],
    extras_require={
        'dev': [
            'setuptools==45.1.0',
        ],
    },
    zip_safe=True,
    entry_points={
        'console_scripts': [
            'rsd-updater = rsdupdater.main:main',
        ],
    },
)
