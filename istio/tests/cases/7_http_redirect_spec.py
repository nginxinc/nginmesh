import requests
import configuration
import performance
from mamba import description, context, it
from expects import expect, be_true, have_length, equal, be_a, have_property, be_none

rule_name="route-rule-http-redirect.yaml"
Rule=configuration.Rule()

with description('Testing HTTP Redirect'):
    with before.all:
         #Read Config file
         configuration.setenv(self)

    with context('Set environment'):
         with it('Add routing rule'):
            Rule.add(rule_name)

    with context('Starting test'):
        with it('Testing HTTP Redirect'):

            while self.total_count < self.request_count:
                r = requests.get(self.url,allow_redirects=False)
                r.status_code
                expect(r.status_code).to(equal(301))
                self.total_count += 1
            configuration.generate_request(self,rule_name)

    with context('Clean environment'):
        with it('Delete routing rule'):
            Rule.delete(rule_name)
