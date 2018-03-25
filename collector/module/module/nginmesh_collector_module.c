/**
 * @file   nginmesh_collector_module.c
 * @author Sehyo Chang <sehyo@nginx.com>
 * @date   Wed Aug 19 2017
 *
 * @brief  Collector module for nginmesh
 *
 * @section LICENSE
 *
 * Copyright (C) 2017,2018 by Nginx
 *
 */
#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>


typedef struct {
    ngx_str_t     topic;
    ngx_str_t     destination_service;

} ngx_http_collector_loc_conf_t;


/**
 * @brief element collector configuration
 */
typedef struct {
    ngx_str_t collector_server;              /**< collector server */
} ngx_http_collector_main_conf_t;


typedef struct  {
    ngx_str_t     destination_service;        // destination service
    ngx_str_t     destination_uid;           // destination service
    ngx_str_t     destination_ip;           // destination ip address
    ngx_str_t     source_ip;                // source ip
    ngx_str_t     source_uid;               // source uid
    ngx_str_t     source_service;           // source service
    ngx_uint_t     source_port;              // source port

} ngx_http_collector_srv_conf_t;


static ngx_int_t ngx_http_collector_report_handler(ngx_http_request_t *r);


static ngx_int_t ngx_http_collector_filter_init(ngx_conf_t *cf);

// create configuration
static void *ngx_http_collector_create_loc_conf(ngx_conf_t *cf);
static char *ngx_http_collector_merge_loc_conf(ngx_conf_t *cf, void *parent,void *child);

static void *ngx_http_collector_create_srv_conf(ngx_conf_t *cf);
static char *ngx_http_collector_merge_srv_conf(ngx_conf_t *cf, void *parent, void *child);

static void *ngx_http_collector_create_main_conf(ngx_conf_t *cf);
static char *nginmesh_http_collector_server_post(ngx_conf_t *cf, void *data, void *conf);    

static ngx_conf_post_handler_pt  ngx_http_collector_server_p =
    nginmesh_http_collector_server_post;

// handlers in rust
void  nginmesh_set_collector_server_config(ngx_str_t *server);
void  nginmesh_collector_report_handler(ngx_http_request_t *r, 
        ngx_http_collector_main_conf_t *main_conf,
        ngx_http_collector_srv_conf_t *srv_conf,
        ngx_http_collector_loc_conf_t *loc_conf);

ngx_int_t  nginmesh_collector_init(ngx_cycle_t *cycle);
void  nginmesh_collector_exit();



/**
 * This module provide callback to istio for http traffic
 *
 */
static ngx_command_t ngx_http_collector_commands[] = {

    { 
      ngx_string("collector_report"),   /* report directive */
      NGX_HTTP_LOC_CONF | NGX_CONF_FLAG, 
      ngx_conf_set_str_slot, /* configuration setup function */
      NGX_HTTP_LOC_CONF_OFFSET, 
      offsetof(ngx_http_collector_loc_conf_t, topic),  // store in the location configuration
      NULL
    },
    {
       ngx_string("collector_destination_service"), /* directive */
       NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
       ngx_conf_set_str_slot, /* configuration setup function */
       NGX_HTTP_SRV_CONF_OFFSET,
       offsetof(ngx_http_collector_srv_conf_t, destination_service),  // store in the location configuration
       NULL
     },
     {
        ngx_string("collector_destination_uid"), /* directive */
        NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
        ngx_conf_set_str_slot, /* configuration setup function */
        NGX_HTTP_SRV_CONF_OFFSET,
        offsetof(ngx_http_collector_srv_conf_t, destination_uid),  // store in the location configuration
        NULL
     },
     {
      ngx_string("collector_destination_ip"), /* directive */
      NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
      ngx_conf_set_str_slot, /* configuration setup function */
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_collector_srv_conf_t, destination_ip),  // store in the location configuration
      NULL
    },
    {
      ngx_string("collector_source_ip"),
      NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
      ngx_conf_set_str_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_collector_srv_conf_t, source_ip),  // store in the location configuration
      NULL
    },

    {
      ngx_string("collector_source_uid"),
      NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
      ngx_conf_set_str_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_collector_srv_conf_t, source_uid),  // store in the location configuration
      NULL
    },
    {
      ngx_string("collector_source_service"),
      NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
      ngx_conf_set_str_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_collector_srv_conf_t, source_service),  // store in the location configuration
      NULL
    },
    {
      ngx_string("collector_source_port"),
      NGX_HTTP_SRV_CONF | NGX_CONF_TAKE1,
      ngx_conf_set_num_slot,
      NGX_HTTP_SRV_CONF_OFFSET,
      offsetof(ngx_http_collector_srv_conf_t, source_port),  // store in the location configuration
      NULL
    },
    { 
      ngx_string("collector_server"), /* directive */
      NGX_HTTP_MAIN_CONF|NGX_CONF_TAKE1,  // server takes 1 //
      ngx_conf_set_str_slot, /* configuration setup function */
      NGX_HTTP_MAIN_CONF_OFFSET, 
      offsetof(ngx_http_collector_main_conf_t,collector_server),
      &ngx_http_collector_server_p
    },
    ngx_null_command /* command termination */
};


