import subprocess
import grequests
import requests
import performance
import time
import os


rule_apply_time=5
nginmesh_rule_path="../../release/samples/bookinfo/kube/"
nginmesh_install_path="../../release/install/"
count_init=0
request_count=10
performance_status='on'
performance_thread='1'
performance_connection='10'
performance_duration='1s'
app_namespace=os.environ.get('app_namespace','default')

def run_shell(self,type):
    if type=="check":
        return str(subprocess.check_output(self, universal_newlines=True,shell=True)).rstrip()
    elif type=="run":
        subprocess.call(self+" > /dev/null 2>&1 | exit 0",universal_newlines=True,shell=True)
        time.sleep(rule_apply_time)
        return

GATEWAY_URL =run_shell("kubectl get svc -n istio-system | grep -E 'istio-ingress' | awk '{ print $4 }'","check")

def setenv(self):
    self.url = "http://"+GATEWAY_URL+"/productpage"
    self.zipkin="http://localhost:9411/api/v2/services"
    self.performance=performance_status
    self.v1_count=count_init
    self.v2_count=count_init
    self.v3_count=count_init
    self.total_count=count_init
    self.request_count=request_count
    return self.performance,self.v1_count,self.v2_count,self.v3_count,self.request_count,self.total_count

def generate_request(self, rule_name=None):
    if rule_name !="route-rule-reviews-test-v2.yaml" and rule_name !="route-rule-http-redirect.yaml" and rule_name !="route-rule-http-retry.yaml" :
        urls = [self.url for i in range(self.request_count)]
        rs = (grequests.get(self.url,allow_redirects=False) for url in urls)
        results = grequests.map(rs)
        for r in results:
            if r.status_code==200 and 'color="black"' not in r.text and 'color="red"' not in r.text:
                self.total_count += 1
                self.v1_count+=1
            elif r.status_code==200 and 'color="black"' in r.text:
                self.total_count += 1
                self.v2_count+=1
            elif r.status_code==200 and 'color="red"' in r.text:
                self.total_count += 1
                self.v3_count+=1
            else:
                self.total_count += 1
        print(" | V1 Hit="+str(self.v1_count)+" | V2 Hit="+str(self.v2_count)+" | V3 Hit="+str(self.v3_count)+" | Total Hit="+str(self.total_count)+ " |")
    else:
        pass

    if self.performance=='on':
        print performance.wrecker(GATEWAY_URL,performance_thread,performance_connection,performance_duration)
    else:
        pass

class Rule:
     def add(self,rule_name):
         run_shell("kubectl create -f "+nginmesh_rule_path+rule_name+" -n"+app_namespace,"run")
     def delete(self,rule_name):
         run_shell("kubectl delete -f "+nginmesh_rule_path+rule_name+" -n"+app_namespace,"run")





