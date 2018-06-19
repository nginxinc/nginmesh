package nginx

const httpTemplate = `{{range $upstream := .Upstreams}}
upstream {{$upstream.Name}} {
    {{if $upstream.LBMethod}}{{$upstream.LBMethod}};{{end}}
	{{range $server := $upstream.Servers}}
    server {{$server.Address}}:{{$server.Port}} max_conns={{$server.MaxConns}} max_fails={{$server.MaxFails}}{{if $server.FailTimeout}} fail_timeout={{$server.FailTimeout}}s{{end}};{{end}}

    keepalive 32;
}{{end}}

{{range $map := .Maps}}
map {{$map.Source}} {{$map.Variable}} {
    {{range $k, $v := $map.Params}}
    {{if eq $k "default"}}default{{else}}"{{$k}}"{{end}} {{$v}};
    {{end}}
}
{{end}}


{{range $server := .VirtualServers}}

{{range $split := $server.SplitClients}}
split_clients $request_id {{$split.Variable}} {
    {{range $dist := $split.Distributions}}
    {{$dist.Weight}}% {{$dist.Value}};
    {{end}}
}
{{end}}

server {
    {{if $server.SSL}}
    listen {{$server.Address}}:{{$server.Port}} ssl;
    ssl_certificate {{$server.SSL.Certificate}};
    ssl_certificate_key {{$server.SSL.Key}};
    ssl_trusted_certificate {{$server.SSL.TrustedCertificate}};
    ssl_client_certificate {{$server.SSL.TrustedCertificate}};

    ssl_verify_client on;
    {{else}}
    listen {{$server.Address}}:{{$server.Port}}{{if $server.IsTarget}} default_server{{end}};
    {{end}}

    server_name {{range $name := $server.Names}} {{$name}}{{end}};

 {{if $.Mixer}}
#    collector_destination_ip  {{$.Mixer.DestinationIP}};
#    collector_destination_uid {{$.Mixer.DestinationUID}};
    {{if $.Mixer.DestinationService}}
#    collector_destination_service {{$.Mixer.DestinationService}};
    {{end}}
#    collector_source_ip {{$.Mixer.SourceIP}};
#    collector_source_uid {{$.Mixer.SourceUID}};
    {{end}}


    {{range $location := $server.Locations}}
    location {{$location.Path}} {
        {{range $set := $location.Sets}}
        set {{$set.Variable}} '{{$set.Value}}';
        {{end}}

        {{if $location.Internal}}internal;{{end}}

        {{if $location.MixerReport}}
        {{if $location.CollectorTopic}}
        collector_report {{$location.CollectorTopic}};
        {{end}}
        {{end}}


       
        {{if $location.Tracing}}
        opentracing_operation_name $host:$server_port;
        opentracing_trace_locations off;
        {{end}}

        {{if $location.Redirect}}
        return {{$location.Redirect.Code}} {{$location.Redirect.URL}};
        {{end}}

        {{range $expression := $location.Expressions}}
        if ({{$expression.Condition}}) {
            {{$expression.Result}};
        }
        {{end}}
        
        {{if $location.Upstream}}
        proxy_set_header Host {{if $location.Host}}{{$location.Host}}{{else}}$host{{end}};

        {{if $.Mixer}}
        {{if $.Mixer.SourceIP}}
        proxy_set_header X-ISTIO-SRC-IP {{$.Mixer.SourceIP}};
        {{end}}
        
        {{if $.Mixer.SourceUID}}
        proxy_set_header X-ISTIO-SRC-UID {{$.Mixer.SourceUID}};
        {{end}}
        {{end}}

        # WebSocket and KeepAlives
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        
        {{if ne $location.ConnectTimeout 0}}	proxy_connect_timeout {{$location.ConnectTimeout}}ms;{{end}}

        {{if $location.ProxyNextUpstream}}
        proxy_next_upstream {{$location.ProxyNextUpstream.Condition}}; 
        proxy_next_upstream_timeout {{$location.ProxyNextUpstream.Timeout}}ms; 
        proxy_next_upstream_tries {{$location.ProxyNextUpstream.Tries}};
        {{end}}

        {{if $location.Rewrite}}
        rewrite ^{{$location.Rewrite.Prefix}}(.*)$ {{$location.Rewrite.Replacement}}$1 break;
        {{end}}

        {{if $location.SSL}}
        proxy_ssl_certificate {{$location.SSL.Certificate}};
        proxy_ssl_certificate_key {{$location.SSL.Key}};
        proxy_ssl_trusted_certificate {{$location.SSL.TrustedCertificate}};
        #{{if $location.SSL.Name}}proxy_ssl_name {{$location.SSL.Name}};{{end}}
        #proxy_ssl_verify on;
        #proxy_ssl_server_name on;
        proxy_ssl_session_reuse on;
        proxy_pass https://{{$location.Upstream}};
        {{else}}
        proxy_pass http://{{$location.Upstream}};
        {{end}}
        {{end}}
    }{{end}}
}{{end}}`