/* The module context. */
static ngx_http_module_t ngx_http_collector_module_ctx = {
    NULL, /* preconfiguration */
    ngx_http_collector_filter_init, /* postconfiguration */
    ngx_http_collector_create_main_conf, /* create main configuration */
    NULL, /* init main configuration */

    ngx_http_collector_create_srv_conf, /* create server configuration */
    ngx_http_collector_merge_srv_conf, /* merge server configuration */

    ngx_http_collector_create_loc_conf, /* create location configuration */
    ngx_http_collector_merge_loc_conf /* merge location configuration */
};

/* Module definition. */
ngx_module_t ngx_http_collector_module = {
    NGX_MODULE_V1,
    &ngx_http_collector_module_ctx, /* module context */
    ngx_http_collector_commands, /* module directives */
    NGX_HTTP_MODULE, /* module type */
    NULL, /* init master */
    NULL, /* init module */
    nginmesh_collector_init, /* init process */
    NULL, /* init thread */
    NULL, /* exit thread */
    NULL, /* exit process */
    NULL, /* exit master */
    NGX_MODULE_V1_PADDING
};

// install log phase handler for collector
static ngx_int_t ngx_http_collector_filter_init(ngx_conf_t *cf) {


    ngx_http_handler_pt        *h1;
    //ngx_http_handler_pt          *h2;
    ngx_http_core_main_conf_t  *cmcf;
    //ngx_http_core_loc_conf_t  *clcf;

    cmcf = ngx_http_conf_get_module_main_conf(cf, ngx_http_core_module);

    h1 = ngx_array_push(&cmcf->phases[NGX_HTTP_LOG_PHASE].handlers);
    if (h1 == NULL) {
        return NGX_ERROR;
    }
    *h1 = ngx_http_collector_report_handler;

    ngx_log_debug(NGX_LOG_DEBUG_EVENT, ngx_cycle->log, 0, "registering collector report handler");

    return NGX_OK;   
}

/**
 * collector report handler.
 *
 */
static ngx_int_t ngx_http_collector_report_handler(ngx_http_request_t *r)
{
    ngx_http_collector_loc_conf_t  *loc_conf;
    ngx_http_collector_main_conf_t *main_conf;
    ngx_http_collector_srv_conf_t *srv_conf;

    ngx_log_debug(NGX_LOG_DEBUG_HTTP,  r->connection->log, 0, "start invoking collector report handler");

    loc_conf = ngx_http_get_module_loc_conf(r, ngx_http_collector_module);
    srv_conf = ngx_http_get_module_srv_conf(r,ngx_http_collector_module);
    main_conf = ngx_http_get_module_main_conf(r, ngx_http_collector_module);

    ngx_log_debug2(NGX_LOG_DEBUG_HTTP,  r->connection->log, 0, "using collector server: %*s",main_conf->collector_server.len,main_conf->collector_server.data);

    // invoke mix client
    nginmesh_collector_report_handler(r,main_conf,srv_conf,loc_conf);

    ngx_log_debug(NGX_LOG_DEBUG_HTTP,  r->connection->log, 0, "finish calling collector report handler");


   return NGX_OK;

} 

