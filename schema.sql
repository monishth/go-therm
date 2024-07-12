-- Create the zones table
CREATE TABLE zone (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
    friendly_name TEXT NOT NULL
);

-- Create the thermostats table
CREATE TABLE thermostat (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    friendly_name TEXT NOT NULL,
    topic TEXT NOT NULL,
    zone_id INTEGER,
    FOREIGN KEY (zone_id) REFERENCES zone (id) ON DELETE CASCADE
);

-- Create the valves table
CREATE TABLE valve (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    friendly_name TEXT NOT NULL,
    state_topic TEXT NOT NULL,
    command_topic TEXT NOT NULL,
    zone_id INTEGER,
    FOREIGN KEY (zone_id) REFERENCES zone (id) ON DELETE CASCADE
);