const mainTemplate = `load_module /etc/nginx/modules/ngx_stream_nginmesh_dest_module.so;

load_module /etc/nginx/modules/ngx_http_opentracing_module.so;
load_module /etc/nginx/modules/ngx_http_zipkin_module.so;


worker_processes  auto;

error_log  /dev/stdout {{.LOGLEVEL}};
pid        /etc/istio/proxy/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for" '
                      'scr_ip="$http_x_istio_src_ip" src_uid="$http_x_istio_src_uid" host="$host" '
                      'verify="$ssl_client_verify" sni="$ssl_server_name"';

    access_log  /dev/stdout  main;

    {{if .Tracing}}
   # Enable tracing for all requests.
    opentracing on;

    # Zipkin additional tags
    opentracing_tag guid:x-request-id $request_id;
    opentracing_tag http.protocol $server_protocol;
    opentracing_tag user_agent $http_user_agent;
    opentracing_tag response_size $bytes_sent;
    opentracing_tag request_size $request_length;
    opentracing_tag node_id {{.ServiceNode}};

    #Zipkin IP/Port/Service_name
    zipkin_collector_host zipkin.istio-system;
    zipkin_collector_port 9411;
    zipkin_service_name  {{.ServiceCluster}};
    {{end}}

    proxy_temp_path /etc/istio/proxy/cache/proxy_temp;
    client_body_temp_path /etc/istio/proxy/cache/client_temp;
    fastcgi_temp_path /etc/istio/proxy/cache/fastcgi_temp;
    uwsgi_temp_path /etc/istio/proxy/cache/uwsgi_temp;
    scgi_temp_path /etc/istio/proxy/cache/scgi_temp;


    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    server_names_hash_bucket_size 128;
    variables_hash_bucket_size 128;

  #  collector_server {{.CollectorServer}};


    # Support for Websocket
    map $http_upgrade $connection_upgrade {
        default upgrade;
        '' '';
    }

    include /etc/istio/proxy/conf.d/*.conf;
}

stream {
    log_format basic '$remote_addr [$time_local] '
                     '$protocol $status $bytes_sent $bytes_received '
                     '$session_time $nginmesh_dest $nginmesh_server';
    map $nginmesh_dest $nginmesh_server {
        {{range $dm := .DestinationMaps}}
        {{- $dm.Remote}} {{$dm.Local}};
        {{end}}
    }
    server {
        listen 15001;
        access_log /dev/stdout basic;
        nginmesh_dest on;
        proxy_pass $nginmesh_server;
    }
    include /etc/istio/proxy/conf.d/*.stream-conf;
}`

const tcpTemplate = `{{range $upstream := .Upstreams}}
upstream {{$upstream.Name}} {
    {{if $upstream.LBMethod}}{{$upstream.LBMethod}};{{end}}
	{{range $server := $upstream.Servers}}
	server {{$server.Address}}:{{$server.Port}};{{end}}
}{{end}}


{{range $server := .Servers}}
server {
   listen {{$server.Address}}:{{$server.Port}};
   {{if ne $server.ConnectTimeout 0}} proxy_connect_timeout {{$server.ConnectTimeout}}ms;{{end}}
   proxy_pass {{$server.Upstream}};

}
{{end}}`
