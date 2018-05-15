import requests
import configuration
import performance
from mamba import description, context, it
from expects import expect, be_true, have_length, equal, be_a, have_property, be_none

rule_name="route-rule-reviews-test-v2.yaml"
Rule=configuration.Rule()

with description('Testing route "jason" user to V2'):
    with before.all:
         #Read Config file
         configuration.setenv(self)

    with context('Set environment'):
         with it('Add routing rule'):
               Rule.add(rule_name)

    with context('Starting test'):
        with it('Testing route "jason" user to V2'):
            while self.total_count < self.request_count:
                cookie={'user':'jason'}
                r = requests.get('http://'+configuration.GATEWAY_URL+'/productpage', cookies =cookie )
                expect(r.status_code).to(equal(200))
                if 'color="black"' not in r.text and 'color="red"' not in r.text:
                    self.total_count += 1
                    self.v1_count+=1
                elif 'color="black"' in r.text:
                    self.total_count += 1
                    self.v2_count+=1
                elif 'color="red"' in r.text:
                    self.total_count += 1
                    self.v3_count+=1
                else:
                     self.total_count += 1

            print(" | V1 Hit="+str(self.v1_count)+" | V2 Hit="+str(self.v2_count)+" | V3 Hit="+str(self.v3_count)+" | Total Hit="+str(self.total_count)+ " |")
            expect(self.v1_count).to(equal(0))
            expect(self.v2_count).not_to(equal(0))
            expect(self.v3_count).to(equal(0))
            configuration.generate_request(self,rule_name)

    with context('Clean environment'):
        with it('Delete routing rule'):
              Rule.delete(rule_name)
