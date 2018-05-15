import requests
import time
import configuration
from mamba import description, context, it
from expects import expect, be_true, have_length, equal, be_a, have_property, be_none

headers = {'content-type': 'application/json','accept': 'application/json'}
with description('Zipkin tracing functionality'):
    with before.all:
         #Read Config file
         configuration.setenv(self)

    with context('Deploy Zipkin and make sure port forwarded to localhost'):
        with it('Bookinfo Zipkin tracing feature'):
            for _ in range(10):
                r = requests.get(self.url)
                r.status_code
                expect(r.status_code).to(equal(200))
            r1=requests.get(self.zipkin)
            r1.status_code
            expect(r1.status_code).to(equal(200))
            if 'productpage' in r1.text:
                expect(0).to(equal(0))
            else:
                expect(0).not_to(equal(0))
            configuration.generate_request(self)







