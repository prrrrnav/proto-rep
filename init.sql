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
