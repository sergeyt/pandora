#!/usr/bin/env python
from time import sleep
from tasks import index_file

alice = 'https://www.adobe.com/be_en/active-use/pdf/Alice_in_Wonderland.pdf'
result = index_file.delay(alice)

while not result.ready():
    sleep(1)

print(result.get())
