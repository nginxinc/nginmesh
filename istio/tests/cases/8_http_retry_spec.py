import requests
import configuration
import performance
from mamba import description, context, it
from expects import expect, be_true, have_length, equal, be_a, have_property, be_none

rule_name="route-rule-http-retry.yaml"
nginx_conf_path="/etc/istio/proxy/conf.d/http_0.0.0.0_9080.conf"
Rule=configuration.Rule()

with description('Testing HTTP Retry'):
    with before.all:
         #Read Config file
         configuration.setenv(self)

    with context('Set environment'):
         with it('Add routing rule'):
            Rule.add(rule_name)

    with context('Starting test'):
        with it('Testing HTTP Retry'):
            while self.total_count < self.request_count:
                r = requests.get(self.url,allow_redirects=False)
                r.status_code
                expect(r.status_code).to(equal(200))
                self.total_count += 1
                output=configuration.run_shell("kubectl exec -it $(kubectl get pod | grep productpage | awk '{ print $1 }') -c istio-proxy cat "+nginx_conf_path,"check")

            if 'proxy_next_upstream_timeout' in output and 'proxy_next_upstream_tries' in output :
                    print("Total Retry Hit="+str(self.total_count))
                    expect(self.total_count).not_to(equal(0))
            configuration.generate_request(self,rule_name)

    with context('Clean environment'):
        with it('Delete routing rule'):
            Rule.delete(rule_name)


