CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS stadiums (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    capacity INTEGER NOT NULL,
    city TEXT NOT NULL,
    country_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS leagues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    number_of_teams INTEGER NOT NULL,
    country_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name TEXT NOT NULL UNIQUE,
    medium_name TEXT NOT NULL,
    acronym TEXT NOT NULL,
    nickname TEXT,
    year_founded INTEGER NOT NULL,
    city TEXT,
    country_id UUID NOT NULL,
    stadium_id UUID NOT NULL,
    league_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS persons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    middle_names TEXT,
    last_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    country_id UUID NOT NULL,
    UNIQUE (first_name, last_name, date_of_birth)
);

CREATE TABLE IF NOT EXISTS team_players (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    person_id UUID NOT NULL,
    team_id UUID NOT NULL,
    squad_number INTEGER NOT NULL,
    general_position TEXT NOT NULL,
    specific_position TEXT,
    started DATE NOT NULL,
    ended DATE,
    UNIQUE (person_id, team_id, started)
);

CREATE TABLE IF NOT EXISTS team_managers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    person_id UUID NOT NULL,
    team_id UUID NOT NULL,
    started DATE NOT NULL,
    ended DATE
);
