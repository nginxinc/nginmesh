import configuration
def wrecker(url,thread=1,connection=10,duration="1s"):
    return configuration.run_shell("wrk -t"+thread+" -c"+connection+" -d"+duration+" http://"+url+":80/productpage | grep -E 'Requests|Transfer|requests|responses'","check")