// create loc conf for collector
static void *ngx_http_collector_create_loc_conf(ngx_conf_t *cf) {

    ngx_http_collector_loc_conf_t  *conf;

    conf = ngx_pcalloc(cf->pool, sizeof(ngx_http_collector_loc_conf_t));
    if (conf == NULL) {
        return NULL;
    }

    ngx_log_debug(NGX_LOG_DEBUG_EVENT, ngx_cycle->log, 0, "set up  collector location config");

    return conf;
}

static char *ngx_http_collector_merge_loc_conf(ngx_conf_t *cf, void *parent, void *child)
{
    ngx_log_debug(NGX_LOG_DEBUG_EVENT, ngx_cycle->log, 0, "merging loc conf");

    ngx_http_collector_loc_conf_t  *prev = parent;
    ngx_http_collector_loc_conf_t  *conf = child;

    ngx_conf_merge_str_value(conf->topic, prev->topic, "");
    ngx_conf_merge_str_value(conf->destination_service,prev->destination_service,"")

    return NGX_CONF_OK;
}

static void *ngx_http_collector_create_srv_conf(ngx_conf_t *cf) {

    ngx_http_collector_srv_conf_t  *conf;

    conf = ngx_pcalloc(cf->pool, sizeof(ngx_http_collector_srv_conf_t));
    if (conf == NULL) {
        return NULL;
    }

    conf->source_port = NGX_CONF_UNSET_UINT;

    ngx_log_debug(NGX_LOG_DEBUG_EVENT, ngx_cycle->log, 0, "set up collector srv config");

    return conf;
}


static char *ngx_http_collector_merge_srv_conf(ngx_conf_t *cf, void *parent, void *child)
{
    ngx_log_debug(NGX_LOG_DEBUG_EVENT, ngx_cycle->log, 0, "merging srv conf");

    ngx_http_collector_srv_conf_t  *prev = parent;
    ngx_http_collector_srv_conf_t  *conf = child;

    ngx_conf_merge_str_value(conf->destination_service,prev->destination_service,"");
    ngx_conf_merge_str_value(conf->destination_uid,prev->destination_uid,"");
    ngx_conf_merge_str_value(conf->source_ip,prev->source_ip,"");
    ngx_conf_merge_str_value(conf->source_uid,prev->source_uid,"");
    ngx_conf_merge_str_value(conf->source_service,prev->source_service,"");
    ngx_conf_merge_uint_value(conf->source_port, prev->source_port, 0);

    return NGX_CONF_OK;
}


static void *ngx_http_collector_create_main_conf(ngx_conf_t *cf)
{
  ngx_http_collector_main_conf_t *conf;

  ngx_log_debug(NGX_LOG_DEBUG_EVENT, ngx_cycle->log, 0, "setting up main config");

  conf = ngx_pcalloc(cf->pool, sizeof(ngx_http_collector_main_conf_t));
  if (conf == NULL) {
    return NULL;
  }

  return conf;
}

// set collector server
static char *nginmesh_http_collector_server_post(ngx_conf_t *cf, void *post, void *data)
{
    ngx_str_t  *server = data;


    ngx_log_debug(NGX_LOG_DEBUG_HTTP,  cf->log, 0, "start invoking main collector post");


    nginmesh_set_collector_server_config(server);

    ngx_log_debug(NGX_LOG_DEBUG_HTTP, cf->log, 0, "finish calling main collector report handler");


   return NGX_OK;
}
