
use std::collections::HashMap;

// This is from global_dictionary.yaml


pub const REQUEST_HEADER: &str = "request.headers";
pub const TARGET_SERVICE: &str = "target.service";
pub const REQUEST_HOST: &str = "request.host";
pub const REQUEST_METHOD: &str = "request.method";
pub const REQUEST_PATH: &str =  "request.path";
pub const REQUEST_REFER: &str = "request.referer";
pub const REQUEST_SCHEME: &str = "request.scheme";
pub const REQUEST_SIZE: &str = "request.size";
pub const REQUEST_TIME: &str = "request.time";
pub const REQUEST_USERAGENT: &str = "request.useragent";
pub const RESPONSE_CODE: &str = "response.code";
pub const RESPONSE_DURATION: &str = "response.duration";
pub const RESPONSE_SIZE: &str = "response.size";
pub const RESPONSE_HEADERS: &str = "response.headers";
pub const SOURCE_IP: &str = "source.ip";
pub const SOURCE_UID: &str = "source.uid";
pub const SOURCE_PORT: &str = "source.port";
pub const SOURCE_SERVICE: &str = "source.service";
pub const SRC_IP_HEADER: &str = "X-ISTIO-SRC-IP";
pub const SRC_UID_HEADER: &str = "X-ISTIO-SRC-UID";
pub const DESTINATION_SERVICE: &str = "destination.service";
pub const DESTINATION_UID: &str = "destination.uid";
pub const DESTINATION_IP: &str = "destination.ip";



const GLOBAL_LIST: [&'static str; 159] = [
    "source.ip",
    "source.port",
    "source.name",
    "source.uid",
    "source.namespace",
    "source.labels",
    "source.user",
    "target.ip",
    "target.port",
    "target.service",
    "target.name",
    "target.uid",
    "target.namespace",
    "target.labels",
    "target.user",
    "request.headers",
    "request.id",
    "request.path",
    "request.host",
    "request.method",
    "request.reason",
    "request.referer",
    "request.scheme",
    "request.size",
    "request.time",
    "request.useragent",
    "response.headers",
    "response.size",
    "response.time",
    "response.duration",
    "response.code",
    ":authority",
    ":method",
    ":path",
    ":scheme",
    ":status",
    "access-control-allow-origin",
    "access-control-allow-methods",
    "access-control-allow-headers",
    "access-control-max-age",
    "access-control-request-method",
    "access-control-request-headers",
    "accept-charset",
    "accept-encoding",
    "accept-language",
    "accept-ranges",
    "accept",
    "access-control-allow",
    "age",
    "allow",
    "authorization",
    "cache-control",
    "content-disposition",
    "content-encoding",
    "content-language",
    "content-length",
    "content-location",
    "content-range",
    "content-type",
    "cookie",
    "date",
    "etag",
    "expect",
    "expires",
    "from",
    "host",
    "if-match",
    "if-modified-since",
    "if-none-match",
    "if-range",
    "if-unmodified-since",
    "keep-alive",
    "last-modified",
    "link",
    "location",
    "max-forwards",
    "proxy-authenticate",
    "proxy-authorization",
    "range",
    "referer",
    "refresh",
    "retry-after",
    "server",
    "set-cookie",
    "strict-transport-sec",
    "transfer-encoding",
    "user-agent",
    "vary",
    "via",
    "www-authenticate",
    "GET",
    "POST",
    "http",
    "envoy",
    "'200'",
    "Keep-Alive",
    "chunked",
    "x-envoy-service-time",
    "x-forwarded-for",
    "x-forwarded-host",
    "x-forwarded-proto",
    "x-http-method-override",
    "x-request-id",
    "x-requested-with",
    "application/json",
    "application/xml",
    "gzip",
    "text/html",
    "text/html; charset=utf-8",
    "text/plain",
    "text/plain; charset=utf-8",
    "'0'",
    "'1'",
    "true",
    "false",
    "gzip, deflate",
    "max-age=0",
    "x-envoy-upstream-service-time",
    "x-envoy-internal",
    "x-envoy-expected-rq-timeout-ms",
    "x-ot-span-context",
    "x-b3-traceid",
    "x-b3-sampled",
    "x-b3-spanid",
    "tcp",
    "connection.id",
    "connection.received.bytes",
    "connection.received.bytes_total",
    "connection.sent.bytes",
    "connection.sent.bytes_total",
    "connection.duration",
    "context.protocol",
    "context.timestamp",
    "context.time",
    "0",
    "1",
    "200",
    "302",
    "400",
    "401",
    "403",
    "404",
    "409",
    "429",
    "499",
    "500",
    "501",
    "502",
    "503",
    "504",
    "destination.ip",
    "destination.port",
    "destination.service",
    "destination.name",
    "destination.uid",
    "destination.namespace",
    "destination.labels",
    "destination.user",
    "source.service",
];

#[allow(dead_code,non_snake_case)]
pub fn Get_Global_Words() -> [&'static str; 159] {
    GLOBAL_LIST
}


pub struct GlobalDictionary   {

    global_dict: HashMap<String,i32>,

    #[allow(dead_code)]
    top_index: i32
}

impl GlobalDictionary  {


    pub fn new() -> GlobalDictionary  {

        let global_words = GLOBAL_LIST;
        let mut global_dict: HashMap<String,i32> = HashMap::new();
        for  i in 0..global_words.len() {
            let key = GLOBAL_LIST[i];
            global_dict.insert(String::from(key),i as i32);
        }

        GlobalDictionary {
            global_dict,
            top_index: global_words.len() as i32
        }
    }

    // find index in the global dictionary
    pub fn index_of(&self, name: &str ) -> Option<&i32> {
        self.global_dict.get(name)
    }


    pub fn size(&self) -> usize  {
        GLOBAL_LIST.len()
    }

}



