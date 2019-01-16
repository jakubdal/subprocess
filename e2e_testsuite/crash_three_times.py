#!/usr/bin/env python
import os.path, sys

for i in range(3):
    checked_file = os.path.join('crash_test', str(i))
    if not os.path.exists(checked_file):
        open(checked_file, "w+").close()
            exit(1)

sys.stdout.write('success')
