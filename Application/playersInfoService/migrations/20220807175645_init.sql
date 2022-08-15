-- +goose Up
-- +goose StatementBegin
CREATE TABLE Clubs (
    clubId serial PRIMARY KEY,
    clubName text NOT NULL
);
CREATE TABLE Nationalities (
    nationalityId serial PRIMARY KEY,
    nationalityName text NOT NULL
);
CREATE TABLE Players (
    id serial PRIMARY KEY,
    name text NOT NULL,
    club_id integer REFERENCES Clubs(clubId),
    nationality_id integer REFERENCES Nationalities(nationalityId)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Players;
DROP TABLE Clubs;
DROP TABLE Nationalities;
-- +goose StatementEnd
