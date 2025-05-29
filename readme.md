# Trade-offs made:

Only company-year pairs with data in all sources are scored (happy path).

Missing or zero denominator results in blank metric output (as per requirements).

Code focuses on readability and direct mapping of config logic, not extensibility or speed.

# Scaling:

With 100x more data, use batch or streaming ETL, possibly loading into a DB for fast joins, or distributed processing (Spark, Dask, Dataflow).

# Production:

Deploy as a batch job, with monitoring (metrics on input/output count, time, errors).

Add retries and error handling for data quality or missing files.

Add logging, and unit/integration tests for each metric type and config parse.


## Architecture Diagram

```text
[External Data Sources: APIs, CSV, S3, Partner Feeds]
             |
             |
[ETL/ELT Pipelines] (Golang, Python, Airflow/Argo)
             |
             |
[Data Storage Layer] (Data Warehouse, RDBMS, NoSQL)
             |
             |
[Analytics & Scoring Services](Microservices in Go/Python)
             |
             |
[API Gateway Layer] (REST/gRPC APIs, AuthN/Z)
             |
             |
[Clients, Internal Teams, Dashboards, Integrations]