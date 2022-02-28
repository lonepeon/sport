CREATE TABLE runs (
  id TEXT PRIMARY KEY,
  ran_at TEXT,
  duration TEXT,
  distance REAL,
  speed REAL,
  gpx_path TEXT,
  map_path TEXT,
  created_at TEXT
)
