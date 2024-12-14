The official ClickHouse Docker image doesn’t automatically run an initialization script out-of-the-box like some other databases (e.g. Postgres). However, there are a few common approaches to ensure your table is created at container startup:
Option 1: Using /docker-entrypoint-initdb.d/ (if supported)

Some versions or forks of ClickHouse Docker images support running initialization SQL files found in /docker-entrypoint-initdb.d/ at startup. If you're using such an image (for example, clickhouse/clickhouse-server:latest from ClickHouse’s official repository on Docker Hub), you can do:

    Create a file named init.sql with your CREATE TABLE statement:

-- init.sql
CREATE TABLE default.events (
    activity_id Int32,
    category_uid Int32,
    class_uid Int32,
    severity_id Int32,
    time DateTime,
    type_uid Int32,
    metadata_product_name String,
    metadata_product_vendor_name String,
    metadata_version String,
    app_name Nullable(String),
    app_uid Nullable(String),
    app_vendor_name Nullable(String),
    web_resources_name Array(String),
    web_resources_type Array(String)
)
ENGINE = MergeTree()
ORDER BY time;

Extend the official ClickHouse image with a custom Dockerfile:

    FROM clickhouse/clickhouse-server:latest
    COPY init.sql /docker-entrypoint-initdb.d/init.sql

    Build and run this image. The initialization script should be executed when the container starts, creating the table automatically.

Option 2: Use a Custom Entrypoint

If the above approach doesn’t work due to the image not supporting docker-entrypoint-initdb.d, you can:

    Create a shell script (e.g. entrypoint.sh) that starts the ClickHouse server in the background, waits for it to become available, then runs the CREATE TABLE command, and finally foregrounds the server.

#!/bin/bash
# entrypoint.sh

# Start ClickHouse in the background
clickhouse-server &

# Wait until ClickHouse is ready to accept connections
until clickhouse-client --query="SELECT 1" &>/dev/null; do
    echo "Waiting for ClickHouse to start..."
    sleep 1
done

# Run the initialization SQL
clickhouse-client --query="
CREATE TABLE IF NOT EXISTS default.events (
    activity_id Int32,
    category_uid Int32,
    class_uid Int32,
    severity_id Int32,
    time DateTime,
    type_uid Int32,
    metadata_product_name String,
    metadata_product_vendor_name String,
    metadata_version String,
    app_name Nullable(String),
    app_uid Nullable(String),
    app_vendor_name Nullable(String),
    web_resources_name Array(String),
    web_resources_type Array(String)
) ENGINE = MergeTree() ORDER BY time;
"

# Bring ClickHouse to the foreground
wait

Use a Dockerfile to copy this script into the container and make it the entrypoint:

    FROM clickhouse/clickhouse-server:latest
    COPY entrypoint.sh /entrypoint.sh
    RUN chmod +x /entrypoint.sh
    ENTRYPOINT ["/entrypoint.sh"]

When you run this container, it will start ClickHouse, wait for it to be ready, create the table, and then run indefinitely.
Option 3: Run a Sidecar or One-Time Job

Another approach is to run ClickHouse normally and have another container or script (in your Docker Compose setup) that waits until ClickHouse is ready and then executes the CREATE TABLE statement. For example, using depends_on in docker-compose.yml and a simple container that runs clickhouse-client once and exits.

Recommended Approach:
Try Option 1 first, as it’s the simplest if your chosen ClickHouse image supports initialization scripts via /docker-entrypoint-initdb.d/. If not, Option 2 with a custom entrypoint script is a reliable fallback method.



// replace github.com/MishraLokesh/grpc-server/audit-logging => ./github.com/MishraLokesh/grpc-server/audit-logging/proto
