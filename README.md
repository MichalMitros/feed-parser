# feed-parser
Feed fetching and parsing microservice for Heureka recruitment test

### Starting server
Simply use 
`
docker-compose up
`
To start main app server, RabbitMQ and Prometheus.

### RabbitMQ
When all services are up, navigate to http://localhost:15672/ or click [here](http://localhost:15672/) to get to RabbitMQ management console.
Console credentials:
```
username: guest
password: guest
```

### Prometheus
When all services are up, navigate [here](http://localhost:9090/graph?g0.expr=feedparser_parsing_feeds_jobs_current&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h&g1.expr=rate(go_memstats_alloc_bytes_total%5B5m%5D)&g1.tab=0&g1.stacked=0&g1.show_exemplars=0&g1.range_input=1h&g2.expr=rate(go_sched_goroutines_goroutines%5B5m%5D)&g2.tab=0&g2.stacked=0&g2.show_exemplars=0&g2.range_input=1h&g3.expr=rate(feedparser_fetched_xml_files_total%5B5m%5D)&g3.tab=0&g3.stacked=0&g3.show_exemplars=0&g3.range_input=1h&g4.expr=rate(feedparser_fetched_xml_files_failures_total%5B5m%5D)&g4.tab=1&g4.stacked=0&g4.show_exemplars=0&g4.range_input=1h&g5.expr=rate(feedparser_parsed_objects_total%5B5m%5D)&g5.tab=0&g5.stacked=0&g5.show_exemplars=0&g5.range_input=1h&g6.expr=rate(feedparser_rabbitmq_published_items_total%5B5m%5D)&g6.tab=0&g6.stacked=0&g6.show_exemplars=0&g6.range_input=1h&g7.expr=rate(feedparser_rabbitmq_published_items_failures_total%5B5m%5D)&g7.tab=1&g7.stacked=0&g7.show_exemplars=0&g7.range_input=1h&g8.expr=rate(feedparser_requests_total%5B1m%5D)&g8.tab=0&g8.stacked=0&g8.show_exemplars=0&g8.range_input=1h) to get to Prometheus.

### Testing the app
There is Postman collection in the repository with two requests. There are two ways of testing parser:
##### Async request
Asynchronous request is processed in the background. All feeds are parsed concurently. Response should be returned after some miliseconds and cointain just the information that the request was accepted for processing. There is no other ways to check if all files were parsed, than check main app console logs, Prometheus metric `feedparser_parsing_feeds_jobs_current` or RabbitMQ incoming messages rate.  This method is more practical for parsing large feed files when being called by schedulers or other similar services, when information about parsing results is not required immediately.

Postman request: `POST ParseFeedAsync`

cURL: `curl --location --request POST 'localhost:8080/parse-feed-async' --header 'Content-Type: application/json' --data-raw '{
    "feedUrls": [
        "https://e.mall.cz/cz-mall-heureka.xml",
        "https://exports.conviu.com/download/6x0iqs2sav7psvzm2l0uyu5wv2jztbrn/writer/yj582due89jrulj2gj9j8dipnbv497nx.xml"
    ]
}'`

##### Nonasync request
In this case all feed files are also processed cuncurently, but the response will be returned when processing of all files is done. The rsponse contain processing final status and parsing time of all urls from the request. As it waits for all files to process, this request is not too practical to use with services with short timeout, because if the largest feed takes 10 minutes to proceed, then the response is returned after 10 minutes. Nevertheless it's more convinient endpoint for testing.

Postman request: `POST ParseFeed`

cURL: `curl --location --request POST 'localhost:8080/parse-feed' --header 'Content-Type: application/json' --data-raw '{
    "feedUrls": [
        "https://e.mall.cz/cz-mall-heureka.xml",
        "https://exports.conviu.com/download/6x0iqs2sav7psvzm2l0uyu5wv2jztbrn/writer/yj582due89jrulj2gj9j8dipnbv497nx.xml"
    ]
}'`






