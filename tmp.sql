CREATE TABLE snapshot_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    snapshot_id INTEGER NOT NULL,
    data TEXT NOT NULL,
    FOREIGN KEY (snapshot_id) REFERENCES snapshots(id)
);
