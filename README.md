# dalv
Querying AWS ALB S3 logs with duckdb


## DuckDB で AWS ALB の S3 ログを読み込むための準備

DuckDB で AWS ALB S3 logs をクエリするためのセットアップと、必要なテーブル作成
```
INSTALL aws;
LOAD aws;
INSTALL httpfs;
LOAD httpfs;
CREATE SECRET (
    TYPE S3,
    PROVIDER CREDENTIAL_CHAIN
);
```

ログを読み込みためのテーブル作成
```aa
CREATE TABLE alb_log_20250303 AS
SELECT *
FROM read_csv(
    's3://{S3_BUCKET_NAME}/xxxxxx/AWSLogs/{ACCOUNT_ID}/elasticloadbalancing/{REGION}/2025/03/03/*.log.gz',
    columns={
        'type': 'VARCHAR',
        'timestamp': 'TIMESTAMP',
        'elb': 'VARCHAR',
        'client_ip_port': 'VARCHAR',
        'target_ip_port': 'VARCHAR',
        'request_processing_time': 'DOUBLE',
        'target_processing_time': 'DOUBLE',
        'response_processing_time': 'DOUBLE',
        'elb_status_code': 'INTEGER',
        'target_status_code': 'VARCHAR',
        'received_bytes': 'BIGINT',
        'sent_bytes': 'BIGINT',
        'request': 'VARCHAR',
        'user_agent': 'VARCHAR',
        'ssl_cipher': 'VARCHAR',
        'ssl_protocol': 'VARCHAR',
        'target_group_arn': 'VARCHAR',
        'trace_id': 'VARCHAR',
        'domain_name': 'VARCHAR',
        'chosen_cert_arn': 'VARCHAR',
        'matched_rule_priority': 'VARCHAR',
        'request_creation_time': 'TIMESTAMP',
        'actions_executed': 'VARCHAR',
        'redirect_url': 'VARCHAR',
        'error_reason': 'VARCHAR',
        'target_port_list': 'VARCHAR',
        'target_status_code_list': 'VARCHAR',
        'classification': 'VARCHAR',
        'classification_reason': 'VARCHAR',
        'conn_trace_id': 'VARCHAR'
    },
    delim=' ',
    quote='"',
    escape='"',
    header=False,
    auto_detect=False
);
```

## DuckDB で AWS ALB の S3 ログをクエリする
```
select * from alb_log_20250303 where elb_status_code != 431;
```
