# ESGBOOK POC
## Trade-offs
- Just scoring company-years where all data exists (assumes clean data, “happy path” only, minimal error handlings/abstractions for POC).

- If any source is missing or divisor is zero, that metric comes out blank.

- Tried to keep code as readable as possible and close to config structure.

## Scaling
- If there was 100x more data, will avoid keeping all CSVs into memory.

- Would probably use a real database for joins, or Spark.

- For streaming/production scale, could build ETL to preprocess and store the latest company-year records, then run the scoring as a batch process.

## Production
- This would run as a scheduled batch job.

- integration with existing logging infra, input/output counts, maybe some simple checks/alerts if anything goes wrong.

- Add error handling and retries for data/file issues.

- Add tests for the metric calculations and config parsing to catch edge cases or config changes.

## notes:
- hardcoded windows path for POC

## Protential Architecture Diagram

```text
[External Data Sources: like APIs, CSV, S3]
             |
             |
[ETL/ELT Pipelines] (Airflow)
             |
             |
[OLAP: Data Storage] (Data Warehouse like Snowflake might be the choice, because the data has been transformed by ETL/ELT pipeline)
             |
             |
[Analytics & Scoring Services](Microservices in Go, this is what POC implemented)
             |
             |
[OLTP: Data Storage](Database like cassandra or Mysql depends on the needs, to serve the webserver)
             |
             |
[API Gateway Layer] (REST/gRPC APIs) 
             |
             |
[Clients, Internal Teams, etc]