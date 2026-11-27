CREATE TABLE FEATURES (
    id VARCHAR(200) NOT NULL,
    date_loaded TIMESTAMP,
    document JSON,
    PRIMARY KEY (id)
);

CREATE INDEX idx_features_date_loaded ON FEATURES (date_loaded);