-- init metadata
SELECT InitSpatialMetadata(1);

CREATE TABLE areas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fid INTEGER UNIQUE,
    name TEXT
);

SELECT AddGeometryColumn('areas', 'geom', 0, 'MULTIPOLYGON', 'XY');

CREATE TABLE content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    area_fid INTEGER,
    obj_type TEXT,
    obj_name TEXT,
    description TEXT,
    FOREIGN KEY(area_fid) REFERENCES areas(fid)
);

CREATE TABLE points (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    x REAL,
    y REAL
);

SELECT AddGeometryColumn('points', 'geom', 0, 'POINT', 'XY');
