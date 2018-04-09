from __future__ import absolute_import
import time
import performance
import configuration
import logo
from mamba import description, context, it
from expects import expect, be_true, have_length, equal, be_a, have_property, be_none

with description('Testing basic functionality'):
    with before.all:
         #Read Config file
         configuration.setenv(self)

    with context('Starting test'):
        with it('Testing basic functionality'):
            configuration.generate_request(self)

            expect(self.v1_count).not_to(equal(0))
            expect(self.v2_count).not_to(equal(0))
            expect(self.v3_count).not_to(equal(0))


